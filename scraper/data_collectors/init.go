package data_collectors

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger

func init() {
	log = config.Logger
}
