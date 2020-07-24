package resolvers

import (
	"fmt"
	"weather-on-mars-interface/types"
)

type SolResolver struct {
	sol *types.Sol
}

func (a *SolResolver) At() *AtResolver {
	fmt.Println("TEST:::::")
	atResult := a.sol.Data.At
	return &AtResolver{&atResult}
}
