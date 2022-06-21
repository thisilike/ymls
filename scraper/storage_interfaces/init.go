package storage_interfaces

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var StorageInterfaceRegister map[string]func(
	config map[string]interface{},
) (StorageInterface, error)

func init() {
	log = config.Logger
	StorageInterfaceRegister = make(map[string]func(config map[string]interface{}) (StorageInterface, error))
	StorageInterfaceRegister["json"] = NewJsonInterface
}
