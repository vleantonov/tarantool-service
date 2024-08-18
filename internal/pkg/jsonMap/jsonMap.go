package jsonMap

import "errors"

var (
	ErrConvertMap = errors.New("can't convert map to JsonMap")
)

// Obj is JSON map
// Alias for map[string]interface{}
type Obj map[string]interface{}

// convertMapToObj lets get JSON map from map[interface{}]interface{}
// It's useful when you want to get JSON encoding of m.
func ConvertMapToObj(m map[interface{}]interface{}) (Obj, error) {

	// Result of conversion
	converted := make(Obj)
	var err error

	// For each field there is a type check and conversion map[interface{}]interface{}
	for key, value := range m {
		strKey, ok := key.(string)
		if !ok {
			// Return error if it is not JSON map
			return nil, ErrConvertMap
		}

		// Type switching for recursive conversion
		switch v := value.(type) {
		case map[interface{}]interface{}:
			converted[strKey], err = ConvertMapToObj(v)
			if err != nil {
				return nil, err
			}
		case []interface{}:
			tempSl := make([]interface{}, len(v))

			// For each element
			// there is a type check and conversion map[interface{}]interface{}
			for idx, elem := range v {
				if _, ok := elem.(map[interface{}]interface{}); ok {
					tempSl[idx], err = ConvertMapToObj(elem.(map[interface{}]interface{}))
					if err != nil {
						return nil, err
					}
				} else {
					tempSl[idx] = elem
				}
			}
			converted[strKey] = tempSl
		default:
			// Other value remains untouched
			converted[strKey] = v
		}
	}

	return converted, nil
}
