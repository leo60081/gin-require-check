package middleware

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func checkAndBuildSturct(part string, getValue func(string) string, _a interface{}) (err error) {
	v := reflect.Indirect(reflect.ValueOf(_a))
	count := v.NumField()
	errMsg := make([]string, 0)
	for i := 0; i < count; i++ {
		f := v.Field(i)
		n := v.Type().Field(i)

		field := n.Name
		if n.Tag.Get(part) != "" {
			l := strings.Split(n.Tag.Get(part), ",")
			field = l[0]
		}
		r := getValue(field)

		if r == "" && n.Tag.Get("binding") == "required" {
			errMsg = append(errMsg, fmt.Sprintf("'%s' is missing.", field))
			continue
		} else if r != "" {
			if reflect.TypeOf(r).Kind() != f.Kind() {
				if reflect.TypeOf(r).ConvertibleTo(f.Type()) {
					f.Set(reflect.ValueOf(r).Convert(f.Type()))
				} else {
					switch f.Kind() {
					case reflect.Int, reflect.Int32, reflect.Int64:
						i, err := strconv.ParseInt(r, 10, 0)
						if err != nil {
							errMsg = append(errMsg, fmt.Sprintf("'%s' invalid value: %s", field, err.Error()))
							continue
						}
						f.Set(reflect.ValueOf(i).Convert(f.Type()))
					case reflect.Float32, reflect.Float64:
						i, err := strconv.ParseFloat(r, 0)
						if err != nil {
							errMsg = append(errMsg, fmt.Sprintf("'%s' invalid value: %s", field, err.Error()))
							continue
						}
						f.Set(reflect.ValueOf(i).Convert(f.Type()))
					}
				}
			} else {
				f.Set(reflect.ValueOf(r))
			}
		}
	}
	if len(errMsg) > 0 {
		return fmt.Errorf("%s %s", part, strings.Join(errMsg, "\n"))
	}
	return nil
}

func newStruct(i interface{}) (s interface{}) {
	return reflect.New(
		reflect.Indirect(
			reflect.ValueOf(i),
		).Type(),
	).Interface()
}
