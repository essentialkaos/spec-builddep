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

// Versions returns map with versions of installed packages
func Versions(packages []string) map[string]string {
	return getPackagesVersions(packages, getRealNames(packages))
}

// ////////////////////////////////////////////////////////////////////////////////// //

// getRealNames returns slice with real packages names
func getRealNames(packages []string) []string {
	realNames := getRealPackagesNames(packages)

	var result []string

	for _, p := range packages {
		if realNames[p] != "" {
			result = append(result, realNames[p])
		} else {
			result = append(result, p)
		}
	}

	return result
}

// getRealPackagesNames returns map package name → real package name
func getRealPackagesNames(packages []string) map[string]string {
	result := map[string]string{}

	cmd := exec.Command("rpm", "-q", "--whatprovides", "--qf", "%{name}\n")
	cmd.Args = append(cmd.Args, packages...)
	data, _ := cmd.CombinedOutput()

	if len(data) == 0 {
		return result
	}

	lines := sliceutil.Deduplicate(strings.Split(string(data), "\n"))

	for i := 0; i < len(packages); i++ {
		if strings.Contains(lines[i], " ") {
			continue
		}

		result[packages[i]] = lines[i]
	}

	return result
}

// getPackagesVersions returns map with installed packages versions
func getPackagesVersions(packages, realPackages []string) map[string]string {
	result := map[string]string{}

	cmd := exec.Command("rpm", "-q", "--qf", "%{name} %{version}\n")
	cmd.Args = append(cmd.Args, realPackages...)
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
		result[name] = strings.Trim(version, "\n\r")
	}

	return result
}
