package resolvers

import (
	"weather-on-mars-interface/types"
)

type SolResolver struct {
	sol types.Sol
}

func (a *SolResolver) At() *AtResolver {
	atResult := a.sol.At
	return &AtResolver{&atResult}
}

func (a *SolResolver) Hws() *HwsResolver {
	hwsResult := a.sol.Hws
	return &HwsResolver{&hwsResult}
}

func (a *SolResolver) Pre() *PreResolver {
	preResult := a.sol.Pre
	return &PreResolver{&preResult}
}

func (a *SolResolver) SolID() string {
	solIDResult := a.sol.SolId
	return solIDResult
}
