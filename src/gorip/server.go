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

func (s *Server) GetRouter() *router {
	return s.router
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	timeStart := time.Now()

	// TODO

	timeEnd := time.Now()

	log.Printf("Response : %2.2f ms", timeEnd.Sub(timeStart).Seconds()*1000)

}
