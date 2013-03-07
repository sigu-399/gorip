// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description		A resource is an implementation of a REST method.
// 
// created      	07-03-2013

package gorip

import (
	"net/http"
)

type Resource interface {
	Execute(writer http.ResponseWriter, request *http.Request)
}
