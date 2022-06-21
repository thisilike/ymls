package data_collectors

import (
	"bytes"
	"errors"

	"github.com/PuerkitoBio/goquery"
	"github.com/thisilike/ymls/scraper/data_collectors/extractors"
	"github.com/thisilike/ymls/scraper/data_collectors/selectors"
	"github.com/thisilike/ymls/scraper/data_collectors/transformers"
)

type DataCollector struct {
	Name         string
	Selectors    []selectors.Selector
	Extractor    extractors.Extractor
	Transformers []transformers.Transformer
}

func (dataCollector DataCollector) Collect(data []byte) (interface{}, error) {
	if len(dataCollector.Name) == 0 {
		// save file
		log.Errorf("files support is still missing")
		return nil, errors.New("files are not supported yet")
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		log.Debug(string(data))
		log.Errorf("failed to read data to html: '%s'", err.Error())
		return nil, errors.New("invalid html")
	}
	sel := doc.Find("html")
	for _, selector := range dataCollector.Selectors {
		sel = selector.Select(sel)
	}
	extData, err := dataCollector.Extractor.Extract(sel)
	if err != nil {
		log.Errorf("failed to extract data: '%s'", err.Error())
		return nil, errors.New("failed to extract data")
	}
	for _, transformer := range dataCollector.Transformers {
		extData, err = transformer.Transform(extData, err)
	}
	return extData, err
}

func NewDataCollector(
	dataCollecterRegister map[string]bool,
	config map[string]interface{},
) (
	DataCollector,
	error,
) {
	dataCollector := DataCollector{}
	dataCollector.Name = config["name"].(string)
	if exists := dataCollecterRegister[dataCollector.Name]; exists {
		log.Errorf("duplicate collector name: '%s'", dataCollector.Name)
		return dataCollector, errors.New("duplicate collector name")
	}
	if _, ok := config["selectors"]; ok {
		tmp, err := selectors.NewSelectors(config["selectors"].([]interface{}))
		if err != nil {
			log.Errorf("failed to load selectors: '%s'", err)
			return dataCollector, errors.New("failed to load selectors")
		}
		dataCollector.Selectors = tmp
	}

	if _, ok := config["transformers"]; ok {
		tmp, err := transformers.NewTransformers(config["transformers"].([]interface{}))
		if err != nil {
			log.Errorf("failed to load transformers: '%s'", err)
			return dataCollector, errors.New("failed to load transformers")
		}
		dataCollector.Transformers = tmp
	}
	tmp, err := extractors.NewExtractor(config["extractor"].(map[string]interface{}))
	if err != nil {
		log.Errorf("failed to load extractor: '%s'", err)
		return dataCollector, errors.New("failed to load extractor")
	}
	dataCollector.Extractor = tmp

	return dataCollector, nil
}

func NewDataCollectors(
	dataCollecterRegister map[string]bool,
	config []interface{},
) (
	[]DataCollector,
	error,
) {
	collectors := []DataCollector{}
	for _, collector := range config {
		col, err := NewDataCollector(dataCollecterRegister, collector.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		collectors = append(collectors, col)
	}
	return collectors, nil
}
