package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/creativenucleus/strudel-sample-server/internal/types"
)

type samplepack struct {
	pathBase  string
	sampleMap map[string]types.StringOrSlice
}

// readToStrudelSamplePack reads a strudel.json file and returns a sampleMap
func readToStrudelSamplePack(strudelJSONFilepath string) (*samplepack, error) {
	var newSamplePack samplepack

	raw, err := os.ReadFile(strudelJSONFilepath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(raw, &newSamplePack.sampleMap)
	if err != nil {
		return nil, err
	}

	newSamplePack.pathBase = filepath.Dir(strudelJSONFilepath)

	return &newSamplePack, nil
}

func (s samplepack) toData(urlBase string) ([]byte, error) {
	s.sampleMap["_base"] = []string{urlBase}

	out, err := json.Marshal(s.sampleMap)
	if err != nil {
		return nil, err
	}

	return out, nil
}
