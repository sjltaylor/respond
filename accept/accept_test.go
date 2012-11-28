package accept

import (
	"testing"
)

/*

type/subtype
type/*
*\/*

followed by any number of parameters separated by ;

each media range and its parameters are separeated by ,

the q parameter mist be first

q is the relative preference for the media type

A missing accept parameter is assumed * / *

media ranges with no q paremters are assumed q=1.0

ordered by
*/

var (
	exampleHeader1 = `text/*;q=0.3, text/html;q=0.7, TEXT/html;level=1;pref=44, text/html;level=2;q=0.4, */*;q=0.5, application/json;q=0.1`
	exampleHeader2 = `text/*;q=0.3, text/html;q=0.7, TEXT/html;level=1;pref=44, text/html;level=2;q=0.4, application/json;q=0.1`
)

func parseAcceptHeaderStringForTesting(t *testing.T, s string) *Header {
	h, err := parseAcceptHeaderString(s)

	if err != nil {
		t.Fatalf("could not parse header string: %s: %s", s, err)
	}
	return h
}

func TestParseValidAcceptHeaderString(t *testing.T) {

	h := parseAcceptHeaderStringForTesting(t, exampleHeader1)

	if h == nil {
		t.Fatal("header not returned")
	}

	if len(h.MediaRanges) != 6 {
		t.Fatal(`media range types not parsed correctly`)
	}

	if h.MediaRanges[0].Q != 0.3 {
		t.Fatalf("Q not parsed correctly, got: %f, expected: %f", h.MediaRanges[0].Q, 0.3)
	}

	if h.MediaRanges[2].Params["pref"] != "44" {
		t.Fatalf("Params not parsed correctly, got: pref=%s, expected: pref=%s", h.MediaRanges[2].Params["pref"], "44")
	}

}

func TestParseInvalidAcceptHeaderString(t *testing.T) {

	_, err := parseAcceptHeaderString("invalid accept header")

	if err == nil {
		t.Fatal("expected an error to be returned")
	}
}

func TestAcceptsAllMediaTypesWithWildcardRange(t *testing.T) {

	h := parseAcceptHeaderStringForTesting(t, exampleHeader1)

	if !h.AcceptsMediaType("foo/bar") {
		t.Fatal("should accept all media types")
	}
}

func TestAcceptsMediaType(t *testing.T) {

	h := parseAcceptHeaderStringForTesting(t, exampleHeader2)

	if !h.AcceptsMediaType("text/html") {
		t.Fatal("should accept media type text/html")
	}

	if !h.AcceptsMediaType("application/json") {
		t.Fatal("should accept media type application/json")
	}

	if !h.AcceptsMediaType("text/html") {
		t.Fatal("should accept media type text/html")
	}

	if h.AcceptsMediaType("image/jpeg") {
		t.Fatal("should not accept media type image/jpeg")
	}
}
