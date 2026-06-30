package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var (
	// ErrZeroValue is the error returned when variable has zero value
	// and nonzero or nonnil was specified
	ErrZeroValue = errors.New("zero value")
	// ErrMin is the error returned when variable is less than mininum
	// value specified
	ErrMin = errors.New("less than min")
	// ErrMax is the error returned when variable is more than
	// maximum specified
	ErrMax = errors.New("greater than max")
	// ErrLen is the error returned when length is not equal to
	// param specified
	ErrLen = errors.New("invalid length")
	// ErrRegexp is the error returned when the value does not
	// match the provided regular expression parameter
	ErrRegexp = errors.New("regular expression mismatch")
	ErrIn     = errors.New("values not in enumeration")
	// ErrUnsupported is the error error returned when a validation rule
	// is used with an unsupported variable type
	ErrUnsupported = errors.New("unsupported type")
	// ErrBadParameter is the error returned when an invalid parameter
	// is provided to a validation rule (e.g. a string where an int was
	// expected (max=foo,len=bar) or missing a parameter when one is required (len=))
	ErrBadParameter = errors.New("bad parameter")
	// ErrUnknownTag is the error returned when an unknown tag is found
	ErrUnknownTag = errors.New("unknown tag")
	// ErrInvalid is the error returned when variable is invalid
	// (normally a nil pointer)
	ErrInvalid = errors.New("invalid value")
	// ErrCannotValidate is the error returned when a struct is unexported
	ErrCannotValidate = errors.New("cannot validate unexported struct")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	var err error
	if len(ve) == 0 {
		return ""
	}
	if len(v) == 1 {
		return fmt.Errorf("field: %s %w ", v[0].Field, v[0].Err).Error()
	}

	for _, ve := range v {
		err = errors.Join(fmt.Errorf("field: %s %w ", ve.Field, ve.Err), err)
	}
	return err.Error()
}

var ve = make(ValidationErrors, 0)

func Validate(v interface{}) error {
	funcMap := map[string]struct{}{
		"len":    {},
		"min":    {},
		"max":    {},
		"in":     {},
		"regexp": {},
	}
	st := reflect.TypeOf(v)
	clear(ve)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tag := field.Tag.Get("validate")
		val := reflect.ValueOf(v).FieldByName(field.Name) //.Interface()
		if tag != "" {
			vldt := make([]map[string]string, 0)
			andCond := strings.Split(tag, "|")
			for _, pair := range andCond {
				mp := make(map[string]string)
				parts := strings.Split(pair, ":")
				if len(parts) != 2 {
					// err := fmt.Errorf("%w: %s", ErrUnknownTag, parts)
					ve = append(ve, ValidationError{Field: field.Name, Err: ErrUnknownTag}) // err})
					continue
				}
				mp[parts[0]] = parts[1]
				vldt = append(vldt, mp)
				if _, ok := funcMap[parts[0]]; ok {
					checkParam(field.Name, parts[0], parts[1], val)
				} else {
					continue
				}
			}
		}
	}
	if len(ve) == 0 || ve[0].Err == nil {
		return nil
	}
	return ve
}

func appendVE(fieldName string, err error) {
	if len(ve) != 0 && ve[0].Err == nil {
		ve[0] = ValidationError{Field: fieldName, Err: err}
	} else {
		ve = append(ve, ValidationError{Field: fieldName, Err: err})
	}
}

func convToInt(paramVal, fieldName string) (int, error) {
	pv, err := strconv.Atoi(paramVal)
	if err != nil {
		appendVE(fieldName, ErrBadParameter)
		return 0, ErrBadParameter
	}
	return pv, nil
}

func checkParam(fieldName, paramName, paramVal string, v reflect.Value) {
	var valid bool
	var retErr error
	var pv int

	switch {
	case v.Kind() == reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			checkParam(fieldName, paramName, paramVal, v.Index(i))
		}
		valid = true
	case (paramName == "len" && v.Kind() == reflect.String):
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = utf8.RuneCountInString(v.String()) == pv
			retErr = ErrLen
		}
	case (paramName == "min" && v.Kind() == reflect.Int):
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = v.Int() > int64(pv)
			retErr = ErrMin
		}
	case (paramName == "max" && v.Kind() == reflect.Int):
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = v.Int() < int64(pv)
			retErr = ErrMax
		}
	case (paramName == "in" && (v.Kind() == reflect.Int || v.Kind() == reflect.String)):
		sl := strings.Split(paramVal, ",")
		retErr = ErrIn

		for _, vSl := range sl {
			if v.Kind() == reflect.Int {
				if pv, err := convToInt(vSl, fieldName); err == nil {
					valid = v.Int() == int64(pv)
				}
			}
			if v.Kind() == reflect.String {
				valid = v.String() == vSl
			}
			if valid {
				break
			}
		}
	case paramName == "regexp":
		retErr = ErrRegexp
		re := *regexp.MustCompile(paramVal)
		valid = re.MatchString(v.String())
	default:
		appendVE(fieldName, ErrUnsupported)
	}
	if !valid && retErr != ErrBadParameter {
		appendVE(fieldName, retErr)
	}
}
