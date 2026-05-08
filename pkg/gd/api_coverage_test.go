package gd

import (
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestLibGDPublicAPICoverage(t *testing.T) {
	headers := configuredHeaders()
	if len(headers) == 0 {
		t.Skip("libgd headers not found")
	}

	source := packageSource(t)
	missing := make([]string, 0)
	for _, symbol := range publicSymbols(t, headers) {
		if apiCoverageExclusions[symbol] != "" {
			continue
		}
		if strings.Contains(source, symbol) {
			continue
		}
		missing = append(missing, symbol)
	}
	if len(missing) > 0 {
		sort.Strings(missing)
		t.Fatalf("missing libgd API coverage decisions:\n%s", strings.Join(missing, "\n"))
	}
}

var apiCoverageExclusions = map[string]string{
	// These are acknowledged through higher-level wrappers that own the memory
	// and do not expose raw libgd allocation to callers.
	"gdFree": "internal memory cleanup only",

	// The legacy gdSource/gdSink APIs are deprecated by libgd in favor of
	// gdIOCtx. Exposing Go callbacks here would add unsafe cgo ownership rules
	// to the public surface.
	"gdImageCreateFromPngSource": "legacy source API; use gdIOCtx wrappers",
	"gdImagePngToSink":           "legacy sink API; use gdIOCtx wrappers",
	"gdNewSSCtx":                 "legacy source/sink API; use gdIOCtx wrappers",

	// Go callbacks across cgo require explicit handle lifetime management, so
	// this low-level hook is intentionally not part of the stable API.
	"gdImageColorReplaceCallback": "unsafe callback API",
}

func configuredHeaders() []string {
	if env := os.Getenv("LIBGD_HEADERS"); env != "" {
		return existingHeaders(strings.Fields(env))
	}
	return existingHeaders([]string{
		"/opt/homebrew/include/gd.h",
		"/opt/homebrew/include/gdfx.h",
		"/usr/local/include/gd.h",
		"/usr/local/include/gdfx.h",
		"/usr/include/gd.h",
		"/usr/include/gdfx.h",
	})
}

func existingHeaders(candidates []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(candidates))
	for _, path := range candidates {
		if seen[path] {
			continue
		}
		if _, err := os.Stat(path); err == nil {
			seen[path] = true
			out = append(out, path)
		}
	}
	return out
}

func publicSymbols(t *testing.T, headers []string) []string {
	t.Helper()
	re := regexp.MustCompile(`BGD_DECLARE\([^)]*\)\s+(gd[A-Za-z0-9_]+)\s*\(`)
	seen := map[string]bool{}
	for _, header := range headers {
		data, err := os.ReadFile(header)
		if err != nil {
			t.Fatal(err)
		}
		normalized := strings.ReplaceAll(string(data), "\n", " ")
		for _, match := range re.FindAllStringSubmatch(normalized, -1) {
			seen[match[1]] = true
		}
	}
	out := make([]string, 0, len(seen))
	for symbol := range seen {
		out = append(out, symbol)
	}
	sort.Strings(out)
	return out
}

func packageSource(t *testing.T) string {
	t.Helper()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var builder strings.Builder
	for _, file := range files {
		if strings.HasSuffix(file, "_test.go") {
			continue
		}
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		builder.Write(data)
		builder.WriteByte('\n')
	}
	return builder.String()
}
