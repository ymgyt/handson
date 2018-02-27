package reflect

import (
	"fmt"
	"math"
	"reflect"
	"strings"
)

type Validator interface {
	Validate() error
}

type PaymentTransaction struct {
	Amount      float64 `validate:"positive"`
	Description string  `validate:"max_length:250"`
}

func (p *PaymentTransaction) Validate() error {
	fmt.Println("Validating payment transaction")
	return nil
}

func Validate(obj interface{}) error {
	v := reflect.ValueOf(obj).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		vField := v.Field(i)
		tField := t.Field(i)

		tag := tField.Tag.Get("validate")
		if tag == "" {
			continue
		}

		switch vField.Kind() {
		case reflect.Float64:
			value := vField.Float()
			if tag == "positive" && value < 0 {
				value = math.Abs(value)
				vField.SetFloat(value)
			}
		case reflect.String:
			value := vField.String()
			if tag == "upper_case" {
				value = strings.ToUpper(value)
				vField.SetString(value)
			}
		default:
			return fmt.Errorf("unsupported kind %q", vField.Kind())
		}
	}

	return nil
}

func CustomValidate(obj interface{}) error {
	v := reflect.ValueOf(obj)
	t := v.Type()

	interfaceT := reflect.TypeOf((*Validator)(nil)).Elem()
	if !t.Implements(interfaceT) {
		return fmt.Errorf("Validator interface is not implemented")
	}

	validateFunc := v.MethodByName("Validate")
	validateFunc.Call(nil)
	return nil
}
