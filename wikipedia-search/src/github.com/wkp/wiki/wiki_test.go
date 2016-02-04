package wiki_test

import(
	"strings"
	"testing"
	"github.com/wkp/wiki"
)
type testpair struct {
    value string
    expected string
}

type testpairint struct {
    value string
    expected int
}

var testsRemoveCite = []testpair{
    { "xpto {{ xpto2 }} xpto3", "xpto xpto3" },
     { "xpto {{ xpto2 }} xpto3 {{xpto4}} ", "xpto xpto3" },
}

var testsRemoveRef = []testpair{
    { "xpto <ref ffdsdfsdf xpto2> sdfsdff </ref> xpto3", "xpto xpto3" },
     { "xpto <ref> sdfsdff </ref>  xpto3 ", "xpto xpto3" },
}

var testsRemoveMarkup = []testpair{
    { "xpto === == '' ' ''' <br> <br/> [http://www ] <div dfdsfsdf> </div> #REDIRECT <math>1+2</math> xpto3", "xpto xpto3" },
}

var testsRemoveLink = []testpair{
    { "xpto [[xpto| xpto2]] xpto3 [[xpto| xpto4]] xpto5", "xpto xpto2 xpto3 xpto4 xpto5" },
    { "[[xpto2]]", "xpto2" },
}

var testsRemoveDuplication = []testpairint{
    { "!@#$ˆxpto?><,./ xpto xpto, xpto2 xpto2 xpto2 \n\tXPTO2 #$%ˆ&*()", 2 },
}

func TestRemoveCite(t *testing.T) {
    for _, pair := range testsRemoveCite {
        v := wiki.RemoveCite(pair.value)
        if v != pair.expected {
            t.Error(
                "For", pair.value, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}

func TestRemoveRef(t *testing.T) {
    for _, pair := range testsRemoveRef {
        v := wiki.RemoveRef(pair.value)
        if v != pair.expected {
            t.Error(
                "For", pair.value, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}


func TestRemoveMarkup(t *testing.T) {
    for _, pair := range testsRemoveMarkup {
        v := wiki.RemoveMarkup(pair.value)
        if v != pair.expected {
            t.Error(
                "For", pair.value, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}

func TestRemoveLink(t *testing.T) {
    for _, pair := range testsRemoveLink {
        v := wiki.RemoveLink(pair.value)
        if v != pair.expected {
            t.Error(
                "For", pair.value, 
                "expected", pair.expected,
                "got", v,
            )
        }
    }
}

func TestRemoveDuplication(t *testing.T) {
    for _, pair := range testsRemoveDuplication {
        v := wiki.RemoveDuplication(pair.value)
        if len(strings.Split(v, " ")) != pair.expected {
            t.Error(
                "For", pair.value, 
                "expected", pair.expected,
                "got", len(strings.Split(v, " ")),
            )
        }
    }
}