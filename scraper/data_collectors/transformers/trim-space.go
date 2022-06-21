package transformers

import "strings"

type TrimSpace struct {
}

func (transformer TrimSpace) Transform(in string, err error) (string, error) {
	return strings.TrimSpace(in), err
}

func NewTrimSpace() Transformer {
	return TrimSpace{}
}
