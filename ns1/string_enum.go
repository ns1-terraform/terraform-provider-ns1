package ns1

import (
	"fmt"
	"strings"
)

// StringEnum stores a list of strings as keys in a Map, and as a quoted join.
// Has methods for conveniently testing for the presence of a string.
type StringEnum struct {
	ValueMap  map[string]int
	Expecting string
}

// NewStringEnum makes a StringEnum from a list of strings.
func NewStringEnum(values []string) *StringEnum {
	valueMap := make(map[string]int)
	quoted := make([]string, len(values), len(values))
	for i, value := range values {
		_, present := valueMap[value]
		if present {
			panic(fmt.Sprintf("duplicate value %q", value))
		}
		valueMap[value] = i

		quoted[i] = fmt.Sprintf("%q", value)
	}

	return &StringEnum{
		ValueMap:  valueMap,
		Expecting: strings.Join(quoted, ", "),
	}
}

// Check if string is in our list.
func (se *StringEnum) Check(v string) (int, error) {
	i, present := se.ValueMap[v]
	if present {
		return i, nil
	}
	return -1, fmt.Errorf("expecting one of %s; got %q", se.Expecting, v)
}

// ValidateFunc (schema helper)
func (se *StringEnum) ValidateFunc(v interface{}, k string) (ws []string, es []error) {
	_, err := se.Check(v.(string))
	if err != nil {
		return nil, []error{err}
	}
	return nil, nil
}
