// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		A resource is an implementation of a REST method.
// 
// created      	07-03-2013

package rip

import (
	"bytes"
)

type Resource interface {
	Factory() Resource
	Execute(context *ResourceContext) ResourceResult
	GetMethod() string
	GetContentTypeIn() []string
	GetContentTypeOut() []string
	GetQueryParameters() map[string]QueryParameter
}

type ResourceContext struct {
	RouteVariables  map[string]string
	QueryParameters map[string]string
	ContentTypeIn   *string
	ContentTypeOut  *string
	Body            *bytes.Buffer
}

type ResourceResult struct {
	HttpStatus int
	Body       *bytes.Buffer
}
