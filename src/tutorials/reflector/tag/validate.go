package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

type Validator interface {
	Validate(interface{}) (bool, error)
}

type (
	DefaultValidator struct {
	}
	NumberValidator struct {
		Min int
		Max int
	}
	StringValidator struct {
		Min int
		Max int
	}
	EmailValidator struct {
		mailRuler *regexp.Regexp
	}
)

func (v DefaultValidator) Validate(interface{}) (bool, error) {
	return true, nil
}
func (n NumberValidator) Validate(val interface{}) (bool, error) {
	num := val.(int)
	if num < n.Min {
		return false, fmt.Errorf("should be greater than %v", n.Min)
	}
	if n.Max >= n.Min && num > n.Max {
		return false, fmt.Errorf("should be less than %v", n.Max)
	}
	return true, nil
}
func (s StringValidator) Validate(val interface{}) (bool, error) {
	l := len(val.(string))
	if l == 0 {
		return false, fmt.Errorf("cannot be blank")
	}
	if l < s.Min {
		return false, fmt.Errorf("should be at least %v chars long", s.Min)
	}
	if s.Max >= s.Min && l > s.Max {
		return false, fmt.Errorf("should be less than %v chars long", s.Max)
	}
	return true, nil
}
func (email EmailValidator) Validate(val interface{}) (bool, error) {
	if !email.mailRuler.MatchString(val.(string)) {
		return false, fmt.Errorf("is not a valid email address")
	}
	return true, nil
}

func NewEmailValidator() EmailValidator {
	return EmailValidator{
		mailRuler: regexp.MustCompile(`\A[\w+\-.]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+\z`),
	}
}

func validatorResolver(tag string) Validator {
	args := strings.Split(tag, ",")
	switch args[0] {
	case "number":
		validator := NumberValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "string":
		validator := StringValidator{}
		fmt.Sscanf(strings.Join(args[1:], ","), "min=%d,max=%d", &validator.Min, &validator.Max)
		return validator
	case "email":
		validator := NewEmailValidator()
		return validator
	}
	return DefaultValidator{}
}

func ValidateStruct(s interface{}) []error {
	errs := []error{}

	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		tag := v.Type().Field(i).Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}

		validator := validatorResolver(tag)
		valid, err := validator.Validate(v.Field(i).Interface())
		if !valid && err != nil {
			errs = append(errs, fmt.Errorf("%s %s", v.Type().Field(i).Name, err.Error()))
		}
	}

	return errs
}
