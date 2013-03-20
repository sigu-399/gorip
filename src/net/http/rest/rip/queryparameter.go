// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description	    Query parameters are variables given in the url.
// 
// created      	09-03-2013

package rip

import (
	"strconv"
	"strings/validation"
)

const (
	QueryParameterInt    = "int"
	QueryParameterFloat  = "float"
	QueryParameterString = "string"
	QueryParameterBool   = "bool"
)

type QueryParameter struct {
	Kind             string
	DefaultValue     string
	FormatValidation validation.Validator
}

func (q *QueryParameter) IsValidType(value string) bool {

	switch q.Kind {

	case QueryParameterInt:
		_, err := strconv.Atoi(value)
		return err == nil

	case QueryParameterFloat:
		_, err := strconv.ParseFloat(value, 64)
		return err == nil

	case QueryParameterString:
		return true

	case QueryParameterBool:
		return value == `true` || value == `false`
	}

	return false
}

func GetQueryParameterStringValue(value string) (bool, string) {
	return true, value
}

func GetQueryParameterIntValue(value string) (bool, int) {
	i, err := strconv.Atoi(value)
	if err != nil {
		return false, 0
	}
	return true, i
}

func GetQueryParameterFloatValue(value string) (bool, float64) {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false, 0
	}
	return true, f
}

func GetQueryParameterBoolValue(value string) (bool, bool) {
	if value == `true` {
		return true, true
	}
	if value == `false` {
		return true, false
	}
	return false, false
}
