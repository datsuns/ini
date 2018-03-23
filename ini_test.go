package ini

import (
	"testing"
)

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
