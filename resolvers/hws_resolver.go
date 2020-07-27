package resolvers

import "weather-on-mars-interface/types"

type HwsResolver struct {
	hws *types.HWS
}

func (h *HwsResolver) Av() *float64 {
	return &h.hws.Av
}

func (h *HwsResolver) Ct() *float64 {
	return &h.hws.Ct
}

func (h *HwsResolver) Mn() *float64 {
	return &h.hws.Mn
}

func (h *HwsResolver) Mx() *float64 {
	return &h.hws.Mx
}
