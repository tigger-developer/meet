// ABOUTME: Regression tests for issue #9 — meet help text is sourced from
// ABOUTME: docs/help/, not hardcoded inline in Go source.

package regression

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const helpDocsDir = "docs/help"

func repoRootForHelp(t *testing.T) string {
	t.Helper()
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(cwd, "..", ".."))
}

// AC9.1 — top-level `meet -h` prose comes from docs/help/meet.txt. The full
// contents of the docs file appear in the help output.
func TestMeetHelp_TopLevel_SourcedFromDocs_RT9_1(t *testing.T) {
	root := repoRootForHelp(t)
	docSrc, err := os.ReadFile(filepath.Join(root, helpDocsDir, "meet.txt"))
	if err != nil {
		t.Fatalf("read docs source: %v", err)
	}
	stdout, stderr, _ := runMeet(t, t.TempDir(), "-h")
	output := stdout + stderr
	if !strings.Contains(output, strings.TrimSpace(string(docSrc))) {
		t.Errorf("meet -h does not contain docs/help/meet.txt content\n--- docs ---\n%s\n--- output ---\n%s",
			string(docSrc), output)
	}
}

// AC9.1 — each subcommand's `-h` prose comes from docs/help/meet-<sub>.txt.
func TestMeetHelp_SubcommandsSourcedFromDocs_RT9_2(t *testing.T) {
	root := repoRootForHelp(t)
	subs := []string{"serve", "token", "create", "cancel", "list"}
	for _, sub := range subs {
		t.Run(sub, func(t *testing.T) {
			docPath := filepath.Join(root, helpDocsDir, "meet-"+sub+".txt")
			docSrc, err := os.ReadFile(docPath)
			if err != nil {
				t.Fatalf("read docs source %s: %v", docPath, err)
			}
			stdout, stderr, _ := runMeet(t, t.TempDir(), sub, "-h")
			output := stdout + stderr
			if !strings.Contains(output, strings.TrimSpace(string(docSrc))) {
				t.Errorf("meet %s -h does not contain docs/help/meet-%s.txt content\n--- docs ---\n%s\n--- output ---\n%s",
					sub, sub, string(docSrc), output)
			}
		})
	}
}

// AC9.1 — help documents double-hyphen flags and does not expose Go flag's
// single-hyphen PrintDefaults output.
func TestMeetHelp_DocOwnedDoubleHyphenFlags_RT9_4(t *testing.T) {
	stdout, stderr, _ := runMeet(t, t.TempDir(), "token", "--help")
	output := stdout + stderr
	if !strings.Contains(output, "--room <name>") {
		t.Errorf("meet token --help does not include docs-owned --room option; output=%q", output)
	}
	if strings.Contains(output, "\n  -room") || strings.Contains(output, "\n  -config") {
		t.Errorf("meet token --help exposes single-hyphen flag defaults; output=%q", output)
	}
}
