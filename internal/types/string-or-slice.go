package types

import (
	"encoding/json"
	"fmt"
)

type StringOrSlice []string

func (s *StringOrSlice) UnmarshalJSON(data []byte) error {
	var single string
	if err := json.Unmarshal(data, &single); err == nil {
		// It's a string, wrap it in a slice
		*s = []string{single}
		return nil
	}

	var slice []string
	if err := json.Unmarshal(data, &slice); err == nil {
		*s = slice
		return nil
	}

	return fmt.Errorf("value is not a string or []string: %s", string(data))
}

func (s StringOrSlice) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}

	return json.Marshal([]string(s))
}
