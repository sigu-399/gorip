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

package rip

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

func (e *endpoint) GetResources() []Resource {
	return e.resources
}

func (e *endpoint) FindMatchingResource(method string, contentTypeParser *contentTypeHeaderParser, acceptParser *acceptHeaderParser) (Resource, *string, *string) {

	var resultContentTypeIn *string
	var resultContentTypeOut *string

	// Loop through accepted OUT content types, highest priority first
	for _, acceptElement := range acceptParser.contentTypes {
		// Find a resource for given method
		for _, v := range e.resources {
			if v.GetMethod() == method {
				allContentTypeOut := v.GetContentTypeOut()
				allContentTypeIn := v.GetContentTypeIn()

				// If OUT content type matches or 'matching everything' */* then the resource matches
				for _, contentTypeOut := range allContentTypeOut {
					if contentTypeOut == acceptElement.contentType || acceptElement.contentType == `*/*` {

						resultContentTypeOut = &contentTypeOut

						// Also the IN content type must match
						matchesIn := false

						// No content type given, and none expected : OK
						if !contentTypeParser.HasContentType() && len(allContentTypeIn) == 0 {
							matchesIn = true
							resultContentTypeIn = nil
						}

						// Content type is given and was found in resource : OK
						if contentTypeParser.HasContentType() && len(allContentTypeIn) > 0 {
							for _, contentTypeIn := range allContentTypeIn {
								if contentTypeIn == contentTypeParser.GetContentType() {
									matchesIn = true
									resultContentTypeIn = &contentTypeIn
								}
							}
						}
						if matchesIn {
							return v, resultContentTypeIn, resultContentTypeOut
						}
					}
				}
			}
		}
	}

	return nil, nil, nil
}
