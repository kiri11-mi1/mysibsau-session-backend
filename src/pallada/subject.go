package pallada

import (
	"encoding/json"
	"fmt"
)

type SubjectArray []interface{}

func (sa SubjectArray) MarshalJSON() ([]byte, error) {
	if len(sa) < 2 {
		return nil, fmt.Errorf("cannot marshal subject array value %v into a string", sa)
	}
	return json.Marshal(sa[1])
}
