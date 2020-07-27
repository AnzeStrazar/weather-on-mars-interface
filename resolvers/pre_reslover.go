package resolvers

import "weather-on-mars-interface/types"

type PreResolver struct {
	pre *types.PRE
}

func (p *PreResolver) Av() *float64 {
	return &p.pre.Av
}

func (p *PreResolver) Ct() *float64 {
	return &p.pre.Ct
}

func (p *PreResolver) Mn() *float64 {
	return &p.pre.Mn
}

func (p *PreResolver) Mx() *float64 {
	return &p.pre.Mx
}
