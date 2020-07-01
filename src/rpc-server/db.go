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
	return nil, nil
}

func (db *MongoDBConn) GetLeaksByEmail(email string) ([]*LeakEntry, error) {
	return nil, nil
}

func (db *MongoDBConn) GetLeaksByDomain(domain string) ([]*LeakEntry, error) {
	return nil, nil
}

func (db *MongoDBConn) GetEmailsByDomainAndLeakID(domain string, leakId string) ([]*EmailEntry, error) {
	return nil, nil
}
