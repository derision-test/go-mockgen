package testutil

import "reflect"

// CallInstance holds the arguments and results of a single mock function call.
type CallInstance interface {
	Args() []interface{}
	Results() []interface{}
}

// GetCallHistory extracts the history from the given mock function and returns the
// set of call instances.
func GetCallHistory(v interface{}) ([]CallInstance, bool) {
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
	if history.Kind() != reflect.Slice || !history.Type().Elem().Implements(reflect.TypeOf((*CallInstance)(nil)).Elem()) {
		return nil, false
	}

	calls := []CallInstance{}
	for i := 0; i < history.Len(); i++ {
		calls = append(calls, history.Index(i).Interface().(CallInstance))
	}

	return calls, true
}

// GetCallHistoryWith extracts the history from the given mock function and returns the
// set of call instances that match the given function.
func GetCallHistoryWith(v interface{}, matcher func(v CallInstance) bool) (matching []CallInstance, _ bool) {
	history, ok := GetCallHistory(v)
	if !ok {
		return nil, false
	}

	for _, call := range history {
		if matcher(call) {
			matching = append(matching, call)
		}
	}

	return matching, true
}
