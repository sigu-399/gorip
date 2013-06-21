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
// repository-name  gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
//
// description      Server implementation.
//
// created          03-03-2013

package gorip

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type Server struct {
	pattern string
	address string
	router  *router

	documentationEndpointEnabled bool
	documentationEndpointUrl     string

	debugEnableLogRequestDump       bool
	debugEnableLogRequestIdentifier bool
	debugEnableLogRequestDuration   bool
}

func NewServer(pattern string, address string) *Server {

	log.Printf("=== Create RIP Server\n")
	return &Server{pattern: pattern, address: address, router: newRouter()}

}

func (s *Server) NewEndpoint(route string, resources ...Resource) error {

	endp := &endpoint{route: route}

	if len(resources) == 0 {
		return errors.New("Endpoint must have at least one resource")
	}

	for _, res := range resources {
		endp.AddResource(res)
	}

	log.Printf("New endpoint : %s\n", endp.GetRoute())

	return s.router.NewEndpoint(endp)
}

func (s *Server) DebugEnableLogRequestDump(b bool) {
	s.debugEnableLogRequestDump = b
}

func (s *Server) DebugEnableLogRequestIdentifier(b bool) {
	s.debugEnableLogRequestIdentifier = b
}

func (s *Server) DebugEnableLogRequestDuration(b bool) {
	s.debugEnableLogRequestDuration = b
}

func (s *Server) ListenAndServe() error {

	log.Printf("=== Listening on %s\n", s.address)

	http.Handle(s.pattern, s)
	return http.ListenAndServe(s.address, nil)
}

func (s *Server) DebugPrintRouterTree() {

	log.Printf("=== Router Tree ================= \n")
	s.router.PrintRouterTree()
	log.Printf("=== End of Router Tree ========== \n")

}

func (s *Server) EnableDocumentationEndpoint(url string) {

	log.Printf("Enable documentation on endpoint %s\n", url)

	s.documentationEndpointEnabled = true
	s.documentationEndpointUrl = url

}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	var timeStart time.Time
	var timeEnd time.Time

	if s.debugEnableLogRequestDuration {
		timeStart = time.Now()
	}

	requestId := "o" // No request id
	if s.debugEnableLogRequestIdentifier {
		requestId = s.generateRequestId(timeStart)
	}

	urlPath := request.URL.Path
	method := request.Method

	log.Printf("[%s] Request %s %s", requestId, method, urlPath)

	if s.debugEnableLogRequestDump {
		s.dumpRequest(request, requestId)
	}

	// Serves documentation if requested and enabled
	if s.documentationEndpointEnabled && s.documentationEndpointUrl == urlPath {
		s.serveDocumentation(writer)
		return
	}

	// Find route node and associated route variables
	node, routeVariables, err := s.router.FindNodeByRoute(urlPath)
	if err != nil {
		errorMessage := err.Error()
		log.Printf(errorMessage)
		s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(errorMessage)}, `text/plain`, requestId)
	}

	resourceContext := ResourceContext{}

	if node == nil {
		message := fmt.Sprintf("[%s] Could not find route for %s", requestId, urlPath)
		log.Printf(message)
		s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
	} else {

		// Route was found:

		// Add route variables to the context
		resourceContext.RouteVariables = routeVariables

		if node.GetEndpoint() == nil {
			message := fmt.Sprintf("[%s] No endpoint found for this route %s", requestId, urlPath)
			log.Printf(message)
			s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
		} else {

			// Endpoint was found:

			// Parse Content-Type and Accept headers

			contentTypeParser, err := newContentTypeHeaderParser(request.Header.Get(`Content-Type`))
			if err != nil {
				message := fmt.Sprintf("[%s] Invalid Content-Type header : %s", requestId, err.Error())
				log.Printf(message)
				s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
			} else {

				acceptParser, err := newAcceptHeaderParser(request.Header.Get(`Accept`))
				if err != nil {
					message := fmt.Sprintf("[%s] Invalid Accept header : %s", requestId, err.Error())
					log.Printf(message)
					s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
				} else {

					if !acceptParser.HasAcceptElement() {
						message := fmt.Sprintf("[%s] No valid Accept header was given", requestId)
						log.Printf(message)
						s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
					} else {

						// Headers are OK:

						endp := node.GetEndpoint()
						availableResourceImplementations := endp.GetResources()

						if len(availableResourceImplementations) == 0 {
							message := fmt.Sprintf("[%s] No resource found on this route %s", requestId, urlPath)
							log.Printf(message)
							s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
						} else {

							matchingResource, contentTypeIn, contentTypeOut := endp.FindMatchingResource(method, &contentTypeParser, &acceptParser)

							if matchingResource == nil {
								message := fmt.Sprintf("[%s] No available resource for this Content-Type", requestId)
								log.Printf(message)
								s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
							} else {

								// Found a matching resource implementation:

								// Add expected content type to the context
								resourceContext.ContentTypeIn = contentTypeIn
								resourceContext.ContentTypeOut = contentTypeOut

								// Read request body

								bodyInBytes, err := ioutil.ReadAll(request.Body)
								if err != nil {
									message := fmt.Sprintf("[%s] Could not read request body", requestId)
									log.Printf(message)
									s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
								} else {

									if resourceContext.ContentTypeIn == nil && len(bodyInBytes) > 0 {
										message := fmt.Sprintf("[%s] Body is not allowed for this resource", requestId)
										log.Printf(message)
										s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
									} else {

										resourceContext.Body = bytes.NewBuffer(bodyInBytes)

										// Create a new instance from factory and executes it
										resource := matchingResource.Factory()
										if resource == nil {
											message := fmt.Sprintf("[%s] Resource factory must instanciate a valid Resource", requestId)
											log.Printf(message)
											s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
										} else {

											// Check and provide query parameters

											resourceContext.QueryParameters = make(map[string]string)
											urlValues := request.URL.Query()

											queryParameterOk := true

											for qpKey, qpObject := range resource.QueryParameters() {
												qpValue := urlValues.Get(qpKey)
												if qpValue == `` {
													qpValue = qpObject.DefaultValue
													if !qpObject.IsValidType(qpValue) {
														message := fmt.Sprintf("[%s] Query parameter %s default value must be of kind %s", requestId, qpKey, qpObject.Kind)
														log.Printf(message)
														s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
														queryParameterOk = false
														break
													}
												}
												if !qpObject.IsValidType(qpValue) {
													message := fmt.Sprintf("[%s] Query parameter %s must be of kind %s", requestId, qpKey, qpObject.Kind)
													log.Printf(message)
													s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
													queryParameterOk = false
													break
												} else {

													// Validate query param
													validator := qpObject.FormatValidator
													if validator != nil {
														if !validator.IsValid(qpValue) {
															message := fmt.Sprintf("[%s] Invalid Query Parameter, %s", requestId, validator.GetErrorMessage())
															log.Printf(message)
															s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
															break
														}
													}

													resourceContext.QueryParameters[qpKey] = qpValue
												}

											}

											// Finally...
											if queryParameterOk {
												result := resource.Execute(&resourceContext)
												s.renderResourceResult(writer, &result, *resourceContext.ContentTypeOut, requestId)
											}
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	if s.debugEnableLogRequestDuration {
		timeEnd = time.Now()
		durationMs := timeEnd.Sub(timeStart).Seconds() * 1000
		log.Printf("[%s] Response Duration : %2.2f ms", requestId, durationMs)
	}
}

func (s *Server) generateRequestId(t time.Time) string {
	hasher := sha1.New()
	hasher.Write([]byte(t.String()))
	return base64.URLEncoding.EncodeToString(hasher.Sum(nil))
}

func (s *Server) dumpRequest(request *http.Request, requestId string) {
	jsonRequest, _ := json.MarshalIndent(request, "", "")
	log.Printf("[%s] === Request dump =================", requestId)
	fmt.Printf("%s\n", jsonRequest)
	log.Printf("[%s] === End of Request dump ==========", requestId)
}

func (s *Server) NewRouteVariableType(kind string, rvtype RouteVariableType) error {
	return s.router.NewRouteVariableType(kind, rvtype)
}

func (s *Server) renderResourceResult(writer http.ResponseWriter, result *ResourceResult, contentType string, requestId string) {

	bodyOutLen := 0
	if result.Body != nil {
		bodyOutLen = result.Body.Len()
	}

	writer.Header().Set(`Content-Length`, strconv.Itoa(bodyOutLen))

	if bodyOutLen > 0 {
		writer.Header().Add(`Content-Type`, contentType)
	}

	writer.WriteHeader(result.HttpStatus)

	if bodyOutLen > 0 {
		_, err := result.Body.WriteTo(writer)
		if err != nil {
			log.Printf("[%s] Error while writing the body %s", requestId, err.Error())
		}
	}

	jsonResult, _ := json.Marshal(result)
	log.Printf("[%s] Response result : %s", requestId, jsonResult)

}
