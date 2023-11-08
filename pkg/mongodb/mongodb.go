package mongodb

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

var (
	Database = New()
)

func New() *MongoDB {
	db := MongoDB{}
	var err error

	mongoOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	db.Client, err = mongo.NewClient(mongoOptions)
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(0)
		return nil
	}

	err = db.Client.Connect(db.Ctx)
	if err != nil {
		log.Error().Msg(err.Error())
		os.Exit(0)
		return nil
	}

	db.Database = db.Client.Database("ogusers-gg")
	db.Collections.Members = db.Database.Collection("members")
	db.Collections.Tags = db.Database.Collection("tags")

	return &db
}

func (c *MongoDB) Find(queryData bson.M, collection *mongo.Collection) (*mongo.Cursor, error) {
	cursor, err := collection.Find(context.Background(), queryData)
	if err != nil {
		return nil, err
	}

	return cursor, nil
}

func (c *MongoDB) FindOne(queryData bson.M, collection *mongo.Collection) (bson.M, error) {
	var reply bson.M

	err := collection.FindOne(context.Background(), queryData).Decode(&reply)
	if err != nil {
		return nil, err
	}

	delete(reply, "_id")
	return reply, nil
}

func (c *MongoDB) FindOneAndUpdate(queryData bson.M, update bson.M, collection *mongo.Collection) (bson.M, error) {
	var reply bson.M

	err := collection.FindOneAndUpdate(context.Background(), queryData, update).Decode(&reply)
	if err != nil {
		return nil, err
	}

	delete(reply, "_id")
	return reply, nil
}

func (c *MongoDB) InsertOne(queryData bson.M, collection *mongo.Collection) (*mongo.InsertOneResult, error) {
	reply, err := collection.InsertOne(context.Background(), queryData)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (c *MongoDB) DeleteOne(queryData bson.M, collection *mongo.Collection) (*mongo.DeleteResult, error) {
	reply, err := collection.DeleteOne(context.Background(), queryData)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
