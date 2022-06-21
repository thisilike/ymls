package selectors

import "github.com/PuerkitoBio/goquery"

type CssSelector struct {
	CssSelectorString string
}

func (cssSel CssSelector) Select(orSel *goquery.Selection) *goquery.Selection {
	return orSel.Find(cssSel.CssSelectorString)
}

func NewCssSelector(config map[string]interface{}) Selector {
	sel := CssSelector{
		CssSelectorString: config["cssSelectString"].(string),
	}
	return sel
}
