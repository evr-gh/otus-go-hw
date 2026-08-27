package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrLenValidation    = errors.New("value does not match required length")
	ErrRegexpValidation = errors.New("value does not match regexp rule")
	ErrEnumValidation   = errors.New("value is not in set")
	ErrMinValidation    = errors.New("value is smaller than min rule")
	ErrMaxValidation    = errors.New("value is larger than max rule")

	ErrUnsupportedType  = errors.New("unsupported type, expected struct or poiter to struct")
	ErrUnsupportedField = errors.New("unsupported field type for validate")
	ErrInvalidRule      = errors.New("invalid rule")
	ErrInvalidRuleValue = errors.New("invalid rule value")
)

type ValidationError struct {
	Field string
	Err   error
}

func (ve ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", ve.Field, ve.Err)
}

func (ve ValidationError) Unwrap() error {
	return ve.Err
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}

	var b strings.Builder
	b.WriteString("Validation errors:\n")

	for _, item := range v {
		b.WriteString(item.Error())
		b.WriteString("\n")
	}
	return b.String()
}

func Validate(v interface{}) error {
	value := reflect.ValueOf(v)
	var t reflect.Type

	switch {
	case value.Kind() == reflect.Pointer:
		t = value.Type().Elem()
		if t.Kind() == reflect.Struct {
			value = reflect.Indirect(value)
		} else {
			return ErrUnsupportedType
		}
	case value.Kind() == reflect.Struct:
		t = reflect.TypeOf(v)
	default:
		return ErrUnsupportedType
	}

	errorList := make(ValidationErrors, 0, value.NumField())
	var ve *ValidationErrors
	err := validateFields(value, t, &errorList)
	if err != nil && !errors.As(err, &ve) {
		return err
	}

	if len(errorList) > 0 {
		return errorList
	}

	return nil
}

func validateFields(value reflect.Value, t reflect.Type, errorList *ValidationErrors) error {
	for i := 0; i < value.NumField(); i++ {
		field := value.Field(i)
		fieldType := t.Field(i)

		var ve *ValidationErrors
		switch {
		case field.Kind() == reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				fieldValue := field.Index(j)
				err := validateField(fieldType, fieldValue, errorList)
				if err != nil && !errors.As(err, &ve) {
					return err
				}
			}
		default:
			err := validateField(fieldType, field, errorList)
			if err != nil && !errors.As(err, &ve) {
				return err
			}
		}
	}

	return errorList
}

func validateField(fieldType reflect.StructField, fieldValue reflect.Value, errorList *ValidationErrors) error {
	validateRules := getFieldValidationRules(fieldType)

	for _, rule := range validateRules {
		if err := validate(rule, fieldType.Name, fieldValue); err != nil {
			if errors.As(err, &ValidationError{}) {
				*errorList = append(*errorList, ValidationError{Field: fieldType.Name, Err: errors.Unwrap(err)})
			} else {
				return err
			}
		}
	}

	return errorList
}

func getFieldValidationRules(field reflect.StructField) []string {
	validateTag, ok := field.Tag.Lookup("validate")
	if !ok {
		return []string{}
	}
	return strings.Split(validateTag, "|")
}

func validate(rule, fieldName string, value reflect.Value) error {
	rules := strings.Split(rule, ":")
	if len(rules) != 2 {
		return ErrInvalidRule
	}

	ruleName := rules[0]
	ruleValue := rules[1]

	if ruleValue == "" {
		return ErrInvalidRuleValue
	}

	switch {
	case value.Kind() == reflect.String:
		switch ruleName {
		case "len":
			return validateStrAgainstLen(ruleValue, fieldName, value)
		case "regexp":
			return validateStrAgainstRegexp(ruleValue, fieldName, value)
		case "in":
			return validateValAgainstEnum(ruleValue, fieldName, value)
		default:
			return ErrInvalidRule
		}
	case value.Kind() == reflect.Int:
		switch ruleName {
		case "in":
			return validateValAgainstEnum(ruleValue, fieldName, value)
		case "min":
			return validateIntAgainstMin(ruleValue, fieldName, value)
		case "max":
			return validateIntAgainstMax(ruleValue, fieldName, value)
		default:
			return ErrInvalidRule
		}
	default:
		return ErrUnsupportedField
	}
}

func validateStrAgainstLen(ruleValue, fieldName string, value reflect.Value) error {
	length, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if value.Len() != length {
		return ValidationError{fieldName, ErrLenValidation}
	}

	return nil
}

func validateStrAgainstRegexp(ruleValue, fieldName string, value reflect.Value) error {
	matched, err := regexp.MatchString(ruleValue, value.String())
	if err != nil {
		return ErrInvalidRuleValue
	}
	if !matched {
		return ValidationError{fieldName, ErrRegexpValidation}
	}

	return nil
}

func validateValAgainstEnum(ruleValue, fieldName string, value reflect.Value) error {
	validValues := strings.Split(ruleValue, ",")
	if value.Kind() == reflect.String {
		if !isInSlice(value.String(), validValues) {
			return ValidationError{fieldName, ErrEnumValidation}
		}
	} else if value.Kind() == reflect.Int {
		intVal := strconv.Itoa(int(value.Int()))
		if !isInSlice(intVal, validValues) {
			return ValidationError{fieldName, ErrEnumValidation}
		}
	}

	return nil
}

func validateIntAgainstMin(ruleValue, fieldName string, value reflect.Value) error {
	minVal, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if int(value.Int()) < minVal {
		return ValidationError{fieldName, ErrMinValidation}
	}

	return nil
}

func validateIntAgainstMax(ruleValue, fieldName string, value reflect.Value) error {
	maxVal, err := strconv.Atoi(ruleValue)
	if err != nil {
		return ErrInvalidRuleValue
	}
	if int(value.Int()) > maxVal {
		return ValidationError{fieldName, ErrMaxValidation}
	}

	return nil
}

func isInSlice(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
