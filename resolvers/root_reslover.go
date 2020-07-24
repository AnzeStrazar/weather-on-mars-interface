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

func (r *RootResolver) Sol(args struct{ ID string }) *SolResolver {
	r.cache.SolMutex.RLock()
	defer r.cache.SolMutex.RUnlock()
	solResult := r.cache.SolCache[args.ID]

	return &SolResolver{solResult}
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
