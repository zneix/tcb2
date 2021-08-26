package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Connection struct {
	client       *mongo.Client
	ctx          context.Context
	databaseName string
}

//const databaseName = "tcb2"
type CollectionName string

const (
	CollectionNameChannels = CollectionName("channels")
	// TODO: Add handling for the (admin) users in the "users" collection
	//CollectionNameUsers    = CollectionName("users")
)
