package binding

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	"github.com/go-playground/validator/v10"
)

type StructValidator interface {
	ValidateStruct(any) error
	Engine() any
}

var Validator StructValidator = &defaultValidator{}

func validate(a any) error {
	if Validator == nil {
		return nil
	}
	return Validator.ValidateStruct(a)
}

type SliceValidationError []error

func (err SliceValidationError) Error() string {
	n := len(err)
	switch n {
	case 0:
		return ""
	default:
		var b strings.Builder
		if err[0] != nil {
			fmt.Fprintf(&b, "[%d]: %s", 0, err[0].Error())
		}
		if n > 1 {
			for i := 1; i < n; i++ {
				if err[i] != nil {
					b.WriteString("\n")
					fmt.Fprintf(&b, "[%d]: %s", i, err[i].Error())
				}
			}
		}
		return b.String()
	}
}

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

func (v *defaultValidator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")
	})
}

func (v *defaultValidator) validateStruct(a any) error {
	v.lazyInit()
	return v.validate.Struct(a)
}

func (v *defaultValidator) ValidateStruct(a any) error {
	if a == nil {
		return nil
	}

	value := reflect.ValueOf(a)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.validateStruct(a)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(SliceValidationError, 0)
		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}
		if len(validateRet) == 0 {
			return nil
		}
		return validateRet
	default:
		return nil
	}
}

func (v *defaultValidator) Engine() any {
	v.lazyInit()
	return v.validate
}
