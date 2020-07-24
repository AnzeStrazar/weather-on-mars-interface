package resolvers

import (
	"weather-on-mars-interface/cache"
)

type RootResolver struct {
	cache *cache.Cache
}

func NewRootResolver(cache *cache.Cache) *RootResolver {
	return &RootResolver{cache: cache}
}

func (r *RootResolver) Sol(args struct{ ID *string }) *SolResolver {
	if args.ID == nil {
		solResult := r.cache.SolCache["583"]
		return &SolResolver{&solResult}
	}
	solResult := r.cache.SolCache[*args.ID]

	return &SolResolver{&solResult}
}
