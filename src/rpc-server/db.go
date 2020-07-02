package main

import (
	context "context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const _LeakCollectionName = "leaks"
const _EmailCollectionName = "emails"
const _RelationCollectionName = "rels"

// LeakEntry is used for temporary data holding for 'leaks' collection per database communications
type LeakEntry struct {
	ID   primitive.ObjectID `bson:"_id"`  // Object ID
	Name string             `bson:"name"` // Leak name
}

// EmailEntry is used for temporary data holding for 'emails' collection per database communications
type EmailEntry struct {
	ID        primitive.ObjectID `bson:"_id"`        // Object ID
	Domain    string             `bson:"domain"`     // Email domain
	CreatedAt int64              `bson:"created_at"` // Email creation date
	UpdatedAt int64              `bson:"updated_at"` // Last operation including this email date
	Email     string             `bson:"email"`      // Email itself
}

// LeakEmailRelationEntry is used for temporary data holding for 'rels' collection per database communications
type LeakEmailRelationEntry struct {
	EmailID     primitive.ObjectID `bson:"email_id"`     // Email Object ID
	LeakID      primitive.ObjectID `bson:"leak_id"`      // Leak Object ID
	EmailDomain string             `bson:"email_domain"` // Email domain
}

// MongoDBConn - Custom DB connector class
type MongoDBConn struct {
	conn *mongo.Database
}

// Connect performs DB connection for the client object initialization
func (db *MongoDBConn) Connect(uri string, database string) error {
	if db == nil {
		return errors.New("I am not even alive")
	}

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		return err
	}

	db.conn = client.Database(database)
	return nil
}

// GetAllLeaks returns an array of all leak documents in the DB
func (db *MongoDBConn) GetAllLeaks() ([]*LeakEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var leaks []*LeakEntry

	ctx := context.TODO()
	cur, err := db.conn.Collection(_LeakCollectionName).Find(ctx, bson.D{{}})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var l LeakEntry
		err := cur.Decode(&l)

		if err != nil {
			return nil, err
		}

		leaks = append(leaks, &l)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	// once exhausted, close the cursor
	cur.Close(ctx)

	return leaks, nil
}

// GetEmailsByLeakID returns an array of email documents in the DB
// which have a relation with the Leak specified
func (db *MongoDBConn) GetEmailsByLeakID(leakId string) ([]*EmailEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	leakOID, err := primitive.ObjectIDFromHex(leakId)

	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "leak_id", Value: leakOID}})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var r LeakEmailRelationEntry
		err := cur.Decode(&r)

		if err != nil {
			return nil, err
		}

		relations = append(relations, &r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	var emails []*EmailEntry

	for _, rel := range relations {
		var emailEntry EmailEntry

		err = db.conn.Collection(_EmailCollectionName).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: rel.EmailID}}).Decode(&emailEntry)

		if err != nil {
			return nil, err
		}

		emails = append(emails, &emailEntry)
	}

	return emails, nil
}

// GetLeaksByEmailID returns an array of leak documents in the DB
// which have a relation with the Email specified
func (db *MongoDBConn) GetLeaksByEmailID(emailId string) ([]*LeakEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	ctx := context.TODO()

	emailOID, err := primitive.ObjectIDFromHex(emailId)

	if err != nil {
		return nil, err
	}

	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "email_id", Value: emailOID}})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var r LeakEmailRelationEntry
		err := cur.Decode(&r)

		if err != nil {
			return nil, err
		}

		relations = append(relations, &r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	var leaks []*LeakEntry

	for _, rel := range relations {
		var leakEntry LeakEntry

		err = db.conn.Collection(_LeakCollectionName).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: rel.LeakID}}).Decode(&leakEntry)

		if err != nil {
			return nil, err
		}

		leaks = append(leaks, &leakEntry)
	}

	return leaks, nil
}

// GetLeaksByDomain returns an array of leak documents in the DB
// which have a relation with the email domain specified
func (db *MongoDBConn) GetLeaksByDomain(domain string) ([]*LeakEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "email_domain", Value: domain}})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var r LeakEmailRelationEntry
		err := cur.Decode(&r)

		if err != nil {
			return nil, err
		}

		relations = append(relations, &r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	var leaks []*LeakEntry

	for _, rel := range relations {
		var leakEntry LeakEntry

		err = db.conn.Collection(_LeakCollectionName).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: rel.LeakID}}).Decode(&leakEntry)

		if err != nil {
			return nil, err
		}

		_leakExists := false

		for _, leakPtr := range leaks {
			if leakPtr.ID.Hex() == leakEntry.ID.Hex() {
				_leakExists = true
				break
			}
		}

		if !_leakExists {
			leaks = append(leaks, &leakEntry)
		}
	}

	return leaks, nil
}

// GetEmailsByDomainAndLeakID returns an array of email documents in the DB
// which have a relation with both the email domain and the Leak specified
func (db *MongoDBConn) GetEmailsByDomainAndLeakID(domain string, leakId string) ([]*EmailEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	leakOID, err := primitive.ObjectIDFromHex(leakId)

	if err != nil {
		return nil, err
	}

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "email_domain", Value: domain}, primitive.E{Key: "leak_id", Value: leakOID}})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var r LeakEmailRelationEntry
		err := cur.Decode(&r)

		if err != nil {
			return nil, err
		}

		relations = append(relations, &r)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	cur.Close(ctx)

	var emails []*EmailEntry

	for _, rel := range relations {
		var emailEntry EmailEntry

		err = db.conn.Collection(_EmailCollectionName).FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: rel.EmailID}}).Decode(&emailEntry)

		if err != nil {
			return nil, err
		}

		emails = append(emails, &emailEntry)
	}

	return emails, nil
}

// GetEmailIDFromEmail returns the Object ID in hex string format of the email specified
func (db *MongoDBConn) GetEmailIDFromEmail(email string) (string, error) {
	var emailEnt EmailEntry

	err := db.conn.Collection(_EmailCollectionName).FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: email}}).Decode(&emailEnt)

	if err != nil {
		return "", err
	}

	return emailEnt.ID.Hex(), nil
}
