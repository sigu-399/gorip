// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description	    Generates the server endpoints documentation.
// 
// created      	19-03-2013

package rip

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) serveDocumentation(writer http.ResponseWriter) {

	documentation := new(bytes.Buffer)

	currentNode := s.router.rootNode
	currentPath := `(root)`

	documentation.WriteString(`<h1>Documentation</h1>`)

	s.serveDocumentationRecursive(currentPath, currentNode, documentation)

	bodyOutLen := len(documentation.Bytes())

	writer.Header().Set(`Content-Length`, strconv.Itoa(bodyOutLen))

	if bodyOutLen > 0 {
		writer.Header().Add(`Content-Type`, `text/html`)
	}

	writer.WriteHeader(http.StatusOK)

	if bodyOutLen > 0 {
		documentation.WriteTo(writer)
	}

}

func (s *Server) serveDocumentationRecursive(currentPath string, currentNode routerNode, buffer *bytes.Buffer) {

	path := currentPath + currentNode.GetPart()

	endpoint := currentNode.GetEndpoint()
	if endpoint != nil {

		buffer.WriteString(`<hr/>`)

		buffer.WriteString(`<h2>` + path + `</h2>`)

		resources := endpoint.GetResources()

		for _, r := range resources {

			buffer.WriteString(`<hr/>`)

			buffer.WriteString(`<p>Method : ` + r.GetMethod() + `</p>`)
			buffer.WriteString(`<p>In : ` + strings.Join(r.GetContentTypeIn(), `,`) + `</p>`)
			buffer.WriteString(`<p>Out : ` + strings.Join(r.GetContentTypeOut(), `,`) + `</p>`)

			qps := r.GetQueryParameters()
			for key, q := range qps {
				buffer.WriteString(`<p>QueryParam : ` + key + ` [type ` + q.Kind + `][default ` + q.DefaultValue + `]</p>`)
			}
		}
	}

	children := currentNode.GetChildren()
	for _, c := range children {
		s.serveDocumentationRecursive(path+`/`, c, buffer)
	}

}
