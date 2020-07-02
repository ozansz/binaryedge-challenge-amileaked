package main

import (
	context "context"
	"errors"

	empty "github.com/golang/protobuf/ptypes/empty"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// LeakServiceServerHandler implements proper data structure for service server
type LeakServiceServerHandler struct {
	DBConnURI    string       // Database server connection URI
	DatabaseName string       // Database name to connect
	db           *MongoDBConn // MongoDB client connector ptr
}

// DBConnect - Makes connection to database cluster
func (s *LeakServiceServerHandler) DBConnect() error {
	// Perform nil pointer check
	if s == nil {
		return errors.New("Uninitialized server stub")
	}

	s.db = &MongoDBConn{}

	if err := s.db.Connect(s.DBConnURI, s.DatabaseName); err != nil {
		return err
	}

	return nil
}

// ListLeaks - Lists all Leak objects with their related Emails
func (s *LeakServiceServerHandler) ListLeaks(req *empty.Empty, srv LeakService_ListLeaksServer) error {
	// Perform nil pointer check
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	// Perform nil pointer check
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

// GetLeaksByEmailStreamed - Open a one direction stream and send Leak objects related to the Email specified
func (s *LeakServiceServerHandler) GetLeaksByEmailStreamed(req *GetLeaksByEmailRequest, srv LeakService_GetLeaksByEmailStreamedServer) error {
	// Perform nil pointer check
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	// Perform nil pointer check
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

// GetLeaksByDomainStreamed - Open a one direction stream and send Leak objects related to the email domain specified
func (s *LeakServiceServerHandler) GetLeaksByDomainStreamed(req *GetLeaksByDomainRequest, srv LeakService_GetLeaksByDomainStreamedServer) error {
	// Perform nil pointer check
	if s == nil {
		return status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	// Perform nil pointer check
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

// GetLeaksByEmail - Return an array of Leak objects related to the Email specified
func (s *LeakServiceServerHandler) GetLeaksByEmail(ctx context.Context, req *GetLeaksByEmailRequest) (*GetLeaksByEmailResponse, error) {
	// Perform nil pointer check
	if s == nil {
		return nil, status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	// Perform nil pointer check
	if s.db == nil {
		return nil, status.Errorf(codes.Aborted, "Uninitialized database connection")
	}

	emailID, err := s.db.GetEmailIDFromEmail(req.GetEmail())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to get Email object ID")
	}

	leaks, err := s.db.GetLeaksByEmailID(emailID)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to get Leak objects")
	}

	leaksRet := make([]*Leak, 0)

	for _, leak := range leaks {
		leakPkg := &Leak{
			Id:         leak.ID.Hex(),
			Name:       leak.Name,
			Emails:     make([]*Leak_Email, 0),
			EmailCount: 0,
		}

		leaksRet = append(leaksRet, leakPkg)
	}

	return &GetLeaksByEmailResponse{
		Leaks: leaksRet,
	}, nil
}

// GetLeaksByDomain - Return an array of Leak objects related to the email domain specified
func (s *LeakServiceServerHandler) GetLeaksByDomain(ctx context.Context, req *GetLeaksByDomainRequest) (*GetLeaksByDomainResponse, error) {
	// Perform nil pointer check
	if s == nil {
		return nil, status.Errorf(codes.Aborted, "Uninitialized server stub")
	}

	// Perform nil pointer check
	if s.db == nil {
		return nil, status.Errorf(codes.Aborted, "Uninitialized database connection")
	}

	leaks, err := s.db.GetLeaksByDomain(req.GetDomain())

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Unable to get Leak objects")
	}

	leaksRet := make([]*Leak, 0)

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
			return nil, status.Errorf(codes.Internal, "Unable to get Email objects")
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

		leaksRet = append(leaksRet, leakPkg)
	}

	return &GetLeaksByDomainResponse{
		Leaks: leaksRet,
	}, nil
}
