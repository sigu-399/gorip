package myRouteVariableTypes

import (
	"strconv"
)

// "id" is a route variable
// Its definition is simple : it must be an integer > 0

type IdType struct {
}

// IdType defintion should match the gorip.RouteVariableType interface
// Simple though, it is only a function definition:
//
func (v *IdType) Matches(value string) bool {
	id, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return id > 0
}
