package types

type Sol struct {
	SolId string
	Data  struct {
		At  AT  `json:"AT"`
		Hws HWS `json:"HWS"`
		Pre PRE `json:"PRE"`
	}
}

type AT struct {
	Av float64 `json:"av"`
	Ct float64 `json:"ct"`
	Mn float64 `json:"mn"`
	Mx float64 `json:"mx"`
}

type HWS struct {
	Av float64 `json:"av"`
	Ct float64 `json:"ct"`
	Mn float64 `json:"mn"`
	Mx float64 `json:"mx"`
}

type PRE struct {
	Av float64 `json:"av"`
	Ct float64 `json:"ct"`
	Mn float64 `json:"mn"`
	Mx float64 `json:"mx"`
}
