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
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Versions returns map with versions of installed packages
func Versions(packages []string) map[string]string {
	return getPackagesVersions(packages, getNormNames(packages))
}

// ////////////////////////////////////////////////////////////////////////////////// //

func getNormNames(packages []string) []string {
	if !hasVirtualPackages(packages) {
		return packages
	}

	virtualNames := normalizeVirtualPackagesNames(getVirtualPackages(packages))

	var result []string

	for _, p := range packages {
		if virtualNames[p] != "" {
			result = append(result, virtualNames[p])
		} else {
			result = append(result, p)
		}
	}

	return result
}

// hasVirtualPackages returns true if slice with packages names contains virtual
// packages
func hasVirtualPackages(packages []string) bool {
	for _, p := range packages {
		if strings.ContainsRune(p, '(') {
			return true
		}
	}

	return false
}

// getVirtualPackages returns slice with virtual packages
func getVirtualPackages(packages []string) []string {
	var result []string

	for _, p := range packages {
		if strings.ContainsRune(p, '(') {
			result = append(result, p)
		}
	}

	return result
}

// normalizeVirtualPackagesNames normalizes virtual packages names
// (e.g. python3dist(setuptools) → platform-python-setuptools)
func normalizeVirtualPackagesNames(packages []string) map[string]string {
	result := map[string]string{}

	cmd := exec.Command("rpm", "-q", "--whatprovides", "--qf", "%{name}\n")
	cmd.Args = append(cmd.Args, packages...)
	data, _ := cmd.Output()

	if len(data) == 0 {
		return result
	}

	buf := bytes.NewBuffer(data)

	for i := 0; i < len(packages); i++ {
		line, err := buf.ReadString('\n')

		if err != nil {
			break
		}

		if strings.Contains(line, " ") {
			continue
		}

		result[packages[i]] = strings.Trim(line, "\n\r")
	}

	return result
}

// getPackagesVersions returns map with installed packages versions
func getPackagesVersions(packages, normPackages []string) map[string]string {
	result := map[string]string{}

	cmd := exec.Command("rpm", "-q", "--qf", "%{version}\n")
	cmd.Args = append(cmd.Args, normPackages...)
	data, _ := cmd.Output()

	if len(data) == 0 {
		return result
	}

	buf := bytes.NewBuffer(data)

	for i := 0; i < len(packages); i++ {
		line, err := buf.ReadString('\n')

		if err != nil {
			break
		}

		if strings.Contains(line, " ") {
			continue
		}

		result[packages[i]] = strings.Trim(line, "\n\r")
	}

	return result
}
