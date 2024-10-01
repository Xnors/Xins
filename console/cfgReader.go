package console

import (
	"encoding/json"
	"os"
)


func ReadConfig(filePath string) (map[string]interface{}, error) {
	file, err := os.Open(filePath)
	if err!= nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := make(map[string]interface{})
	err = decoder.Decode(&config)
	if err!= nil {
		return nil, err
	}

	return config, nil
}

