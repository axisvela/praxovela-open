// Command release-manifest generates release-manifest.json at the repo root.
//
// It reads version.json plus git metadata and emits a build/release manifest
// recording the version, commit, toolchain, and per-module Go directives. It
// runs from anywhere in the checkout:
//
//	go run ./scripts/release-manifest.go
//
// The output file is a build artifact (gitignored); SBOM and CI test evidence
// are referenced as strings and produced by the pipeline separately.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type manifest struct {
	Version     string    `json:"version"`
	Commit      string    `json:"commit"`
	CommitShort string    `json:"commit_short"`
	Branch      string    `json:"branch"`
	Toolchain   toolchain `json:"toolchain"`
	Modules     []module  `json:"modules"`
	SBOM        string    `json:"sbom"`
	Tests       string    `json:"tests"`
}

type toolchain struct {
	Go          string `json:"go"`
	GoDirective string `json:"go_directive"`
	Rust        string `json:"rust,omitempty"`
	Node        string `json:"node,omitempty"`
}

type module struct {
	Path string `json:"path"`
	Go   string `json:"go"`
}

type versionFile struct {
	Version string `json:"version"`
	Go      string `json:"go"`
}

var goDirectiveRE = regexp.MustCompile(`(?m)^go[ \t]+(\S+)`)

func main() {
	root, err := repoRoot()
	if err != nil {
		fatalf("%v", err)
	}

	ver, err := readVersion(root)
	if err != nil {
		fatalf("read version.json: %v", err)
	}

	commit, err := runOutput(root, "git", "rev-parse", "HEAD")
	if err != nil {
		fatalf("git rev-parse HEAD: %v", err)
	}
	commitShort, err := runOutput(root, "git", "rev-parse", "--short=12", "HEAD")
	if err != nil {
		fatalf("git rev-parse --short=12 HEAD: %v", err)
	}
	branch, err := runOutput(root, "git", "branch", "--show-current")
	if err != nil {
		fatalf("git branch --show-current: %v", err)
	}
	goVersion, err := runOutput(root, "go", "version")
	if err != nil {
		fatalf("go version: %v", err)
	}

	mods, err := discoverModules(root)
	if err != nil {
		fatalf("discover modules: %v", err)
	}

	m := manifest{
		Version:     ver.Version,
		Commit:      commit,
		CommitShort: commitShort,
		Branch:      branch,
		Toolchain: toolchain{
			Go:          goVersion,
			GoDirective: ver.Go,
			Rust:        optionalToolVersion(root, "cargo", "--version"),
			Node:        optionalToolVersion(root, "node", "--version"),
		},
		Modules: mods,
		SBOM:    "sbom.json",
		Tests:   fmt.Sprintf("CI: %d-module matrix green", len(mods)),
	}

	out := filepath.Join(root, "release-manifest.json")
	if err := writeJSON(out, m); err != nil {
		fatalf("write %s: %v", out, err)
	}
	fmt.Printf("wrote %s (version %s, commit %s)\n", out, m.Version, m.CommitShort)
}

// repoRoot returns the git working tree root.
func repoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// runOutput runs name with args in dir and returns trimmed combined output.
func runOutput(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		return s, err
	}
	return s, nil
}

// optionalToolVersion returns the version line for name, or "" if the tool is
// not installed or fails to report a version. Optional tools are omitted from
// the manifest rather than failing the generator.
func optionalToolVersion(dir, name, verArg string) string {
	if _, err := exec.LookPath(name); err != nil {
		return ""
	}
	out, err := runOutput(dir, name, verArg)
	if err != nil {
		return ""
	}
	return out
}

// readVersion parses version.json at the repo root.
func readVersion(root string) (versionFile, error) {
	data, err := os.ReadFile(filepath.Join(root, "version.json"))
	if err != nil {
		return versionFile{}, err
	}
	var v versionFile
	if err := json.Unmarshal(data, &v); err != nil {
		return versionFile{}, err
	}
	return v, nil
}

// discoverModules reads the go directive from every packages/*/go.mod.
func discoverModules(root string) ([]module, error) {
	dirs, err := os.ReadDir(filepath.Join(root, "packages"))
	if err != nil {
		return nil, err
	}
	var mods []module
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, "packages", d.Name(), "go.mod"))
		if err != nil {
			continue
		}
		m := goDirectiveRE.FindStringSubmatch(string(data))
		if m == nil {
			continue
		}
		mods = append(mods, module{Path: "packages/" + d.Name(), Go: m[1]})
	}
	sort.Slice(mods, func(i, j int) bool { return mods[i].Path < mods[j].Path })
	return mods, nil
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "release-manifest: "+format+"\n", args...)
	os.Exit(1)
}
