package selectors

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var SelectorRegister map[string]func(map[string]interface{}) Selector

func init() {
	log = config.Logger
	SelectorRegister = make(map[string]func(map[string]interface{}) Selector)
	SelectorRegister["css-select"] = NewCssSelector
}
