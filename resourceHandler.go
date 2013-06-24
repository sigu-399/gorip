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
// description      A resource is an implementation of a REST method.
//
// created          07-03-2013

package gorip

import (
	"bytes"
)

type ResourceHandler interface {
	Factory() ResourceHandler
	Execute(context *ResourceHandlerContext) ResourceHandlerResult
	Signature() *ResourceHandlerSignature
	QueryParameters() map[string]QueryParameter
	DocumentationNotes() string
}

type ResourceHandlerSignature struct {
	Method         string
	ContentTypeIn  []string
	ContentTypeOut []string
}

type ResourceHandlerContext struct {
	RouteVariables  map[string]string
	QueryParameters map[string]string
	ContentTypeIn   *string
	ContentTypeOut  *string
	Body            *bytes.Buffer
}

type ResourceHandlerResult struct {
	HttpStatus int
	Body       *bytes.Buffer
}
