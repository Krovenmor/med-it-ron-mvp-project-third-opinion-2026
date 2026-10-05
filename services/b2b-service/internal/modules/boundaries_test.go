package modules_test

import (
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

const modulesImport = "github.com/Krovenmor/med-it-ron-mvp-project-third-opinion-2026/services/b2b-service/internal/modules/"

var sharedPackages = map[string]bool{
	"system/history": true,
}

func TestModulesTalkOnlyThroughPublicAPI(t *testing.T) {
	entries, err := os.ReadDir(".")
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		module := entry.Name()
		err := filepath.WalkDir(module, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".go") {
				return err
			}
			file, err := parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly)
			if err != nil {
				return err
			}
			for _, imp := range file.Imports {
				target, _ := strconv.Unquote(imp.Path.Value)
				if violation := crossesBoundary(module, target); violation != "" {
					t.Errorf("%s imports %s: %s", path, target, violation)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
}

func crossesBoundary(module, target string) string {
	rest, ok := strings.CutPrefix(target, modulesImport)
	if !ok {
		return ""
	}
	owner, inner, _ := strings.Cut(rest, "/")
	switch {
	case owner == module, sharedPackages[rest]:
		return ""
	case inner != "api":
		return "other modules are reachable only through their api package"
	}
	return ""
}
