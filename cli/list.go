package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"github.com/essentialkaos/ek/v12/fmtc"

	"github.com/essentialkaos/spec-builddep/rpm"
	"github.com/essentialkaos/spec-builddep/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// listDeps prints list with required dependencies
func listDeps(specFile string) error {
	deps, err := spec.GetDeps(specFile)

	if err != nil {
		return err
	}

	if len(deps) == 0 {
		fmtc.Printf("{g}Spec file %s has no dependencies{!}", specFile)
		return nil
	}

	if !useRawOutput {
		printDepList(deps)
	} else {
		printRawDepList(deps)
	}

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printDepList prints pretty list with required dependencies
func printDepList(deps spec.BuildDeps) {
	installed := rpm.Versions(deps.Names())

	fmtc.NewLine()

	for _, dep := range deps {
		installedVer := installed[dep.Name]

		if installedVer != "" {
			fmtc.Printf(" {s-}•{!} {g}%s{!}", dep.Name)
		} else {
			fmtc.Printf(" {s-}•{!} %s", dep.Name)
		}

		if dep.Version != "" {
			fmtc.Printf(" {s}%s{!} %s", dep.Cond, dep.Version)
		}

		if installedVer != "" {
			fmtc.Printf(" {s-}(%s){!}", installedVer)
		}

		fmtc.NewLine()
	}

	fmtc.NewLine()
}

// printRawDepList prints raw list with required dependencies
func printRawDepList(deps spec.BuildDeps) {
	for _, dep := range deps {
		if dep.Version != "" {
			fmtc.Printf("%s %s %s\n", dep.Name, dep.Cond, dep.Version)
		} else {
			fmtc.Printf("%s\n", dep.Name)
		}
	}
}
