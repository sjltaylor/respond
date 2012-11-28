/*

This package implements a partial Accept header parser as defined by RFC 2616 section 14.1
See: http://www.w3.org/Protocols/rfc2616/rfc2616-sec14.html

Also see: http://www.xml.com/pub/a/2005/06/08/restful.html

Limitations:

* Although detailed in hte RFC, no functionality to indicate media type preferences is implemented here.
* AcceptsMediaType(...) bool disregards media range parameters

*/

package accept

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MediaRange struct {
	Type    string
	Subtype string
	Q       float64
	Params  map[string]string
}

func NewMediaRange() *MediaRange {
	return &MediaRange{
		Params: make(map[string]string),
	}
}

type Header struct {
	MediaRanges []*MediaRange
}

func newHeader() *Header {
	return &Header{make([]*MediaRange, 0)}
}

func ParseAcceptHeader(request *http.Request) (header *Header, err error) {
	return parseAcceptHeaderString(request.Header.Get(`Accept`))
}

func parseAcceptHeaderString(s string) (header *Header, err error) {

	defer func() {
		if recoveredError := recover(); recoveredError != nil {
			header = nil
			err = fmt.Errorf("Accept header not valid: '%s'", s)
		}
	}()

	s = strings.ToLower(s)

	h := newHeader()

	mediaRangeStrings := strings.Split(s, ",")

	for _, mediaRangeString := range mediaRangeStrings {

		mediaRangeString = strings.TrimSpace(mediaRangeString)

		mediaRange := parseMediaRange(mediaRangeString)

		h.MediaRanges = append(h.MediaRanges, mediaRange)
	}

	return h, nil
}

func parseMediaRange(mediaRangeString string) *MediaRange {
	//text/*;q=0.3, text/html;q=0.7, text/html;level=1;pref=44, text/html;level=2;q=0.4, */*;q=0.5, application/json;q=0.1

	parts := strings.Split(mediaRangeString, ";")

	mediaRange := NewMediaRange()

	typeAndSubtype := strings.Split(parts[0], "/")

	mediaRange.Type = typeAndSubtype[0]
	mediaRange.Subtype = typeAndSubtype[1]

	for _, paramString := range parts[1:] {

		nameValue := strings.Split(paramString, "=")
		name := nameValue[0]
		value := nameValue[1]

		if name == "q" {

			if f, err := strconv.ParseFloat(value, 64); err != nil {
				panic(fmt.Errorf("can't parse q parameter as float: %s : %s", value, err))
			} else {
				mediaRange.Q = f
			}

		} else {

			mediaRange.Params[name] = value
		}

	}

	return mediaRange
}

/*
	panics if the mediaType is not valid
*/
func (header *Header) AcceptsMediaType(mediaType string) bool {

	mediaRange := parseMediaRange(mediaType) // (media types can be parsed as media ranges)

	for _, acceptedMediaRange := range header.MediaRanges {

		if acceptedMediaRange.Type == "*" {
			return true
		}

		if acceptedMediaRange.Type != mediaRange.Type {
			continue
		}

		if acceptedMediaRange.Subtype == "*" {
			return true
		}

		if acceptedMediaRange.Subtype == mediaRange.Subtype {
			return true
		}
	}

	return false
}
