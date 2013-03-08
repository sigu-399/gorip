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
	"strconv"
	"strings"
)

type contentTypeHeaderParser struct {
	contentType *string
	charset     *string
}

func NewContentTypeHeaderParser(value string) contentTypeHeaderParser {
	p := contentTypeHeaderParser{}
	p.parse(value)
	return p
}

func (p *contentTypeHeaderParser) parse(value string) {

	if value != `` {

		split := strings.Split(value, `;`)

		if len(split) == 0 {

			trimmed := strings.TrimSpace(value)
			p.contentType = &trimmed

		} else {

			trimmed := strings.TrimSpace(split[0])
			p.contentType = &trimmed

			for i := range split {

				trimmed := strings.TrimSpace(split[i])

				// TODO : check charsets ?
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
	contentTypes []acceptHeaderElementParser
}

type acceptHeaderElementParser struct {
	contentType string
	priority    float64
}

func newAcceptHeaderElementParser(value string) acceptHeaderElementParser {
	p := acceptHeaderElementParser{}
	p.parse(value)
	return p
}

func NewAcceptHeaderParser(value string) acceptHeaderParser {
	p := acceptHeaderParser{}
	p.parse(value)
	return p
}

func (p *acceptHeaderParser) parse(value string) {

	if value != `` {
		split := strings.Split(value, `,`)
		for i := range split {
			element := newAcceptHeaderElementParser(split[i])
			p.contentTypes = append(p.contentTypes, element)
		}
	}
}

func (p *acceptHeaderElementParser) parse(value string) {

	if value != `` {
		p.priority = 1 // default

		split := strings.Split(value, `;`)
		if len(split) == 1 {
			p.contentType = strings.TrimSpace(split[0])
		} else {
			if len(split) == 2 {
				p.contentType = strings.TrimSpace(split[0])

				trimmed := split[1]
				if strings.HasPrefix(trimmed, `q=`) {
					splitPriority := strings.Split(trimmed, `=`)
					if len(splitPriority) == 2 {
						fValue, err := strconv.ParseFloat(splitPriority[1], 64)
						if err == nil {
							p.priority = fValue
						}
					}
				}
			}
		}
	}

}
