package main

import (
	"bytes"
	"fmt"
	rip "gorip"
	"net/http"
)
import (
	"strconv"
)

type RouteVariableIdType struct {
}

func (_ RouteVariableIdType) Matches(value string) bool {
	id, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return id > 0
}

type GetUserResourceHandlerImpl struct {
}

func (_ GetUserResourceHandlerImpl) Execute(context *rip.ResourceHandlerContext) rip.ResourceHandlerResult {
	displayText := fmt.Sprintf(`Hello %s, your id is %s !`, context.QueryParameters["who"], context.RouteVariables["user_id"])
	return rip.ResourceHandlerResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(displayText)}
}

func main() {

	var err error

	myServer := rip.NewServer("/", ":8080")

	err = myServer.NewRouteVariableType("id", RouteVariableIdType{})
	if err != nil {
		panic(err.Error())
	}

	err = myServer.NewEndpoint("/users/{user_id:id}", rip.ResourceHandler{
		Method:         rip.HttpMethodGET,
		ContentTypeIn:  []string{},
		ContentTypeOut: []string{`text/plain`},
		QueryParameters: map[string]rip.QueryParameter{
			"who": rip.QueryParameter{Kind: rip.QueryParameterString, DefaultValue: "World"},
			"age": rip.QueryParameter{Kind: rip.QueryParameterInt, DefaultValue: "18"}},
		Implementation: GetUserResourceHandlerImpl{}})

	if err != nil {
		panic(err.Error())
	}

	myServer.EnableDocumentationEndpoint("/documentation")
	myServer.DebugPrintRouterTree()
	myServer.DebugEnableLogRequestDump(true)
	myServer.DebugEnableLogRequestIdentifier(true)
	myServer.DebugEnableLogRequestDuration(true)

	err = myServer.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
