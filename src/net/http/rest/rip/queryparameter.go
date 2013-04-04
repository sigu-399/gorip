// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	Kind            string
	DefaultValue    string
	FormatValidator validation.Validator
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
