package main

// A basic RIP server structure would be:
// - A package containing our resources ( User, Post, whatever... )
// - A package containing route variable types ( Id, RelationType, etc... )
import (
	rip "github.com/sigu-399/gorip"
	"myResources"
	"myRouteVariablesTypes"
)

func main() {

	var err error

	// Create a server listening to everything on port 8080.
	//
	myServer := rip.NewServer("/", ":8080")

	// Registers our id type, you are free to create any type you need/want	but this type is the most common and used one.
	// See endpoint creation for details in the next block of code.
	//
	// ( The route variable is defined in myRouteVariablesTypes/idType.go )
	//
	err = myServer.NewRouteVariableType("id", &myRouteVariablesTypes.IdType{})
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
	err = myServer.NewEndpoint("/users/{user_id:id}", &myResources.GetUser{})
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
