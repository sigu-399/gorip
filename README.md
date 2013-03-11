# gorip

( Go REST In Peace )

REST Server Framework written in Go language

## Status
Work in progress ( 90% done )

## Usage
```
package main

import (
	"bytes"
	"fmt"
	"net/http/rest/rip"
	"net/http"
	"strconv"
)

// Defines a route variable validator.
// This example one defines the type 'id', wich is a >0 integer.
// You may define any kind of validator depending on your needs 
type ResourceIdValidator struct {
}

func (v *ResourceIdValidator) Matches(value string) bool {
	id, err := strconv.Atoi(value)
	if err != nil {
		return false
	}
	return id > 0
}

// Declares our resource called THING. Must match the rip.Resource interface.
type ResourceThingGET struct {
}

// Factory, must return an new instance of itself.
// So each new API call will create a new instance of this resource.
func (r *ResourceThingGET) Factory() rip.Resource {
	return &ResourceThingGET{}
}

// Wich method do you implement ? here, it is a GET
func (r *ResourceThingGET) GetMethod() string {
	return rip.HttpMethodGET
}

// Allowed Content-Type IN ( None since a GET does not have a body, so no content type)
func (r *ResourceThingGET) GetContentTypeIn() []string {
	return []string{}
}

// Allowed Content-Type OUT, this implementation returns text/plain, but could be JSON, XML, images...
func (r *ResourceThingGET) GetContentTypeOut() []string {
	return []string{`text/plain`}
}

func (r *ResourceThingGET) GetQueryParameters() map[string]rip.QueryParameter {
	return map[string]rip.QueryParameter{
		"who": rip.QueryParameter{Kind: rip.QueryParameterString, DefaultValue: "World"}}
}

// The implementation of the endpoint
func (r *ResourceThingGET) Execute(context *rip.ResourceContext) rip.ResourceResult {
	fmt.Printf("context %s\n", context)
	return rip.ResourceResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(fmt.Sprintf(`Hello %s !`, context.QueryParameters["who"]))}
}

func main() {

	var err error

	// Creates a server listening to everything on port 8080
	server := rip.NewServer("/", ":8080")

	// Registers the validator, so we can use it to register dynamic routes
	server.RegisterRouteVariableValidator("id", &ResourceIdValidator{})

	// Register an endpoint, ex /things/4, things/890 ...
	endpointThing := rip.NewEndpoint("/things/{thing_id:id}")
	endpointThing.AddResource(&ResourceThingGET{})
	// Here you could add more resource to this endpoint...

	err = server.RegisterEndpoint(endpointThing)
	if err != nil {
		panic(err.Error())
	}

	// Starts the server
	server.ListenAndServe()
}

```
