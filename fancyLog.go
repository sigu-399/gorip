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
// description      Fancy log and terminal using colors and log types.
//
// created          20-07-2013

package gorip

import (
	"fmt"
	"log"
)

const (
	TERM_COLOR_BLACK    = 0
	TERM_COLOR_RED      = 1
	TERM_COLOR_GREEN    = 2
	TERM_COLOR_YELLOW   = 3
	TERM_COLOR_BLUE     = 4
	TERM_COLOR_MAGENTA  = 5
	TERM_COLOR_CYAN     = 6
	TERM_COLOR_WHITE    = 7
)

func TermColorEscape(message string, c int) string {
	return fmt.Sprintf("\033[%dm%s\033[0m", 30+c, message)
}

type FLOG_TYPE int8

const (
	FLOG_TYPE_INFO FLOG_TYPE = iota
	FLOG_TYPE_WARNING
	FLOG_TYPE_ERROR
	FLOG_TYPE_DEBUG
	FLOG_TYPE_ACTION
)

func Flog(t FLOG_TYPE, m string) {

	c := TERM_COLOR_BLUE
	ts := "NFO"

	switch t {

	case FLOG_TYPE_INFO:
		c = TERM_COLOR_CYAN
		ts = "NFO"

	case FLOG_TYPE_WARNING:
		c = TERM_COLOR_YELLOW
		ts = "WRN"

	case FLOG_TYPE_ERROR:
		c = TERM_COLOR_RED
		ts = "ERR"

	case FLOG_TYPE_ACTION:
		c = TERM_COLOR_GREEN
		ts = "ACT"

	case FLOG_TYPE_DEBUG:
		c = TERM_COLOR_MAGENTA
		ts = "DBG"

	default:
		log.Printf(m)
		return
	}

	log.Printf(fmt.Sprintf(`%s %s`, TermColorEscape(ts, c), m))

}
