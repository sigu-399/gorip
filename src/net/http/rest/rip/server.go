// author			sigu-399
// author-github	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc	REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		Server implementation.
// 
// created			03-03-2013

package rip

import (
	"bytes"
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
}

func NewServer(pattern string, address string) *Server {
	return &Server{pattern: pattern, address: address, router: newRouter()}
}

func (s *Server) ListenAndServe() error {

	http.Handle(s.pattern, s)
	return http.ListenAndServe(s.address, nil)
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {

	timeStart := time.Now()

	urlPath := request.URL.Path
	method := request.Method

	resourceContext := ResourceContext{}

	log.Printf("Requesting %s %s", method, urlPath)

	// Check if the documentation endpoint is enabled
	if s.documentationEndpointEnabled && s.documentationEndpointUrl == urlPath {
		s.serveDocumentation(writer)
		return
	}

	// Find route node and associated route variables
	node, routeVariables, err := s.router.FindNodeByRoute(urlPath)

	if err != nil {
		message := err.Error()
		log.Printf(message)
		s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
	}

	if node == nil {
		message := fmt.Sprintf("Could not find route for %s", urlPath)
		log.Printf(message)
		s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
	} else {

		// Route was found:

		// Add route variables to the context
		resourceContext.RouteVariables = routeVariables

		if node.GetEndpoint() == nil {
			message := fmt.Sprintf("No endpoint found for this route %s", urlPath)
			log.Printf(message)
			s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`)
		} else {

			// Endpoint was found:

			// Parse Content-Type and Accept headers

			contentTypeParser, err := newContentTypeHeaderParser(request.Header.Get(`Content-Type`))
			if err != nil {
				message := fmt.Sprintf("Invalid Content-Type header : %s", err.Error())
				log.Printf(message)
				s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
			} else {

				acceptParser, err := newAcceptHeaderParser(request.Header.Get(`Accept`))
				if err != nil {
					message := fmt.Sprintf("Invalid Accept header : %s", err.Error())
					log.Printf(message)
					s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
				} else {

					if !acceptParser.HasAcceptElement() {
						message := fmt.Sprintf("No valid Accept header was given")
						log.Printf(message)
						s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
					} else {

						// Headers are OK:

						endp := node.GetEndpoint()
						availableResourceImplementations := endp.GetResources()

						if len(availableResourceImplementations) == 0 {
							message := fmt.Sprintf("No resource found on this route %s", urlPath)
							log.Printf(message)
							s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`)
						} else {

							matchingResource, contentTypeIn, contentTypeOut := endp.FindMatchingResource(method, &contentTypeParser, &acceptParser)

							if matchingResource == nil {
								message := fmt.Sprintf("No available resource for this Content-Type")
								log.Printf(message)
								s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
							} else {

								// Found a matching resource implementation: 

								// Add expected content type to the context 
								resourceContext.ContentTypeIn = contentTypeIn
								resourceContext.ContentTypeOut = contentTypeOut

								// Read request body

								bodyInBytes, err := ioutil.ReadAll(request.Body)
								if err != nil {
									message := fmt.Sprintf("Could not read request body")
									log.Printf(message)
									s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`)
								} else {

									if resourceContext.ContentTypeIn == nil && len(bodyInBytes) > 0 {
										message := fmt.Sprintf("Body is not allowed for this resource")
										log.Printf(message)
										s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
									} else {

										resourceContext.Body = bytes.NewBuffer(bodyInBytes)

										// Create a new instance from factory and executes it
										resource := matchingResource.Factory()
										if resource == nil {
											message := fmt.Sprintf("Resource factory must instanciate a valid Resource")
											log.Printf(message)
											s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusInternalServerError, Body: bytes.NewBufferString(message)}, `text/plain`)
										} else {

											// Check and provide query parameters

											resourceContext.QueryParameters = make(map[string]string)
											urlValues := request.URL.Query()

											queryParameterOk := true

											for qpKey, qpObject := range resource.GetQueryParameters() {
												qpValue := urlValues.Get(qpKey)
												if qpValue == `` {
													qpValue = qpObject.DefaultValue
													if !qpObject.IsValidType(qpValue) {
														message := fmt.Sprintf("Query parameter %s default value must be of kind %s", qpKey, qpObject.Kind)
														log.Printf(message)
														s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
														queryParameterOk = false
														break
													}
												}
												if !qpObject.IsValidType(qpValue) {
													message := fmt.Sprintf("Query parameter %s must be of kind %s", qpKey, qpObject.Kind)
													log.Printf(message)
													s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
													queryParameterOk = false
													break
												} else {

													// Validate query param													
													validator := qpObject.FormatValidator
													if validator != nil {
														if !validator.IsValid(qpValue) {
															message := fmt.Sprintf("Query Parameter is invalid : ", validator.GetErrorMessage())
															log.Printf(message)
															s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`)
															break
														}
													}

													resourceContext.QueryParameters[qpKey] = qpValue
												}

											}

											// Finally...
											if queryParameterOk {
												result := resource.Execute(&resourceContext)
												s.renderResourceResult(writer, &result, *resourceContext.ContentTypeOut)
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

	timeEnd := time.Now()
	durationMs := timeEnd.Sub(timeStart).Seconds() * 1000

	log.Printf("Response time : %2.2f ms", durationMs)

}

func (s *Server) RegisterEndpoint(e *endpoint) error {

	if e == nil {
		panic(`Endpoint cannot be nil`)
	}

	log.Printf("Registering endpoint : %s\n", e.GetRoute())
	return s.router.RegisterEndpoint(e)
}

func (s *Server) RegisterDocumentationEndpoint(url string) {
	s.documentationEndpointEnabled = true
	s.documentationEndpointUrl = url
}

func (s *Server) RegisterRouteVariableValidator(kind string, validator RouteVariableValidator) error {
	return s.router.RegisterRouteVariableValidator(kind, validator)
}

func (s *Server) renderResourceResult(writer http.ResponseWriter, result *ResourceResult, contentType string) {

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
			log.Printf("Error while writing the body %s", err.Error())
		}
	}

}
