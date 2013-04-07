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
// description      Utilities for dealing with http headers.
// 
// created          08-03-2013

package gorip

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

type contentTypeHeaderParser struct {
	contentType *string
	parameters  map[string]string
}

func newContentTypeHeaderParser(value string) (contentTypeHeaderParser, error) {
	p := contentTypeHeaderParser{}
	p.parameters = make(map[string]string)

	err := p.parse(value)
	if err != nil {
		return contentTypeHeaderParser{}, err
	}
	return p, nil
}

func (p *contentTypeHeaderParser) parse(value string) error {

	if value != `` {

		split := strings.Split(value, `;`)

		if len(split) == 1 {

			trimmed := strings.TrimSpace(value)
			p.contentType = &trimmed

		} else {

			trimmed := strings.TrimSpace(split[0])
			p.contentType = &trimmed

			for i := 1; i != len(split); i++ {
				trimmedParameter := strings.TrimSpace(split[i])
				splitParameter := strings.Split(trimmedParameter, `=`)
				if len(splitParameter) == 2 {
					p.parameters[splitParameter[0]] = splitParameter[1]
				} else {
					return errors.New(`Invalid Content-Type parameter : expecting key-value`)
				}
			}
		}
	}

	return nil
}

func (p *contentTypeHeaderParser) HasContentType() bool {
	return p.contentType != nil
}

func (p *contentTypeHeaderParser) GetContentType() string {
	return *p.contentType
}

type acceptHeaderParser struct {
	contentTypes []acceptHeaderElementParser
}

type acceptHeaderElementParsers []acceptHeaderElementParser

func (s acceptHeaderElementParsers) Len() int      { return len(s) }
func (s acceptHeaderElementParsers) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

type sortByPriority struct{ acceptHeaderElementParsers }

func (s sortByPriority) Less(i, j int) bool {
	return s.acceptHeaderElementParsers[i].priority > s.acceptHeaderElementParsers[j].priority
}

func newAcceptHeaderParser(value string) (acceptHeaderParser, error) {
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

	sort.Sort(sortByPriority{p.contentTypes})

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
