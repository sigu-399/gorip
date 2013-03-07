// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		Defines a REST endpoint.
//					An endpoint is a url route and associated resources.
// 
// created      	07-03-2013

package gorip

import ()

type endpoint struct {
	route     string
	resources []Resource
}

func NewEndpoint(route string) *endpoint {
	return &endpoint{route: route}
}

func (e *endpoint) GetRoute() string {
	return e.route
}

func (e *endpoint) AddResource(resource Resource) {
	e.resources = append(e.resources, resource)
}
