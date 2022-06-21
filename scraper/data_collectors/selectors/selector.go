package selectors

import "github.com/PuerkitoBio/goquery"

type Selector interface {
	Select(*goquery.Selection) *goquery.Selection
}

func NewSelector(config map[string]interface{}) (Selector, error) {
	selector := SelectorRegister[config["type"].(string)](config)
	return selector, nil
}

func NewSelectors(config []interface{}) ([]Selector, error) {
	selectors := []Selector{}
	for _, cnf := range config {
		sel, err := NewSelector(cnf.(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		selectors = append(selectors, sel)
	}
	return selectors, nil
}
