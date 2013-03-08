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

package gorip

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	pattern string
	address string
	router  *router
}

func NewServer(pattern string, address string) *Server {
	return &Server{pattern: pattern, address: address, router: NewRouter()}
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

	// Find route node and associated route variables
	node, routeVariables, err := s.router.FindNodeByRoute(urlPath)
	if err != nil {
		log.Printf("Warning : %s", err.Error())
	}

	if node == nil {
		log.Printf("Warning : Could not find route for %s", urlPath)
	} else {

		// Route was found:

		// Add route variables to the context
		resourceContext.routeVariables = routeVariables

		if node.GetEndpoint() == nil {
			log.Printf("Warning : No endpoint found for this route")
		} else {

			// Endpoint was found:

			// Parse Content-Type and Accept headers

			contentTypeParser, err := NewContentTypeHeaderParser(request.Header.Get(`Content-Type`))
			log.Printf("contentTypeParser : %s", contentTypeParser)
			if err != nil {
				log.Printf(`Invalid Content-Type header : ` + err.Error())
			}

			fmt.Printf("Accept : %s\n", request.Header.Get(`Accept`))
			acceptParser, err := NewAcceptHeaderParser(request.Header.Get(`Accept`))
			if err != nil {
				log.Printf(`Invalid Accept header : ` + err.Error())
			}

			if !acceptParser.HasAcceptElement() {
				log.Printf(`No valid Accept header was given`)
			} else {

				// Headers are OK:

				endp := node.GetEndpoint()
				availableResourceImplementations := endp.GetResources()

				if len(availableResourceImplementations) == 0 {
					log.Printf(`No resource found on this route`)
				} else {

					matchingResource := endp.FindMatchingResource(method, &acceptParser)

					if matchingResource == nil {
						log.Printf(`No available resource for this Content-Type`)
					} else {

						// Found a matching resource implementation: 

						// Executes it
						matchingResource.Execute(&resourceContext)
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

func (s *Server) RegisterRouteVariableValidator(kind string, validator RouteVariableValidator) error {
	return s.router.RegisterRouteVariableValidator(kind, validator)
}
