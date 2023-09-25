package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2023 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/essentialkaos/ek/v12/fmtc"
	"github.com/essentialkaos/ek/v12/fsutil"
	"github.com/essentialkaos/ek/v12/options"
	"github.com/essentialkaos/ek/v12/system"
	"github.com/essentialkaos/ek/v12/version"

	"github.com/essentialkaos/spec-builddep/rpm"
	"github.com/essentialkaos/spec-builddep/spec"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// installDeps installs required dependencies
func installDeps(specFile string) error {
	currentUser, err := system.CurrentUser(true)

	if err != nil {
		return fmt.Errorf("Can't check user permissions: %v", err)
	}

	if !currentUser.IsRoot() {
		return fmt.Errorf("Superuser (root) permissions is required to install packages")
	}

	deps, err := spec.GetDeps(specFile)

	if err != nil {
		return err
	}

	quiet := options.GetB(OPT_QUIET)
	depsToInstall := filterRequiredDeps(deps)

	if len(depsToInstall) == 0 {
		fmtc.If(!quiet).Println("{g}All required packages already installed{!}")
		return nil
	}

	if options.GetB(OPT_CLEAN) {
		cleanCache()
	}

	err = installPackages(depsToInstall, len(deps))

	if err != nil {
		return err
	}

	fmtc.If(!quiet).Println("{g}All required packages successfully installed{!}")

	return nil
}

// ////////////////////////////////////////////////////////////////////////////////// //

// filterRequiredDeps filters required dependencies based on currently installed
// packages and provided options
func filterRequiredDeps(deps spec.BuildDeps) spec.BuildDeps {
	if options.GetB(OPT_ACTUAL) {
		return deps
	}

	var result spec.BuildDeps

	installed := rpm.Versions(deps.Names(false))

	for _, dep := range deps {
		ver, ok := installed[dep.Name]

		if !ok {
			result = append(result, dep)
			continue
		}

		if dep.Version == "" {
			continue
		}

		curVer, err := version.Parse(ver)

		if err != nil {
			result = append(result, dep)
			continue
		}

		reqVer, err := version.Parse(dep.Version)

		if err != nil {
			result = append(result, dep)
			continue
		}

		if !isVersionFit(dep, curVer, reqVer) {
			result = append(result, dep)
		}
	}

	return result
}

// isVersionFit returns true if currently installed versions is fit for requirements
func isVersionFit(dep spec.BuildDep, curVer, reqVer version.Version) bool {
	switch dep.Cond {
	case spec.EQ:
		return curVer.Equal(reqVer)
	case spec.LT:
		return curVer.Less(reqVer)
	case spec.LE:
		return curVer.Less(reqVer) || curVer.Equal(reqVer)
	case spec.GT:
		return curVer.Greater(reqVer)
	case spec.GE:
		return curVer.Greater(reqVer) || curVer.Equal(reqVer)
	}

	return false
}

// cleanCache cleans yum/dnf cache
func cleanCache() {
	cmd := exec.Command(getPackageManager())

	cmd.Args = append(cmd.Args, genPackageManagerOptions()...)
	cmd.Args = append(cmd.Args, "clean", "expire-cache")

	cmd.Run()
}

// installPackages install required packages
func installPackages(deps spec.BuildDeps, total int) error {
	quiet := options.GetB(OPT_QUIET)

	fmtc.If(!quiet).Printf(
		"{*}Installing {s}(%d/%d){!}: %s…\n\n",
		len(deps), total, strings.Join(deps.Names(false), ", "),
	)

	cmd := exec.Command(getPackageManager())

	if !quiet {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	cmd.Args = append(cmd.Args, genPackageManagerOptions()...)
	cmd.Args = append(cmd.Args, "--assumeyes", "install")
	cmd.Args = append(cmd.Args, deps.Names(true)...)

	err := cmd.Run()

	fmtc.If(!quiet).NewLine()

	if err != nil {
		return fmt.Errorf("Installation finished with error")
	}

	return nil
}

// getPackageManager returns name of package manager
func getPackageManager() string {
	if fsutil.IsExist("/usr/bin/dnf") || fsutil.IsExist("/bin/dnf") {
		return "dnf"
	}

	return "yum"
}

// genPackageManagerOptions translates utility options to package manager
// options
func genPackageManagerOptions() []string {
	var result []string

	for _, repo := range options.Split(OPT_ENABLEREPO) {
		result = append(result, "--enablerepo="+repo)
	}

	for _, repo := range options.Split(OPT_DISABLEREPO) {
		result = append(result, "--disablerepo="+repo)
	}

	for _, macro := range options.Split(OPT_DEFINE) {
		result = append(result, "--define="+macro)
	}

	return result
}
