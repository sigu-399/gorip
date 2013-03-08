// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description	    Global definitions.
// 
// created      	08-03-2013

package gorip

type HttpMethod string

const (
	GET     HttpMethod = "GET"
	HEAD    HttpMethod = "HEAD"
	POST    HttpMethod = "POST"
	PUT     HttpMethod = "PUT"
	PATCH   HttpMethod = "PATCH"
	DELETE  HttpMethod = "DELETE"
	TRACE   HttpMethod = "TRACE"
	CONNECT HttpMethod = "CONNECT"
)
