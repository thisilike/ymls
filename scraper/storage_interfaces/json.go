package storage_interfaces

import (
	"encoding/json"
	"errors"
	"fmt"
)

type JsonInterface struct {
}

func (jsonInterface JsonInterface) Save(data map[string]interface{}) error {
	d, err := json.Marshal(data)
	if err != nil {
		log.Errorf("failed to marshal json: '%s'", err)
		return errors.New("failed to marshal data")
	}
	fmt.Println(string(d))
	return nil
}

func (jsonInterface JsonInterface) Open() error {
	return nil
}

func (jsonInterface JsonInterface) Close() error {
	return nil
}

func NewJsonInterface(config map[string]interface{}) (StorageInterface, error) {
	return JsonInterface{}, nil
}
