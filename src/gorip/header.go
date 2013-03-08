// author  			sigu-399
// author-github 	https://github.com/sigu-399
// author-mail		sigu.399@gmail.com
// 
// repository-name	gorip
// repository-desc  REST Server Framework - ( gorip: REST In Peace ) - Go language
// 
// description	    Utilities for dealing with http headers.
// 
// created      	08-03-2013

package gorip

import (
	"strings"
)

type contentTypeHeaderParser struct {
	contentType *string
	charset     *string
}

func NewContentTypeHeaderParser(value string) *contentTypeHeaderParser {
	p := &contentTypeHeaderParser{}
	p.parse(value)
	return p
}

func (p *contentTypeHeaderParser) parse(value string) {

	if value != `` {
		split := strings.Split(value, `;`)
		if len(split) == 0 {
			p.contentType = &value
		} else {
			p.contentType = &split[0]
			for i := range split {
				trimmed := strings.TrimSpace(split[i])
				if strings.HasPrefix(trimmed, `charset=`) {
					splitCharset := strings.Split(trimmed, `=`)
					if len(splitCharset) == 2 {
						p.charset = &splitCharset[1]
					}
				}
			}
		}
	}
}

func (p *contentTypeHeaderParser) HasContentType() bool {
	return p.contentType != nil
}

func (p *contentTypeHeaderParser) HasCharset() bool {
	return p.charset != nil
}

func (p *contentTypeHeaderParser) GetContentType() string {
	return *p.contentType
}

func (p *contentTypeHeaderParser) GetCharset() string {
	return *p.charset
}

type acceptHeaderParser struct {
}
