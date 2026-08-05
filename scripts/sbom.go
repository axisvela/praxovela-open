// Command sbom generates sbom.json at the repo root — a CycloneDX 1.5 software
// bill of materials covering the transitive Go dependencies of every module in
// the workspace.
//
// Modules are discovered with the same pattern as release-manifest.go
// (packages/*/go.mod go directives). Each module's transitive dependency set is
// resolved with `go list -m -json all` run in that module's directory, offline
// (GOPROXY=off; versions come from the local module cache). Emitted library
// components are deduplicated by module path + version.
//
//	go run ./scripts/sbom.go
//
// The output file is a build artifact (gitignored).
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type sbomComponent struct {
	Type    string `json:"type"`
	Group   string `json:"group,omitempty"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

type sbomDoc struct {
	BOMFormat    string          `json:"bomFormat"`
	SpecVersion  string          `json:"specVersion"`
	SerialNumber string          `json:"serialNumber"`
	Version      int             `json:"version"`
	Metadata     sbomMetadata    `json:"metadata"`
	Components   []sbomComponent `json:"components"`
}

type sbomMetadata struct {
	Component sbomComponent `json:"component"`
}

type sbomVersion struct {
	Version string `json:"version"`
	Name    string `json:"name"`
}

type sbomModule struct {
	Path string
	Dir  string
	Go   string
}

// goModInfo mirrors a record of `go list -m -json all`.
type goModInfo struct {
	Path    string `json:"Path"`
	Version string `json:"Version"`
	Main    bool   `json:"Main"`
}

var sbomGoDirectiveRE = regexp.MustCompile(`(?m)^go[ \t]+(\S+)`)

func main() {
	root, err := sbomRepoRoot()
	if err != nil {
		sbomFatal("%v", err)
	}

	ver, err := sbomReadVersion(root)
	if err != nil {
		sbomFatal("read version.json: %v", err)
	}

	mods, err := sbomDiscoverModules(root)
	if err != nil {
		sbomFatal("discover modules: %v", err)
	}
	if len(mods) == 0 {
		sbomFatal("no Go modules found under packages/")
	}

	seen := make(map[string]bool)
	var comps []sbomComponent
	for _, m := range mods {
		info, err := sbomListModules(m.Dir)
		if err != nil {
			sbomFatal("%v", err)
		}
		for _, dep := range info {
			if dep.Main || dep.Version == "" {
				continue
			}
			key := dep.Path + "@" + dep.Version
			if seen[key] {
				continue
			}
			seen[key] = true
			comps = append(comps, sbomComponent{
				Type:    "library",
				Group:   sbomGroup(dep.Path),
				Name:    dep.Path,
				Version: dep.Version,
			})
		}
	}

	sort.Slice(comps, func(i, j int) bool {
		if comps[i].Name != comps[j].Name {
			return comps[i].Name < comps[j].Name
		}
		return comps[i].Version < comps[j].Version
	})

	appName := ver.Name
	if appName == "" {
		appName = "AxisRobo Agent / PRAXOVELA"
	}

	doc := sbomDoc{
		BOMFormat:    "CycloneDX",
		SpecVersion:  "1.5",
		SerialNumber: sbomUUID(),
		Version:      1,
		Metadata: sbomMetadata{
			Component: sbomComponent{Type: "application", Name: appName, Version: ver.Version},
		},
		Components: comps,
	}

	out := filepath.Join(root, "sbom.json")
	if err := sbomWriteJSON(out, doc); err != nil {
		sbomFatal("write %s: %v", out, err)
	}
	fmt.Printf("wrote %s (%d modules, %d library components)\n", out, len(mods), len(comps))
}

// sbomRepoRoot returns the git working tree root.
func sbomRepoRoot() (string, error) {
	out, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", fmt.Errorf("git rev-parse --show-toplevel: %w", err)
	}
	return strings.TrimSpace(string(out)), nil
}

// sbomReadVersion parses version.json at the repo root.
func sbomReadVersion(root string) (sbomVersion, error) {
	data, err := os.ReadFile(filepath.Join(root, "version.json"))
	if err != nil {
		return sbomVersion{}, err
	}
	var v sbomVersion
	if err := json.Unmarshal(data, &v); err != nil {
		return sbomVersion{}, err
	}
	return v, nil
}

// sbomDiscoverModules finds every packages/* module via its go directive,
// mirroring discoverModules in release-manifest.go.
func sbomDiscoverModules(root string) ([]sbomModule, error) {
	dirs, err := os.ReadDir(filepath.Join(root, "packages"))
	if err != nil {
		return nil, err
	}
	var mods []sbomModule
	for _, d := range dirs {
		if !d.IsDir() {
			continue
		}
		gf := filepath.Join(root, "packages", d.Name(), "go.mod")
		data, err := os.ReadFile(gf)
		if err != nil {
			continue
		}
		m := sbomGoDirectiveRE.FindStringSubmatch(string(data))
		if m == nil {
			continue
		}
		mods = append(mods, sbomModule{Path: "packages/" + d.Name(), Dir: filepath.Dir(gf), Go: m[1]})
	}
	sort.Slice(mods, func(i, j int) bool { return mods[i].Path < mods[j].Path })
	return mods, nil
}

// sbomListModules resolves the transitive build list of the module in dir with
// `go list -m -json all`, offline. Records are stream-decoded.
func sbomListModules(dir string) ([]goModInfo, error) {
	cmd := exec.Command("go", "list", "-m", "-json", "all")
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GOPROXY=off")
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("go list -m -json all in %s: %v: %s", dir, err, strings.TrimSpace(stderr.String()))
	}
	var mods []goModInfo
	dec := json.NewDecoder(strings.NewReader(string(out)))
	for {
		var m goModInfo
		if err := dec.Decode(&m); err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("go list -m -json all in %s: decode: %w", dir, err)
		}
		mods = append(mods, m)
	}
	return mods, nil
}

// sbomGroup returns the path portion before the last "/" ("" when none).
func sbomGroup(path string) string {
	if i := strings.LastIndex(path, "/"); i >= 0 {
		return path[:i]
	}
	return ""
}

// sbomUUID returns a random version-4 UUID prefixed with "urn:uuid:".
func sbomUUID() string {
	var b [16]byte
	if _, err := rand.Read(b[:]); err != nil {
		sbomFatal("random serial number: %v", err)
	}
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b[:])
	return fmt.Sprintf("urn:uuid:%s-%s-%s-%s-%s", h[0:8], h[8:12], h[12:16], h[16:20], h[20:32])
}

func sbomWriteJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func sbomFatal(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "sbom: "+format+"\n", args...)
	os.Exit(1)
}
