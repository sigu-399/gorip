package main

import (
	"bytes"
	"fmt"
	"github.com/sigu-399/gorip"
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

func (_ GetUserResourceHandlerImpl) Execute(context *gorip.ResourceHandlerContext) gorip.ResourceHandlerResult {
	displayText := fmt.Sprintf(`Hello %s, your id is %s !`, context.QueryParameters["who"], context.RouteVariables["user_id"])
	return gorip.ResourceHandlerResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(displayText)}
}

func main() {

	var err error

	myServer := gorip.NewServer("/", ":8080")

	err = myServer.NewRouteVariableType("id", RouteVariableIdType{})
	if err != nil {
		panic(err.Error())
	}

	err = myServer.NewEndpoint("/users/{user_id:id}", gorip.ResourceHandler{
		Method:         gorip.HttpMethodGET,
		ContentTypeIn:  []string{},
		ContentTypeOut: []string{`text/plain`},
		QueryParameters: map[string]gorip.QueryParameter{
			"who": gorip.QueryParameter{Kind: gorip.QueryParameterString, DefaultValue: "World"},
			"age": gorip.QueryParameter{Kind: gorip.QueryParameterInt, DefaultValue: "18"}},
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
