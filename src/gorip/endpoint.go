// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		Defines a REST endpoint.
// 
// created      	07-03-2013

package gorip

import ()

type Endpoint struct {
	route string
}

func NewEndpoint(route string) *Endpoint {
	return &Endpoint{route: route}
}
