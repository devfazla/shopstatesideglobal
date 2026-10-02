package markdown

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// caseVector mirrors one entry of the shared testdata/clean-html-cases.json,
// which is the single source of truth across the Go, JavaScript, and Python
// ports of the clean-html function.
type caseVector struct {
	Name string `json:"name"`
	In   string `json:"in"`
	Want string `json:"want"`
}

func loadSharedCases(t *testing.T) []caseVector {
	t.Helper()
	path := filepath.FromSlash("../../testdata/clean-html-cases.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading shared vectors: %v", err)
	}
	var doc struct {
		Cases []caseVector `json:"cases"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		t.Fatalf("parsing shared vectors: %v", err)
	}
	if len(doc.Cases) == 0 {
		t.Fatal("no shared vectors found")
	}
	return doc.Cases
}

func TestCleanHTMLSharedVectors(t *testing.T) {
	for _, c := range loadSharedCases(t) {
		t.Run(c.Name, func(t *testing.T) {
			got := strings.TrimSuffix(CleanHTML(c.In), "\n")
			if got != c.Want {
				t.Errorf("CleanHTML(%q)\n got: %q\nwant: %q", c.In, got, c.Want)
			}
		})
	}
}
