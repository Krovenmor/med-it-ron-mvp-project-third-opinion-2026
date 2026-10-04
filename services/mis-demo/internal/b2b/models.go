package b2b

import "encoding/json"

type Delivery struct {
	StudyID  string          `json:"study_id"`
	Status   int             `json:"status"`
	Response json.RawMessage `json:"response"`
}
