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

type CollectionName string

const (
	CollectionNameChannels = CollectionName("channels")
	CollectionNameMOTDs    = CollectionName("motds")

	// Migration collections
	// CollectionNameChannels         = CollectionName("new-channels") // Temporary coll with new channels
	CollectionNameOldSubscriptions = CollectionName("old-subscriptions") // Coll where old subscriptions should be imported

	// TODO: Add handling for the (admin) users in the "users" collection
	// CollectionNameUsers    = CollectionName("users")
)
