// Copyright 2013 sigu-399 ( https://github.com/sigu-399 )
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author           sigu-399
// author-github    https://github.com/sigu-399
// author-mail      sigu.399@gmail.com
//
// repository-nam   gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
//
// description      Defines a REST endpoint.
//                  An endpoint is a url route and associated resourceHandlers.
//
// created          07-03-2013

package gorip

import ()

type endpoint struct {
	route            string
	resourceHandlers []ResourceHandler
}

func (e *endpoint) GetRoute() string {
	return e.route
}

func (e *endpoint) AddResource(rh ResourceHandler) {
	e.resourceHandlers = append(e.resourceHandlers, rh)
}

func (e *endpoint) GetResourceHandlers() []ResourceHandler {
	return e.resourceHandlers
}

func (e *endpoint) FindMatchingResource(method string, contentTypeParser *contentTypeHeaderParser, acceptParser *acceptHeaderParser) (*ResourceHandler, *string, *string) {

	var resultContentTypeIn *string
	var resultContentTypeOut *string

	// Loop through accepted OUT content types, highest priority first
	for _, acceptElement := range acceptParser.contentTypes {
		// Find a resource for given method
		for _, v := range e.resourceHandlers {

			if v.Method == method {

				allContentTypeIn := v.ContentTypeIn
				allContentTypeOut := v.ContentTypeOut

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
							return &v, resultContentTypeIn, resultContentTypeOut
						}
					}
				}
			}
		}
	}

	return nil, nil, nil
}
