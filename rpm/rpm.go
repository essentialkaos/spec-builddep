package rpm

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/essentialkaos/ek/v12/sliceutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

type pkgIndex struct {
	ByName     map[string]string
	ByProvides map[string]string
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Versions returns map with versions of installed packages
func Versions(packages []string) map[string]string {
	packages, index := normalizePackageNames(packages)
	return getVersions(packages, index)
}

// ////////////////////////////////////////////////////////////////////////////////// //

// normalizePackageNames normalizes packages names
func normalizePackageNames(packages []string) ([]string, *pkgIndex) {
	var result []string

	index := getIndex(packages)

	for _, p := range packages {
		if index.ByName[p] != "" {
			result = append(result, index.ByName[p])
		} else {
			result = append(result, p)
		}
	}

	return result, index
}

// getIndex returns index with given packages names and provided packages names
func getIndex(packages []string) *pkgIndex {
	index := &pkgIndex{
		ByName:     make(map[string]string),
		ByProvides: make(map[string]string),
	}

	cmd := exec.Command("rpm", "-q", "--whatprovides", "--qf", "%{name}\n")
	cmd.Args = append(cmd.Args, packages...)
	data, _ := cmd.CombinedOutput()

	if len(data) == 0 {
		return index
	}

	lines := sliceutil.Deduplicate(strings.Split(string(data), "\n"))

	for i := 0; i < len(packages); i++ {
		if strings.Contains(lines[i], " ") {
			continue
		}

		index.ByName[packages[i]] = lines[i]
		index.ByProvides[lines[i]] = packages[i]
	}

	return index
}

// getVersions returns map with versions info for installed packages
func getVersions(packages []string, index *pkgIndex) map[string]string {
	result := map[string]string{}

	cmd := exec.Command("rpm", "-q", "--qf", "%{name} %{version}\n")
	cmd.Args = append(cmd.Args, packages...)
	data, _ := cmd.CombinedOutput()

	if len(data) == 0 {
		return result
	}

	buf := bytes.NewBuffer(data)

	for {
		line, err := buf.ReadString('\n')

		if err != nil {
			break
		}

		if strings.Count(line, " ") > 1 {
			continue
		}

		name, version, _ := strings.Cut(line, " ")

		if index.ByProvides[name] == "" {
			result[name] = strings.Trim(version, "\n\r")
		} else {
			result[index.ByProvides[name]] = strings.Trim(version, "\n\r")
		}
	}

	return result
}
