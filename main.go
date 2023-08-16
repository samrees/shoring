package main

import (
	"fmt"
	"os/exec"
	"regexp"

	"github.com/Masterminds/semver/v3"

	//"log"
	"os"
	"path"
)

var homeDir string
var asdfDir string

type VersionType int

const (
	SEMVER VersionType = iota
	LEXICAL
)

type Requirement struct {
	name        string
	constraint  string
	versionType VersionType
}

var brewReqs = []Requirement{
	Requirement{
		name:        "bash",
		constraint:  ">=5",
		versionType: SEMVER,
	},
	Requirement{
		name:        "fd",
		constraint:  ">=8.5",
		versionType: SEMVER,
	},

	Requirement{
		name:        "ripgrep",
		constraint:  ">=13",
		versionType: SEMVER,
	},
}

var asdfPluginReqs = []string{"awscli", "clisso", "fly", "helm-sops", "helm", "kubectl", "kubectx", "sops", "terraform"}

var asdfReqs = []Requirement{
	Requirement{
		name:        "fly",
		constraint:  ">=7.8",
		versionType: SEMVER,
	},
}

func compareSemver(versionString string, constraintString string) bool {
	version, err := semver.NewVersion(versionString)
	if err != nil {
		// not a valid semver version
		return false
	}
	constraint, _ := semver.NewConstraint(constraintString)
	if constraint.Check(version) {
		return true
	}
	return false
}

func compareLexical(versionString string, constraintString string) bool {
	onlydigits := regexp.MustCompile(`^[\d\._-]+$`)
	stripped := regexp.MustCompile(`[\._-]`)
	if onlydigits.MatchString(versionString) {
		if stripped.ReplaceAllString(versionString, "") >= stripped.ReplaceAllString(constraintString, "") {
			return true
		}
	}
	return false
}

func checkBrew(req Requirement) bool {
	files, err := os.ReadDir(path.Join("/usr/local/Cellar/", req.name))
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.IsDir() {
			switch req.versionType {
			case SEMVER:
				if compareSemver(file.Name(), req.constraint) {
					return true
				}
			case LEXICAL:
				if compareLexical(file.Name(), req.constraint) {
				}
			}
		}
	}
	return false
}

func checkAsdf(req Requirement) bool {
	files, err := os.ReadDir(path.Join(asdfDir, req.name))
	if err != nil {
		return false
	}
	for _, file := range files {
		if file.IsDir() {
			switch req.versionType {
			case SEMVER:
				if compareSemver(file.Name(), req.constraint) {
					return true
				}
			case LEXICAL:
				if compareLexical(file.Name(), req.constraint) {
				}
			}
		}
	}
	return false
}

func baselineBrew() bool {
	_, err := exec.LookPath("brew")
	if err != nil {
		return false
	}
	return true
}
func baselineAsdf() bool {
	_, err := exec.LookPath("asdf")
	if err != nil {
		return false
	}
	return true
}

func main() {
	if !baselineBrew() {
		fmt.Println("# brew does not seem installed, run this to fix:")
		fmt.Println("xcode-select --install")
		fmt.Println("/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)\"")
	}
	if !baselineAsdf() {
		fmt.Println("# asdf does not seem installed, run this to fix:")
		fmt.Println("brew install asdf")
	}

	// brew
	missingBrew := make([]Requirement, len(brewReqs))
	for _, req := range brewReqs {
		if !checkBrew(req) {
			missingBrew = append(missingBrew, req)
		}
	}
	if len(missingBrew) > 0 {
		for _, req := range missingBrew {
			fmt.Sprintln("# brew package %s is failed to meet constraint %s, run this to fix:", req.name, req.constraint)
			fmt.Sprintln("brew upgrade %s || brew install %s", req.name, req.name)
		}
	}

	homeDir, _ = os.UserHomeDir()
	asdfDir = path.Join(homeDir, ".asdf")

	// asdf plugins
	//homeDir, _ := os.UserHomeDir()
	//asdfDir := path.Join(homeDir, ".asdf")
	missingAsdfPlugins := make([]string, len(asdfPluginReqs))

	pluginList, err := os.ReadDir(path.Join(asdfDir, "plugins"))
	if err != nil {
		missingAsdfPlugins = asdfPluginReqs
	}

	for _, req := range asdfPluginReqs {
		found := false
		for _, file := range pluginList {
			if file.Name() == req {
				found = true
			}
		}
		if !found {
			missingAsdfPlugins = append(missingAsdfPlugins, req)
		}
	}
	if len(missingAsdfPlugins) > 0 {
		for _, plugin := range missingAsdfPlugins {
			fmt.Sprintln("# asdf plugin %s is missing, run this to fix:")
			fmt.Sprintln("asdf plugin add %s", plugin)
		}
	}

	// asdf installs
	missingAsdf := make([]Requirement, len(asdfReqs))
	for _, req := range asdfReqs {
		if !checkAsdf(req) {
			missingAsdf = append(missingAsdf, req)
		}
	}
	if len(missingAsdf) > 0 {
		for _, req := range missingAsdf {
			fmt.Sprintln("# asdf install %s failed to meet constraint %s, run this to fix:", req.name, req.constraint)
			fmt.Sprintln("asdf install %s %s", req.name, req.constraint)
		}
	}
}
