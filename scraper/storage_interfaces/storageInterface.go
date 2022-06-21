package storage_interfaces

import (
	"errors"
)

type StorageInterface interface {
	Open() error
	Close() error
	Save(map[string]interface{}) error
}

func NewStorageInterface(config map[string]interface{}) (StorageInterface, error) {
	if storageIntefaceFnc, exist := StorageInterfaceRegister[config["interface"].(string)]; exist {
		storageInterface, err := storageIntefaceFnc(config)
		if err != nil {
			log.Errorf("failed to initiate storage interface: '%s'", err.Error())
			return nil, errors.New("failed to initiate storage interface")
		}
		return storageInterface, nil
	} else {
		log.Errorf("did not find storage interface with name: '%s'", config["interface"].(string))
		return nil, errors.New("did not find storage interface with that name")
	}
}

func NewStorageInterfaces(config []interface{}) ([]StorageInterface, error) {
	var storageInterfaces []StorageInterface
	for _, cnf := range config {
		si, err := NewStorageInterface(cnf.(map[string]interface{}))
		if err != nil {
			log.Errorf("failed to initiate storage interfaces: '%s'", err)
			return nil, errors.New("failed to initiated storage interfaces")
		}
		storageInterfaces = append(storageInterfaces, si)
	}
	return storageInterfaces, nil
}
