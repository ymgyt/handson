package main

import (
	"encoding"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/k0kubun/pp"
)

var ErrInvalidSpecification = errors.New("specification must be a struct pointer")

type varInfo struct {
	Name  string
	Alt   string
	Key   string
	Field reflect.Value
	Tags  reflect.StructTag
}

func gatherInfo(prefix string, spec interface{}) ([]varInfo, error) {
	expr := regexp.MustCompile("[^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")
	rv := reflect.ValueOf(spec)

	if rv.Kind() != reflect.Ptr {
		return nil, ErrInvalidSpecification
	}
	rv = rv.Elem()
	if rv.Kind() != reflect.Struct {
		return nil, ErrInvalidSpecification
	}
	rt := rv.Type()

	infos := make([]varInfo, 0, rv.NumField())
	_ = infos
	for i := 0; i < rv.NumField(); i++ {
		fv := rv.Field(i)
		ft := rt.Field(i)
		if !fv.CanSet() || isTrue(ft.Tag.Get("ignored")) {
			continue
		}

		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				if fv.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct
					break
				}
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			fv = fv.Elem()
		}

		info := varInfo{
			Name:  ft.Name,
			Field: fv,
			Tags:  ft.Tag,
			Alt:   strings.ToUpper(ft.Tag.Get("envconfig")),
		}

		info.Key = info.Name

		if isTrue(ft.Tag.Get("split_words")) {
			// UserName => [][]string{
			//	{"User", "User"},
			//  {"Name", "Name"}}
			words := expr.FindAllStringSubmatch(ft.Name, -1)
			if len(words) > 0 {
				var name []string
				for _, words := range words {
					name = append(name, words[0])
				}
				// User_Name
				info.Key = strings.Join(name, "_")
			}
		}
		// UserName `envconfig:"USER_NAME"` が優先される
		if info.Alt != "" {
			info.Key = info.Alt
		}
		if prefix != "" {
			info.Key = fmt.Sprintf("%s_%s", prefix, info.Key)
		}
		info.Key = strings.ToUpper(info.Key)
		infos = append(infos, info)

		if fv.Kind() == reflect.Struct {
			if decoderFrom(fv) == nil && setterFrom(fv) == nil && textUnmarshaler(fv) == nil {
			}

		}

	}
	return nil, nil
}

func isTrue(s string) bool {
	b, _ := strconv.ParseBool(s)
	return b
}

type Decoder interface {
	Decode(value string) error
}

type Setter interface {
	Set(value string) error
}

func interfaceFrom(field reflect.Value, fn func(interface{}, *bool)) {
	if !field.CanInterface() {
		return
	}
	var ok bool
	fn(field.Interface(), &ok)
	if !ok && field.CanAddr() {
		fn(field.Addr().Interface(), &ok)
	}
}

func decoderFrom(field reflect.Value) (d Decoder) {
	interfaceFrom(field, func(v interface{}, ok *bool) { d, *ok = v.(Decoder) })
	return d
}

func setterFrom(field reflect.Value) (s Setter) {
	interfaceFrom(field, func(v interface{}, ok *bool) { s, *ok = v.(Setter) })
	return s
}

func textUnmarshaler(field reflect.Value) (t encoding.TextUnmarshaler) {
	interfaceFrom(field, func(v interface{}, ok *bool) { t, *ok = v.(encoding.TextUnmarshaler) })
	return t
}

func TestRegex(t *testing.T) {
	expr := regexp.MustCompile("([^A-Z]+|[A-Z][^A-Z]+|[A-Z]+)")
	inputs := []string{"HooName"}
	for _, input := range inputs {
		words := expr.FindAllStringSubmatch(input, -1)
		pp.Println(words)
	}
}
