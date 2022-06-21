package extractors

import "github.com/PuerkitoBio/goquery"

type TextExtractor struct {
}

func (textExtractor TextExtractor) Extract(sel *goquery.Selection) (string, error) {
	return sel.Text(), nil
}

func NewTextExtractor(map[string]interface{}) Extractor {
	return TextExtractor{}
}
