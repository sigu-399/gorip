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
	HttpMethodGET     HttpMethod = "GET"
	HttpMethodHEAD    HttpMethod = "HEAD"
	HttpMethodPOST    HttpMethod = "POST"
	HttpMethodPUT     HttpMethod = "PUT"
	HttpMethodPATCH   HttpMethod = "PATCH"
	HttpMethodDELETE  HttpMethod = "DELETE"
	HttpMethodTRACE   HttpMethod = "TRACE"
	HttpMethodCONNECT HttpMethod = "CONNECT"
)
