package main

import (
	"errors"

	empty "github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type LeakServiceServerHandler struct {
	DBConnURI    string
	DatabaseName string
	db           *MongoDBConn
}

func (s *LeakServiceServerHandler) DBConnect() error {
	if s == nil {
		return errors.New("Uninitialized server stub")
	}

	s.db = &MongoDBConn{}

	if err := s.db.Connect(s.DBConnURI, s.DatabaseName); err != nil {
		return err
	}

	return nil
}

func (s *LeakServiceServerHandler) ListLeaks(req *empty.Empty, srv LeakService_ListLeaksServer) error {
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	if s.db == nil {
		return status.Errorf(codes.Aborted, "Uninitialized database connection")
	}

	leaks, err := s.db.GetAllLeaks()

	if err != nil {
		return status.Errorf(codes.Internal, "Unable to get Leak objects")
	}

	for _, leak := range leaks {
		leakPkg := &Leak{
			Id:         leak.ID.Hex(),
			Name:       leak.Name,
			EmailCount: 0,
			Emails:     nil,
		}

		emailsArr := make([]*Leak_Email, 0)
		var emailCount int64 = 0

		emails, err := s.db.GetEmailsByLeakID(leakPkg.GetId())

		if err != nil {
			return status.Errorf(codes.Internal, "Unable to get Email objects")
		}

		for _, email := range emails {
			emailCount++

			emailsArr = append(emailsArr, &Leak_Email{
				Email:            email.Email,
				Domain:           email.Domain,
				FirstOccuranceTs: email.CreatedAt.Unix(),
				LastOccuranceTs:  email.UpdatedAt.Unix(),
			})
		}

		leakPkg.Emails = emailsArr
		leakPkg.EmailCount = emailCount

		if err := srv.Send(leakPkg); err != nil {
			return err
		}
	}

	return nil
}

func (s *LeakServiceServerHandler) GetLeaksByEmail(req *GetLeaksByEmailRequest, srv LeakService_GetLeaksByEmailServer) error {
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	if s.db == nil {
		return status.Errorf(codes.Aborted, "Uninitialized database connection")
	}

	emailID, err := s.db.GetEmailIDFromEmail(req.GetEmail())

	if err != nil {
		return status.Errorf(codes.Internal, "Unable to get Email object ID")
	}

	leaks, err := s.db.GetLeaksByEmailID(emailID)

	if err != nil {
		return status.Errorf(codes.Internal, "Unable to get Leak objects")
	}

	for _, leak := range leaks {
		leakPkg := &Leak{
			Id:         leak.ID.Hex(),
			Name:       leak.Name,
			Emails:     make([]*Leak_Email, 0),
			EmailCount: 0,
		}

		if err := srv.Send(leakPkg); err != nil {
			return err
		}
	}

	return nil
}

func (s *LeakServiceServerHandler) GetLeaksByDomain(req *GetLeaksByDomainRequest, srv LeakService_GetLeaksByDomainServer) error {
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	if s.db == nil {
		return status.Errorf(codes.Aborted, "Uninitialized database connection")
	}

	leaks, err := s.db.GetLeaksByDomain(req.GetDomain())

	if err != nil {
		return status.Errorf(codes.Internal, "Unable to get Leak objects")
	}

	for _, leak := range leaks {
		leakPkg := &Leak{
			Id:         leak.ID.Hex(),
			Name:       leak.Name,
			Emails:     nil,
			EmailCount: 0,
		}

		emailsArr := make([]*Leak_Email, 0)
		var emailCount int64 = 0

		emails, err := s.db.GetEmailsByDomainAndLeakID(req.GetDomain(), leakPkg.GetId())

		if err != nil {
			return status.Errorf(codes.Internal, "Unable to get Email objects")
		}

		for _, email := range emails {
			emailCount++

			emailsArr = append(emailsArr, &Leak_Email{
				Email:            email.Email,
				Domain:           email.Domain,
				FirstOccuranceTs: email.CreatedAt.Unix(),
				LastOccuranceTs:  email.UpdatedAt.Unix(),
			})
		}

		leakPkg.Emails = emailsArr
		leakPkg.EmailCount = emailCount

		if err := srv.Send(leakPkg); err != nil {
			return err
		}
	}

	return nil
}
