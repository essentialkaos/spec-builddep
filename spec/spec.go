package spec

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"bytes"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/sliceutil"
	"github.com/essentialkaos/ek/v12/sortutil"
	"github.com/essentialkaos/ek/v12/strutil"
)

// ////////////////////////////////////////////////////////////////////////////////// //

const (
	EQ Cond = 1 // Equal | =
	LT Cond = 2 // Less than | <
	LE Cond = 3 // Less or equal | <=
	GT Cond = 4 // Greater than | >
	GE Cond = 5 // Greater or equal | >=
)

// ////////////////////////////////////////////////////////////////////////////////// //

type Cond uint8

type BuildDep struct {
	Name    string
	Version string
	Cond    Cond
}

type BuildDeps []BuildDep

// ////////////////////////////////////////////////////////////////////////////////// //

func (p BuildDeps) Len() int           { return len(p) }
func (p BuildDeps) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p BuildDeps) Less(i, j int) bool { return sortutil.NaturalLess(p[i].Name, p[j].Name) }

// ////////////////////////////////////////////////////////////////////////////////// //

// GetDeps returns slice with build dependencies from given spec file
func GetDeps(spec string) (BuildDeps, error) {
	err := fsutil.ValidatePerms("FRS", spec)

	if err != nil {
		return nil, err
	}

	cmd := exec.Command("rpmspec", "-P", spec)
	data, err := cmd.CombinedOutput()

	if err != nil {
		errText := strings.ReplaceAll(string(data), "\n", "")
		errText = strings.Replace(errText, "error: ", "", 1)
		return nil, fmt.Errorf("Spec parsing error: %s", errText)
	}

	return extractBuildDeps(data), nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// Names returns slice with dependencies names
func (d BuildDeps) Names(withVersion bool) []string {
	var result []string

	for _, dep := range d {
		if dep.Version == "" || !withVersion {
			result = append(result, dep.Name)
		} else {
			result = append(result, fmt.Sprintf(
				"%s %s %s", dep.Name, dep.Cond.Clause(), dep.Version,
			))
		}
	}

	return result
}

// String returns string representation of build dependency
func (d BuildDep) String() string {
	if d.Cond != 0 {
		return fmt.Sprintf("{%s %v %s}", d.Name, d.Cond, d.Version)
	}

	return fmt.Sprintf("{%s}", d.Name)
}

// String returns string representation of condition
func (c Cond) String() string {
	switch c {
	case EQ:
		return "="
	case LT:
		return "<"
	case LE:
		return "≤"
	case GT:
		return ">"
	case GE:
		return "≥"
	}

	return ""
}

// Clause returns condition as clause
func (c Cond) Clause() string {
	switch c {
	case EQ:
		return "="
	case LT:
		return "<"
	case LE:
		return "<="
	case GT:
		return ">"
	case GE:
		return ">="
	}

	return ""
}

// ////////////////////////////////////////////////////////////////////////////////// //

// extractBuildDeps extracts dependencies info from spec data
func extractBuildDeps(data []byte) []BuildDep {
	var result BuildDeps

	buf := bytes.NewBuffer(data)

	for {
		line, err := buf.ReadString('\n')

		if err != nil {
			break
		}

		if line == "" || !strings.Contains(line, "BuildRequires: ") {
			continue
		}

		line = strings.TrimRight(line, " \n\r")
		line = strings.ReplaceAll(line, "BuildRequires:", "")
		line = strings.TrimLeft(line, " \t")

		result = append(result, parseDepsLine(line)...)
	}

	sort.Sort(result)

	return sliceutil.Deduplicate(result)
}

// parseDepsLine parses line with one or more dependencies
func parseDepsLine(line string) []BuildDep {
	var result BuildDeps

	deps := strutil.Fields(line)

	for i := 0; i < len(deps); i++ {
		// Check next two items for condition and version
		if i+2 < len(deps) && strings.ContainsAny(deps[i+1], "<>=") {
			result = append(result, BuildDep{
				Name:    deps[i],
				Cond:    parseCond(deps[i+1]),
				Version: deps[i+2],
			})
			i += 2
		} else {
			result = append(result, BuildDep{Name: deps[i]})
		}
	}

	return result
}

// parseCond parses version condition
func parseCond(cond string) Cond {
	switch cond {
	case "=":
		return EQ
	case "<":
		return LT
	case "<=":
		return LE
	case ">":
		return GT
	case ">=":
		return GE
	}

	return 0
}
