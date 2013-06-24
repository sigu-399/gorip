package myResourceHandlers

import (
	"bytes"
	"fmt"
	rip "gorip"
	"net/http"
)

// Definition of our User resource, handling a GET request
// Must match the interface in gorip.Resource
type GetUser struct {
}

// Here you have to return a new instance
func (r *GetUser) Factory() rip.Resource {
	return &GetUser{}
}

// Defines the signature of the resource
// Here:
// - This resource reacts to a GET method
// - It doesnt require a content-type in input ( GET methods do not send data )
// - It returns a response in text/plain format only ( but could be JSON, XML ... )
//
// Note you may accept multiple content-type as input and output
//
func (r *GetUser) Signature() *rip.ResourceSignature {
	return &rip.ResourceSignature{Method: rip.HttpMethodGET, ContentTypeIn: []string{}, ContentTypeOut: []string{`text/plain`}}
}

// Defines the Query Parameters to read from the URL
// Here we expect something like GET /users/4?who=XYZ&age=123
//
func (r *GetUser) QueryParameters() map[string]rip.QueryParameter {
	return map[string]rip.QueryParameter{
		"who": rip.QueryParameter{Kind: rip.QueryParameterString, DefaultValue: "World"},
		"age": rip.QueryParameter{Kind: rip.QueryParameterInt, DefaultValue: "18"}}
}

// You can add your own notes to be displayed in the documentation
func (r *GetUser) DocumentationNotes() string {
	return ""
}

// Implementation of the request handler
// Here we simply return and 200 OK response, with some text as the body
// Note the use of the context to read the query parameters and route variables
//
func (r *GetUser) Execute(context *rip.ResourceContext) rip.ResourceResult {
	displayText := fmt.Sprintf(`Hello %s, your id is %s !`, context.QueryParameters["who"], context.RouteVariables["user_id"])
	return rip.ResourceResult{HttpStatus: http.StatusOK, Body: bytes.NewBufferString(displayText)}
}
