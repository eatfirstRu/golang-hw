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
	ErrZeroValue      = errors.New("zero value")
	ErrMin            = errors.New("less than min")
	ErrMax            = errors.New("greater than max")
	ErrLen            = errors.New("invalid length")
	ErrRegexp         = errors.New("regular expression mismatch")
	ErrIn             = errors.New("values not in enumeration")
	ErrUnsupported    = errors.New("unsupported type")
	ErrBadParameter   = errors.New("bad parameter")
	ErrUnknownTag     = errors.New("unknown tag")
	ErrInvalid        = errors.New("invalid value")
	ErrCannotValidate = errors.New("cannot validate unexported struct")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	if len(v) == 1 {
		return fmt.Errorf("field: %s %w ", v[0].Field, v[0].Err).Error()
	}

	var sb strings.Builder
	for i, ve := range v {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Errorf("field: %s %w ", ve.Field, ve.Err).Error())
	}
	return sb.String()
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
	ve = ve[:0]

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tag := field.Tag.Get("validate")
		val := reflect.ValueOf(v).FieldByName(field.Name)
		if tag != "" {
			andCond := strings.Split(tag, "|")
			for _, pair := range andCond {
				parts := strings.Split(pair, ":")
				if len(parts) != 2 {
					ve = append(ve, ValidationError{Field: field.Name, Err: ErrUnknownTag})
					continue
				}
				if _, ok := funcMap[parts[0]]; ok {
					checkParam(field.Name, parts[0], parts[1], val)
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

func checkIn(fieldName, paramVal string, v reflect.Value) (bool, error) {
	sl := strings.Split(paramVal, ",")
	for _, vSl := range sl {
		if v.Kind() == reflect.Int {
			if pv, err := convToInt(vSl, fieldName); err == nil && v.Int() == int64(pv) {
				return true, ErrIn
			}
		}
		if v.Kind() == reflect.String && v.String() == vSl {
			return true, ErrIn
		}
	}
	return false, ErrIn
}

func checkRegexp(paramVal string, v reflect.Value) (bool, error) {
	re := regexp.MustCompile(paramVal)
	return re.MatchString(v.String()), ErrRegexp
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
	case paramName == "len" && v.Kind() == reflect.String:
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = utf8.RuneCountInString(v.String()) == pv
			retErr = ErrLen
		}
	case paramName == "min" && v.Kind() == reflect.Int:
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = v.Int() > int64(pv)
			retErr = ErrMin
		}
	case paramName == "max" && v.Kind() == reflect.Int:
		if pv, retErr = convToInt(paramVal, fieldName); retErr == nil {
			valid = v.Int() < int64(pv)
			retErr = ErrMax
		}
	case paramName == "in" && (v.Kind() == reflect.Int || v.Kind() == reflect.String):
		valid, retErr = checkIn(fieldName, paramVal, v)
	case paramName == "regexp":
		valid, retErr = checkRegexp(paramVal, v)
	default:
		appendVE(fieldName, ErrUnsupported)
	}
	if !valid && !errors.Is(retErr, ErrBadParameter) {
		appendVE(fieldName, retErr)
	}
}
