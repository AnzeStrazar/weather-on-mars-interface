package resolvers

import (
	"weather-on-mars-interface/types"
)

type SolResolver struct {
	sol types.Sol
}

func (a *SolResolver) At() *AtResolver {
	atResult := a.sol.Data.At
	return &AtResolver{&atResult}
}
