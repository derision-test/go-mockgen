package matchers

import (
	"reflect"
)

type callInstance interface {
	Args() []interface{}
	Results() []interface{}
}

var callInstanceInterface = reflect.TypeOf((*callInstance)(nil)).Elem()

func getCallHistory(v interface{}) ([]callInstance, bool) {
	value := reflect.ValueOf(v)
	if !value.IsValid() {
		return nil, false
	}

	method := value.MethodByName("History")
	if !method.IsValid() {
		return nil, false
	}

	if method.Type().NumIn() != 0 || method.Type().NumOut() != 1 {
		return nil, false
	}

	history := method.Call(nil)[0]
	if history.Kind() != reflect.Slice || !history.Type().Elem().Implements(callInstanceInterface) {
		return nil, false
	}

	calls := []callInstance{}
	for i := 0; i < history.Len(); i++ {
		calls = append(calls, history.Index(i).Interface().(callInstance))
	}

	return calls, true
}
