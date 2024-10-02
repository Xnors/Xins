package console

import (
	"encoding/json"
	"os"
)

type MirrorType map[string]map[string]map[string]map[string]string

func ReadMirrors(filePath string) (MirrorType, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := make(MirrorType)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
