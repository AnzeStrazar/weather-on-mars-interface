package resolvers

import "weather-on-mars-interface/types"

type AtResolver struct {
	at *types.AT
}

func (a *AtResolver) Av() *float64 {
	return &a.at.Av
}

func (a *AtResolver) Ct() *float64 {
	return &a.at.Ct
}

func (a *AtResolver) Mn() *float64 {
	return &a.at.Mn
}

func (a *AtResolver) Mx() *float64 {
	return &a.at.Mx
}
