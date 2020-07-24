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

func (r *RootResolver) Sol(args struct{ ID *string }) *[]*SolResolver {
	resolversResults := make([]*SolResolver, 0)

	if args.ID == nil {
		for _, solResult := range r.cache.SolCache {
			resolversResults = append(resolversResults, &SolResolver{&solResult})
		}

		return &resolversResults
	}
	solResult := r.cache.SolCache[*args.ID]

	resolversResults = append(resolversResults, &SolResolver{&solResult})
	return &resolversResults
}
