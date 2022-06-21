package transformers

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var TransformerRegister map[string]func() Transformer

func init() {
	log = config.Logger
	TransformerRegister = make(map[string]func() Transformer)
	TransformerRegister["trim-space"] = NewTrimSpace
}
