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
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func (s *Server) serveDocumentation(writer http.ResponseWriter) {

	documentation := new(bytes.Buffer)

	currentNode := s.router.rootNode
	currentPath := ``

	documentation.WriteString(`<html>` + "\n")

	documentation.WriteString(`<head>` + "\n")
	documentation.WriteString(`<title>REST Server Documentation - gorip</title>` + "\n")

	documentation.WriteString(`<style type="text/css">` + "\n")
	documentation.WriteString(`body {margin: 0;font-family: Helvetica, Arial, sans-serif;font-size: 16px;color: #222}` + "\n")
	documentation.WriteString(`h1, h2, h3, h4 {color:#37AB5E;margin:20px;padding:0}` + "\n")
	documentation.WriteString(`h1 {font-size: 24px; font-weight: bold}` + "\n")
	documentation.WriteString(`h2 {font-size: 20px;background: #E0F5EB;padding: 2px 5px}` + "\n")
	documentation.WriteString(`h3 {font-size: 18px}` + "\n")
	documentation.WriteString(`h4 {font-size: 16px}` + "\n")
	documentation.WriteString(`p, pre, table{margin: 20px}` + "\n")
	documentation.WriteString(`pre {background: none repeat scroll 0 0 #E9E9E9;border-radius: 5px 5px 5px 5px;padding: 10px;}` + "\n")
	documentation.WriteString(`td {background-color:#F0F5EB}` + "\n")

	documentation.WriteString(`</style>` + "\n")

	documentation.WriteString(`</head>` + "\n")

	documentation.WriteString(`<body>` + "\n")

	documentation.WriteString(`<h1>REST Documentation</h1>` + "\n")

	s.serveDocumentationRecursive(currentPath, currentNode, documentation)

	documentation.WriteString(`</body>` + "\n")
	documentation.WriteString(`</html>` + "\n")

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

		buffer.WriteString(`<h2>` + path + `</h2>` + "\n")

		resourceHandlers := endpoint.GetResourceHandlers()

		for _, r := range resourceHandlers {

			buffer.WriteString(`<h3>` + r.Method + `</h3>` + "\n")

			buffer.WriteString(`<table>` + "\n")
			buffer.WriteString(`<tr><th></th><th>Content Type</th></tr>` + "\n")
			buffer.WriteString(`<tr><td>In</td><td>` + strings.Join(r.ContentTypeIn, `,`) + `</td>` + "\n")
			buffer.WriteString(`<tr><td>Out</td><td>` + strings.Join(r.ContentTypeOut, `,`) + `</td>` + "\n")
			buffer.WriteString(`</table>` + "\n")

			buffer.WriteString(`<table>` + "\n")
			buffer.WriteString(`<tr><th>Query Param</th><th>Kind</th><th>Default</th></tr>` + "\n")
			qps := r.QueryParameters
			for key, q := range qps {
				buffer.WriteString(`<tr><td>` + key + `</td><td>` + q.Kind + `</td><td>` + q.DefaultValue + `</td><tr>`)
				if q.FormatValidator != nil {
					buffer.WriteString(` Validator : ` + q.FormatValidator.GetErrorMessage())
				}
			}
			buffer.WriteString(`</table>` + "\n")

			if r.Documentation != nil {

				rhd := r.Documentation

				if len(rhd.TestURL) > 0 {

					testResultString := ""

					client := &http.Client{}
					req, _ := http.NewRequest(r.Method, rhd.TestURL, nil)
					req.Header.Add(`Accept`, rhd.TestContentType)
					resp, err := client.Do(req)
					if err != nil {
						testResultString = "ERROR : " + err.Error()
					} else {
						bodyInBytes, err := ioutil.ReadAll(resp.Body)
						if err != nil {
							testResultString = "ERROR : " + err.Error()
						} else {
							testResultString = string(bodyInBytes)
						}
					}

					buffer.WriteString(fmt.Sprintf(`<h4>Example (%s)</h4>`+"\n", rhd.TestURL))
					buffer.WriteString(`<pre>` + testResultString + `</pre>` + "\n")
				}

				if len(rhd.AdditionalNotes) > 0 {
					buffer.WriteString(`<h4>Additional Notes</h4>` + "\n")
					buffer.WriteString(`<p>` + rhd.AdditionalNotes + `</p>` + "\n")
				}
			}
		}
	}

	children := currentNode.GetChildren()
	for _, c := range children {
		s.serveDocumentationRecursive(path+`/`, c, buffer)
	}

}
