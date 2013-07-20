package main

import (
	"bytes"
	"fmt"
	"github.com/sigu-399/gorip"
	"net/http"
	"strconv"
)

// This example server handles a GET request on URLs in the form "/users/{user_id:id}"
//
// "{user_id:id}" is called a route variable, it is composed of a name(user_id) and a type(id).
//
// We will first define the type "id" to be recognized by the server.
// It is expected that the "id" is an integer > 0
//
// RouteVariableIdType fits the interface gorip.RouteVariableType ( defined in router.go )
//
// type RouteVariableType interface {
//   Matches(string) bool
// }
//
// You may, of course, to add your own types to fit your needs.
//
type RouteVariableIdType struct {
}

func (_ RouteVariableIdType) Matches(value string) bool {
	id, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return id > 0
}

// Here is our implementation of the GET request on that specific endpoint.
// The implementation fits the interface defined in resourceHandler.go:
//
// type ResourceHandlerImplementation interface {
//	Execute(context *ResourceHandlerContext) ResourceHandlerResult
// }
//
// rip provides a context, containing all informations needed to process the request:
// - the request's body
// - route variables
// - url parameters ( aka query parameters )
//
// and expects a result:
// - http status
// - response body
//
type GetUserResourceHandlerImpl struct {
}

func (_ GetUserResourceHandlerImpl) Execute(context *gorip.ResourceHandlerContext) gorip.ResourceHandlerResult {

	displayText := fmt.Sprintf(`Foo is %s and the user's id is %s`, context.QueryParameters["foo"], context.RouteVariables["user_id"])

	return gorip.ResourceHandlerResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(displayText)}

}

// Utility function to handle errors
func onError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func main() {

	// Setup a server listening to everything on port 8080
	myServer := gorip.NewServer("/", ":8080")

	// Registers our route variable
	err := myServer.NewRouteVariableType("id", RouteVariableIdType{})
	onError(err)

	// Registers our endpoint
	err = myServer.NewEndpoint("/users/{user_id:id}", gorip.ResourceHandler{
		Method:         "GET",                  // GET method
		ContentTypeIn:  []string{},             // No content in
		ContentTypeOut: []string{`text/plain`}, // Content out is plain text in this example, could be json or xml...
		QueryParameters: map[string]gorip.QueryParameter{ // Query parameters from the URL
			"foo": gorip.QueryParameter{Kind: gorip.QueryParameterString, DefaultValue: "George"},
			"bar": gorip.QueryParameter{Kind: gorip.QueryParameterInt, DefaultValue: "1984"}},
		Implementation: GetUserResourceHandlerImpl{}, // Implementation of this handler
		Documentation: &gorip.ResourceHandlerDocumentation{ // Documentation and live test
			TestURL:         "http://localhost:8080/users/1",
			TestContentType: "text/plain",
			AdditionalNotes: "This is a user..."}})
	onError(err)

	// Note: NewEndpoint is a variadic function, allowing multiple handlers to be added
	// For example: myServer.NewEndpoint(	"/users/{user_id:id}",
	// 										gorip.ResourceHandler{...},
	// 										gorip.ResourceHandler{...},
	// 										gorip.ResourceHandler{...} etc... )
	// Letting you bind a GET, POST and DELETE request to the same endpoint for example.

	// Miscellaneous:
	//
	// can auto generate documentation if you want to
	myServer.EnableDocumentationEndpoint("/documentation")

	// Logs an ascii representation of the routing tree
	myServer.DebugPrintRouterTree()

	// Dumps the request in json format to the console
	myServer.DebugEnableLogRequestDump(true)

	// Every log line will be annotated with a unique identifier to track a specific request all the way down
	myServer.DebugEnableLogRequestIdentifier(true)

	// Displays the duration of the request/response
	myServer.DebugEnableLogRequestDuration(true)

	// Finally, starts the server
	err = myServer.ListenAndServe()
	onError(err)

}
