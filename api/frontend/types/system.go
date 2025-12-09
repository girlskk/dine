package types

import "time"

type SystemNowResp struct {
	Now time.Time `json:"now"`
}
