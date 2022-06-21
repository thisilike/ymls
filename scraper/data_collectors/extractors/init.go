package extractors

import (
	"github.com/sirupsen/logrus"
	"github.com/thisilike/ymls/config"
)

var log *logrus.Logger
var ExtractorRegister map[string]func(map[string]interface{}) Extractor

func init() {
	log = config.Logger
	ExtractorRegister = make(map[string]func(map[string]interface{}) Extractor)
	ExtractorRegister["text"] = NewTextExtractor
}
