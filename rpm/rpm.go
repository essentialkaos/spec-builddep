package rpm

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
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
	versions := make(map[string]string)

	for _, pkg := range packages {
		cmd := exec.Command("rpm", "-q", "--whatprovides", "--qf", "%{version}\n", pkg)
		data, _ := cmd.CombinedOutput()

		if len(data) == 0 {
			continue
		}

		nlIndex := bytes.IndexRune(data, '\n')

		if nlIndex == -1 {
			continue
		}

		version := string(data[:nlIndex])

		if !strings.Contains(version, " ") {
			versions[pkg] = version
		}
	}

	return versions
}

// ////////////////////////////////////////////////////////////////////////////////// //
