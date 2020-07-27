package types

type Sol struct {
	SolId string `json:"solID" bson:"solID"`
	At    AT     `json:"AT" bson:"AT"`
	Hws   HWS    `json:"HWS" bson:"HWS"`
	Pre   PRE    `json:"PRE" bson:"PRE"`
}

type AT struct {
	Av float64 `json:"av" bson:"av"`
	Ct float64 `json:"ct" bson:"ct"`
	Mn float64 `json:"mn" bson:"mn"`
	Mx float64 `json:"mx" bson:"mx"`
}

type HWS struct {
	Av float64 `json:"av" bson:"av"`
	Ct float64 `json:"ct" bson:"ct"`
	Mn float64 `json:"mn" bson:"mn"`
	Mx float64 `json:"mx" bson:"mx"`
}

type PRE struct {
	Av float64 `json:"av" bson:"av"`
	Ct float64 `json:"ct" bson:"ct"`
	Mn float64 `json:"mn" bson:"mn"`
	Mx float64 `json:"mx" bson:"mx"`
}
