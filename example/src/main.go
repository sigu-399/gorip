package main

// A basic RIP server structure would be:
// - A package containing our resources ( User, Post, whatever... )
// - A package containing route variable types ( Id, RelationType, etc... )
import (
	"bytes"
	"fmt"
	rip "gorip"
	"net/http"
)
import (
	"strconv"
)

type GetUser struct {
}

func (_ GetUser) Execute(context *rip.ResourceHandlerContext) rip.ResourceHandlerResult {
	displayText := fmt.Sprintf(`Hello %s, your id is %s !`, context.QueryParameters["who"], context.RouteVariables["user_id"])
	return rip.ResourceHandlerResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(displayText)}
}

// "id" is a route variable
// Its definition is simple : it must be an integer > 0

type RouteVariableIdType struct {
}

// IdType defintion should match the gorip.RouteVariableType interface
// Simple though, it is only a function definition:
//
func (_ RouteVariableIdType) Matches(value string) bool {
	id, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return id > 0
}

func main() {

	var err error

	// Create a server listening to everything on port 8080.
	//
	myServer := rip.NewServer("/", ":8080")

	// Registers our id type, you are free to create any type you need/want	but this type is the most common and used one.
	// See endpoint creation for details in the next block of code.
	//
	// ( The route variable is defined in myRouteVariableTypes/idType.go )
	//
	err = myServer.NewRouteVariableType("id", RouteVariableIdType{})
	if err != nil {
		panic(err.Error())
	}

	// Creates a basic endpoint to access a particular User resource.
	// Note that the "URL" pattern uses the id type we just defined ( &myRouteVariablesTypes.IdType{} )
	// It will match URLs like /users/1 and /users/57689 ...
	//
	// &myResources.GetUser{} is an implementation of a Resource + Method.
	// Here for example this endpoint reacts to, for example a "GET /users/78"
	//
	// Note that you could add more resources to handle methods like POST, DELETE, etc...
	// example:
	// err = myServer.NewEndpoint("/users/{user_id:id}", &myResources.GetUser{}, &myResources.PostUser{})
	//
	// ( The Resource is defined in myResources/user.go )
	//
	err = myServer.NewEndpoint("/users/{user_id:id}", rip.ResourceHandler{Method: rip.HttpMethodGET,
		ContentTypeIn:  []string{},
		ContentTypeOut: []string{`text/plain`},
		QueryParameters: map[string]rip.QueryParameter{
			"who": rip.QueryParameter{Kind: rip.QueryParameterString, DefaultValue: "World"},
			"age": rip.QueryParameter{Kind: rip.QueryParameterInt, DefaultValue: "18"}},
		Implementation: GetUser{}})

	if err != nil {
		panic(err.Error())
	}

	// gorip has a built-in documentation - ideal during development
	// You can access it in your favourite browser at http://yourhost:8080/documentation
	myServer.EnableDocumentationEndpoint("/documentation")

	//The following functions can be usefull during development/debugging/tracking :

	// Displays an ascii representation of the router tree
	myServer.DebugPrintRouterTree()

	// Print a request dump to the console ( JSON format ) - on/off
	myServer.DebugEnableLogRequestDump(true)

	// Uses a unique identifier for every request log line - good for tracking - on/off
	myServer.DebugEnableLogRequestIdentifier(true)

	// Outputs the duration from the request to the response - on/off
	myServer.DebugEnableLogRequestDuration(true)

	// Finally starts the server
	err = myServer.ListenAndServe()
	if err != nil {
		panic(err.Error())
	}
}
