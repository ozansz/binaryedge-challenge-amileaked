package main

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LeakEntry struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}
type DomainEntry struct {
	ID   primitive.ObjectID `bson:"_id"`
	Name string             `bson:"name"`
}
type EmailEntry struct {
	ID        primitive.ObjectID `bson:"_id"`
	DomainID  primitive.ObjectID `bson:"domain_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
	Email     string             `bson:"email"`
}

type LeakEmailRelationEntry struct {
	EmailID primitive.ObjectID `bson:"email_id"`
	LeakID  primitive.ObjectID `bson:"eak_id"`
}

type MongoDBConn struct {
	conn *mongo.Database
}

func (db *MongoDBConn) Connect(uri string) error {
	return nil
}

func (db *MongoDBConn) GetLeaksByEmail() (*Leak, error) {
	return nil, nil
}

func (db *MongoDBConn) GetLeaksByDomainName() (*Leak, error) {
	return nil, nil
}
