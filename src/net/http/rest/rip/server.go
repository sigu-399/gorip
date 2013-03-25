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
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
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

	// Generate request id
	hasher := sha1.New()
	hasher.Write([]byte(timeStart.String()))
	requestId := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	log.Printf("[%s] Request %s %s", requestId, method, urlPath)

	jsonRequest, _ := json.Marshal(request)
	log.Printf("[%s] Request dump : %s", requestId, jsonRequest)

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
		s.renderResourceResult(writer, &ResourceResult{HttpStatus: http.StatusBadRequest, Body: bytes.NewBufferString(message)}, `text/plain`, requestId)
	}

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

											for qpKey, qpObject := range resource.GetQueryParameters() {
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

	timeEnd := time.Now()
	durationMs := timeEnd.Sub(timeStart).Seconds() * 1000

	log.Printf("[%s] Response time : %2.2f ms", requestId, durationMs)

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
