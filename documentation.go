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
// description      Generates the server endpoints documentation.
// 
// created          19-03-2013

package gorip

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
				buffer.WriteString(`<p>QueryParam : ` + key + ` [type ` + q.Kind + `][default ` + q.DefaultValue + `]`)
				if q.FormatValidator != nil {
					buffer.WriteString(` Validator : ` + q.FormatValidator.GetErrorMessage())
				}
				buffer.WriteString(`</p>`)
			}

			buffer.WriteString(`<h3>Custom infos</h3>`)
			buffer.WriteString(`<p>` + r.GetCustomDocumentation() + `</p>`)

		}
	}

	children := currentNode.GetChildren()
	for _, c := range children {
		s.serveDocumentationRecursive(path+`/`, c, buffer)
	}

}
