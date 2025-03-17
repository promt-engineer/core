package validator

import (
	"bitbucket.org/play-workspace/base-slot-server/pkg/kernel/constants"
	"fmt"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

type Validator struct {
	once sync.Once
	*validator.Validate
	config *constants.Config

	buildErr error
}

type TaggedError struct {
	Tag string
	Err error
}

func New(config *constants.Config) (*Validator, error) {
	v := &Validator{
		Validate: validator.New(),
		config:   config,
	}

	v.register(constants.GameRuleName, stringOneOfTheListOrNil(v.config.AvailableGames))

	if v.buildErr != nil {
		return nil, v.buildErr
	}

	v.generateErrorMessage()

	return v.Engine().(*Validator), nil
}

func (v *Validator) register(tag string, fn validator.Func) *Validator {
	if v.buildErr != nil {
		return v
	}

	v.buildErr = v.Validate.RegisterValidation(tag, fn)

	return v
}

func (v *Validator) lazyinit() {
	v.once.Do(func() {
		v.SetTagName("validate")
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

			if name == "-" {
				return ""
			}

			return name
		})
	})
}

func (v *Validator) Engine() interface{} {
	v.lazyinit()

	return v
}

func CheckValidationErrors(err error) (e []TaggedError) {
	if _, ok := err.(*validator.InvalidValidationError); ok {
		e = append(e, TaggedError{Tag: InvalidTag, Err: err})
	}

	errs, ok := err.(validator.ValidationErrors)

	if !ok {
		e = append(e, TaggedError{Tag: InvalidTag, Err: err})

		return e
	}

	for _, validationError := range errs {
		message, ok := errorMessages[validationError.Tag()]

		if !ok {
			e = append(e, TaggedError{Tag: validationError.Tag(), Err: validationError})

			continue
		}

		switch strings.Count(message, "%s") {
		case 0:
			e = append(e, TaggedError{Tag: validationError.Tag(), Err: fmt.Errorf(message)})
		case 1:
			e = append(e, TaggedError{
				Tag: validationError.Tag(),
				Err: fmt.Errorf(message, validationError.Field()),
			})
		case 2:
			e = append(e, TaggedError{
				Tag: validationError.Tag(),
				Err: fmt.Errorf(message, validationError.Field(), validationError.Param()),
			})
		}
	}

	return e
}

func (v *Validator) ValidateStruct(obj interface{}) error {
	if obj == nil {
		return nil
	}

	value := reflect.ValueOf(obj)
	switch value.Kind() {
	case reflect.Ptr:
		return v.ValidateStruct(value.Elem().Interface())
	case reflect.Struct:
		return v.Struct(obj)
	case reflect.Slice, reflect.Array:
		count := value.Len()
		validateRet := make(sliceValidateError, 0)

		for i := 0; i < count; i++ {
			if err := v.ValidateStruct(value.Index(i).Interface()); err != nil {
				validateRet = append(validateRet, err)
			}
		}

		if len(validateRet) == 0 {
			return nil
		}

		return validateRet
	case reflect.Bool, reflect.Chan, reflect.Complex128, reflect.Complex64, reflect.Float32, reflect.Float64,
		reflect.Func, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.Interface,
		reflect.Invalid, reflect.Map, reflect.String, reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uint8, reflect.Uintptr, reflect.UnsafePointer:
		fallthrough
	default:
		return nil
	}
}
