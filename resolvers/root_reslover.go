package resolvers

import (
	"context"
	"fmt"
	"log"
	"weather-on-mars-interface/cache"
	"weather-on-mars-interface/types"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type RootResolver struct {
	cache  *cache.Cache
	client *mongo.Client
}

func NewRootResolver(cache *cache.Cache, client *mongo.Client) *RootResolver {
	return &RootResolver{cache: cache, client: client}
}

func (r *RootResolver) Sol(args struct{ ID string }) *SolResolver {
	var result types.Sol

	// TODO: Go to the database only for data that are older than seven days.
	if args.ID < "585" {
		solsCollections := r.client.Database("weather-on-mars").Collection("sols")
		solsFind, err := solsCollections.Find(context.TODO(), bson.M{"solID": bson.M{"$eq": args.ID}})
		if err != nil {
			fmt.Printf("error %s\n", err)
		}

		if solsFind.Next(context.TODO()) {
			err := solsFind.Decode(&result)
			if err != nil {
				log.Print(err)
			}
		}
	} else {
		r.cache.SolMutex.RLock()
		defer r.cache.SolMutex.RUnlock()
		result = r.cache.SolCache[args.ID]
	}

	return &SolResolver{result}
}

func (r *RootResolver) Sols() *[]*SolResolver {
	resolversResults := make([]*SolResolver, 0)

	r.cache.SolMutex.RLock()
	for _, solResult := range r.cache.SolCache {
		resolversResults = append(resolversResults, &SolResolver{solResult})
	}

	defer r.cache.SolMutex.RUnlock()

	return &resolversResults
}
