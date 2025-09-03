package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type stringOrSlice []string
type strudelSampleMap map[string]stringOrSlice

func (s *stringOrSlice) UnmarshalJSON(data []byte) error {
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

func (s stringOrSlice) MarshalJSON() ([]byte, error) {
	if len(s) == 1 {
		return json.Marshal(s[0])
	}
	return json.Marshal([]string(s))
}

func addToStrudelSampleMap(sampleMap strudelSampleMap, path string) (strudelSampleMap, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var loadSampleMap strudelSampleMap
	err = json.Unmarshal(raw, &loadSampleMap)
	if err != nil {
		return nil, err
	}

	for k, v := range loadSampleMap {
		if _, exists := sampleMap[k]; exists {
			return nil, fmt.Errorf("duplicate key [%s] found in [%s]", k, path)
		}

		sampleMap[k] = v
	}

	return sampleMap, nil
}
