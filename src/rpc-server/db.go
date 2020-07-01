package main

import (
	context "context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const _LeakCollectionName = "leaks"
const _EmailCollectionName = "emails"
const _RelationCollectionName = "rels"

type LeakEntry struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}
type EmailEntry struct {
	ID        primitive.ObjectID `bson:"_id"`
	Domain    string             `bson:"domain"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Email     string             `bson:"email"`
}

type LeakEmailRelationEntry struct {
	EmailID     primitive.ObjectID `bson:"email_id"`
	LeakID      primitive.ObjectID `bson:"leak_id"`
	EmailDomain string             `bson:"email_domain"`
}

type MongoDBConn struct {
	conn *mongo.Database
}

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

func (db *MongoDBConn) GetEmailsByLeakID(leakId string) ([]*EmailEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "leak_id", Value: leakId}})

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

func (db *MongoDBConn) GetLeaksByEmailID(emailId string) ([]*LeakEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "email_id", Value: emailId}})

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

		leaks = append(leaks, &leakEntry)
	}

	return leaks, nil
}

func (db *MongoDBConn) GetEmailsByDomainAndLeakID(domain string, leakId string) ([]*EmailEntry, error) {
	if db == nil {
		return nil, errors.New("I am not even alive")
	}

	if db.conn == nil {
		return nil, errors.New("My connection is not set")
	}

	var relations []*LeakEmailRelationEntry

	ctx := context.TODO()
	cur, err := db.conn.Collection(_RelationCollectionName).Find(ctx, bson.D{primitive.E{Key: "$and", Value: bson.A{
		bson.D{primitive.E{Key: "email_domain", Value: domain}},
		bson.D{primitive.E{Key: "leak_id", Value: leakId}},
	}}})

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

func (db *MongoDBConn) GetEmailIDFromEmail(email string) (string, error) {
	var emailEnt EmailEntry

	err := db.conn.Collection(_EmailCollectionName).FindOne(context.TODO(), bson.D{primitive.E{Key: "email", Value: email}}).Decode(&emailEnt)

	if err != nil {
		return "", err
	}

	return emailEnt.ID.Hex(), nil
}
