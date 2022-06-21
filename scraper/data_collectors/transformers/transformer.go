package transformers

import "errors"

// Transform get's an error from the previous transformer and decicdes wether to continue or skip a pass the error through
type Transformer interface {
	Transform(string, error) (string, error)
}

func NewTransformer(t string) (Transformer, error) {
	if transformer, ok := TransformerRegister[t]; ok {
		return transformer(), nil
	} else {
		log.Errorf("invalid transformer name: '%s'", t)
		return nil, errors.New("invalid transformer name")
	}
}

func NewTransformers(config []interface{}) ([]Transformer, error) {
	transformerList := []Transformer{}
	for _, cnfName := range config {
		if name, ok := cnfName.(string); !ok {
			log.Errorf("invalid transformer: '%s' is not a string", cnfName)
			return nil, errors.New("invalid transformer")
		} else {
			transformerList = append(transformerList, TransformerRegister[name]())
		}
	}
	return transformerList, nil
}
