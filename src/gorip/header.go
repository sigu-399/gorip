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
	"errors"
	"strconv"
	"strings"
)

type contentTypeHeaderParser struct {
	contentType *string
	charset     *string
}

func NewContentTypeHeaderParser(value string) (contentTypeHeaderParser, error) {
	p := contentTypeHeaderParser{}
	err := p.parse(value)
	if err != nil {
		return contentTypeHeaderParser{}, err
	}
	return p, nil
}

func (p *contentTypeHeaderParser) parse(value string) error {

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

				if strings.HasPrefix(trimmed, `charset=`) {
					splitCharset := strings.Split(trimmed, `=`)
					if len(splitCharset) == 2 {
						// TODO : check charsets ?
						p.charset = &splitCharset[1]
					} else {
						return errors.New(`Invalid Content-Type parameter : expecting key-value charset`)
					}
				} else {
					return errors.New(`Invalid Content-Type parameter : expecting charset`)
				}

			}
		}
	}

	return nil
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

func NewAcceptHeaderParser(value string) (acceptHeaderParser, error) {
	p := acceptHeaderParser{}
	err := p.parse(value)
	if err != nil {
		return acceptHeaderParser{}, err
	}
	return p, nil
}

func (p *acceptHeaderParser) parse(value string) error {

	if value != `` {
		split := strings.Split(value, `,`)
		for i := range split {
			element, err := newAcceptHeaderElementParser(split[i])
			if err != nil {
				return err
			}
			p.contentTypes = append(p.contentTypes, element)
		}
	}

	return nil
}

func (p *acceptHeaderParser) HasAcceptElement() bool {
	return len(p.contentTypes) > 0
}

type acceptHeaderElementParser struct {
	contentType string
	priority    float64
}

func newAcceptHeaderElementParser(value string) (acceptHeaderElementParser, error) {
	p := acceptHeaderElementParser{}
	err := p.parse(value)
	if err != nil {
		return acceptHeaderElementParser{}, err
	}
	return p, nil
}

func (p *acceptHeaderElementParser) parse(value string) error {

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
							if fValue >= 0 && fValue <= 1 {
								p.priority = fValue
							} else {
								return errors.New(`Invalid accept element : q value must be a float`)
							}
						} else {
							return errors.New(`Invalid accept element : q value must be a float`)
						}
					} else {
						return errors.New(`Invalid accept element : expecting key-value q`)
					}
				} else {
					return errors.New(`Invalid accept element : expecting q=`)
				}
			} else {
				return errors.New(`Invalid accept element : only 2 parameters are allowed`)
			}
		}
	} else {
		return errors.New(`Accept element cannot be empty`)
	}

	return nil
}
