package ini

import (
	"strings"
	"testing"
)

func formatInput(list []string) string {
	raw := strings.Join(list, "\n")
	return raw
}

func TestKeepFileHeader(t *testing.T) {
	text := []string{
		"# ini file header is a text lines before the 1st section",
		"dummy",
		"a",
		"[]",
		";[]",
		"\n",
	}
	input := formatInput(text)
	ini, err := LoadText(input)
	if err != nil {
		t.Fatalf("load error [%v] %v", text, err)
	}
	for i, s := range ini.Header {
		if s != text[i] {
			t.Fatalf("invalid header [%v] != [%v]", text[i], s)
		}
	}
	dest := &strings.Builder{}
	ini.RawWrite(dest)
	if input != dest.String() {
		t.Fatalf("invalid output [%v] != [%v]", input, dest.String())
	}
}

func TestParseSectionName(t *testing.T) {
	list := map[string]string{
		"[title]":   "title",
		"[title] ":  "title",
		" title] ":  "title",
		" [title] ": "title",
	}
	for k, v := range list {
		if ParseSectionName(k) != v {
			t.Errorf("expect <%v> but actual <%v>", v, ParseSectionName(k))
		}
	}

}

func TestNoSection(t *testing.T) {
	text := []string{
		"\n",
	}
	input := formatInput(text)
	ini, err := LoadText(input)
	if err != nil {
		t.Fatalf("load error [%v] %v", text, err)
	}
	if ini.HasValue("dummySection", "dummyEntry", "dummyValue") {
		t.Fatalf("ini should not have any entry")
	}
}
