package extractors

import (
	"errors"

	"github.com/PuerkitoBio/goquery"
)

type Extractor interface {
	Extract(*goquery.Selection) (string, error)
}

func NewExtractor(config map[string]interface{}) (Extractor, error) {
	if extrFnc, ok := ExtractorRegister[config["type"].(string)]; ok {
		return extrFnc(config), nil
	} else {
		log.Errorf("did not find extractor: '%s'", config["type"].(string))
		return nil, errors.New("did not find extractor")
	}
}
