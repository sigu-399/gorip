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

	log.Printf("Requesting %s %s", method, urlPath)

	node, variables, err := s.router.FindNodeByRoute(urlPath)
	if err != nil {
		log.Printf("Warning : %s", err.Error())
	}

	if node == nil {
		log.Printf("Warning : Could not find route for %s", urlPath)
	} else {
		log.Printf("node : %s", node)
		log.Printf("variables : %s", variables)

		if node.GetEndpoint() == nil {
			log.Printf("Warning : No endpoint found on this route")
		} else {
			contentTypeParser, err := NewContentTypeHeaderParser(request.Header.Get(`Content-Type`))
			if err != nil {
				log.Printf(`Invalid Content-Type header : ` + err.Error())
			} else {
				fmt.Printf("%s\n", contentTypeParser)
			}
			acceptParser, err := NewAcceptHeaderParser(request.Header.Get(`Accept`))
			if err != nil {
				log.Printf(`Invalid Accept header : ` + err.Error())
			} else {
				fmt.Printf("%s\n", acceptParser)
			}

			if !acceptParser.HasAcceptElement() {
				log.Printf(`No valid Accept header was given`)
			} else {
				//TODO
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
