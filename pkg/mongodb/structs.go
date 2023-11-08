package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDB struct {
	Ctx      context.Context
	Client   *mongo.Client
	Database *mongo.Database

	Collections struct {
		Members *mongo.Collection
		Tags    *mongo.Collection
	}
}
