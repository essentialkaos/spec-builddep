package app

// ////////////////////////////////////////////////////////////////////////////////// //
//                                                                                    //
//                         Copyright (c) 2025 ESSENTIAL KAOS                          //
//      Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>     //
//                                                                                    //
// ////////////////////////////////////////////////////////////////////////////////// //

import (
	"fmt"
	"os"

	"github.com/essentialkaos/ek/v13/env"
	"github.com/essentialkaos/ek/v13/fmtc"
	"github.com/essentialkaos/ek/v13/fsutil"
	"github.com/essentialkaos/ek/v13/options"
	"github.com/essentialkaos/ek/v13/support"
	"github.com/essentialkaos/ek/v13/support/deps"
	"github.com/essentialkaos/ek/v13/support/pkgs"
	"github.com/essentialkaos/ek/v13/terminal"
	"github.com/essentialkaos/ek/v13/terminal/tty"
	"github.com/essentialkaos/ek/v13/usage"
	"github.com/essentialkaos/ek/v13/usage/completion/bash"
	"github.com/essentialkaos/ek/v13/usage/completion/fish"
	"github.com/essentialkaos/ek/v13/usage/completion/zsh"
	"github.com/essentialkaos/ek/v13/usage/man"
	"github.com/essentialkaos/ek/v13/usage/update"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Basic utility info
const (
	APP  = "spec-builddep"
	VER  = "1.1.0"
	DESC = "Utility for installing dependencies for building an RPM package"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// Options
const (
	OPT_LIST     = "L:list"
	OPT_ACTUAL   = "A:actual"
	OPT_CLEAN    = "C:clean"
	OPT_NO_COLOR = "nc:no-color"
	OPT_HELP     = "h:help"
	OPT_VER      = "v:version"

	OPT_QUIET       = "q:quiet"
	OPT_DEFINE      = "D:define"
	OPT_EXCLUDE     = "x:exclude"
	OPT_ENABLEREPO  = "ER:enablerepo"
	OPT_DISABLEREPO = "DR:disablerepo"

	OPT_VERB_VER     = "vv:verbose-version"
	OPT_COMPLETION   = "completion"
	OPT_GENERATE_MAN = "generate-man"
)

// ////////////////////////////////////////////////////////////////////////////////// //

// optMap contains information about all supported options
var optMap = options.Map{
	OPT_LIST:     {Type: options.BOOL},
	OPT_ACTUAL:   {Type: options.BOOL},
	OPT_CLEAN:    {Type: options.BOOL},
	OPT_NO_COLOR: {Type: options.BOOL},
	OPT_HELP:     {Type: options.BOOL},
	OPT_VER:      {Type: options.MIXED},

	OPT_QUIET:       {Type: options.BOOL},
	OPT_DEFINE:      {Mergeble: true},
	OPT_EXCLUDE:     {Mergeble: true},
	OPT_ENABLEREPO:  {Mergeble: true},
	OPT_DISABLEREPO: {Mergeble: true},

	OPT_VERB_VER:     {Type: options.BOOL},
	OPT_COMPLETION:   {},
	OPT_GENERATE_MAN: {Type: options.BOOL},
}

// useRawOutput is raw output flag (for cli command)
var useRawOutput = false

// ////////////////////////////////////////////////////////////////////////////////// //

// Run is main utility function
func Run(gitRev string, gomod []byte) {
	preConfigureUI()

	options.MergeSymbol = "\x00"
	args, errs := options.Parse(optMap)

	if !errs.IsEmpty() {
		terminal.Error("Options parsing errors:")
		terminal.Error(errs.Error("- "))
		os.Exit(1)
	}

	configureUI()

	switch {
	case options.Has(OPT_COMPLETION):
		os.Exit(printCompletion())
	case options.Has(OPT_GENERATE_MAN):
		printMan()
		os.Exit(0)
	case options.GetB(OPT_VER):
		genAbout(gitRev).Print(options.GetS(OPT_VER))
		os.Exit(0)
	case options.GetB(OPT_VERB_VER):
		support.Collect(APP, VER).
			WithRevision(gitRev).
			WithDeps(deps.Extract(gomod)).
			WithPackages(pkgs.Collect("dnf,yum", "rpm", "rpm-build")).
			Print()
		os.Exit(0)
	case options.GetB(OPT_HELP) || len(args) == 0:
		genUsage().Print()
		os.Exit(0)
	}

	err := checkSystem()

	if err != nil {
		if !options.GetB(OPT_QUIET) {
			terminal.Error(err)
		}

		os.Exit(1)
	}

	err = process(args)

	if err != nil {
		if !options.GetB(OPT_QUIET) {
			terminal.Error(err)
		}

		os.Exit(1)
	}
}

// ////////////////////////////////////////////////////////////////////////////////// //

// preConfigureUI preconfigures UI based on information about user terminal
func preConfigureUI() {
	if !tty.IsTTY() {
		fmtc.DisableColors = true
		useRawOutput = true
	}
}

// configureUI configures user interface
func configureUI() {
	if options.GetB(OPT_NO_COLOR) {
		fmtc.DisableColors = true
	}
}

// checkSystem checks system for required dependencies
func checkSystem() error {
	if env.Which("rpm") == "" {
		return fmt.Errorf("This utility requires rpm")
	}

	if env.Which("rpmspec") == "" {
		return fmt.Errorf("This utility requires rpmspec (part of rpm-build)")
	}

	return nil
}

// process starts arguments processing
func process(args options.Arguments) error {
	specFile := args.Get(0).Clean().String()
	err := fsutil.ValidatePerms("FRS", specFile)

	if err != nil {
		return err
	}

	switch {
	case options.GetB(OPT_LIST):
		err = listDeps(specFile)
	default:
		err = installDeps(specFile)
	}

	return err
}

// ////////////////////////////////////////////////////////////////////////////////// //

// printCompletion prints completion for given shell
func printCompletion() int {
	info := genUsage()

	switch options.GetS(OPT_COMPLETION) {
	case "bash":
		fmt.Print(bash.Generate(info, APP))
	case "fish":
		fmt.Print(fish.Generate(info, APP))
	case "zsh":
		fmt.Print(zsh.Generate(info, optMap, APP))
	default:
		return 1
	}

	return 0
}

// printMan prints man page
func printMan() {
	fmt.Println(man.Generate(genUsage(), genAbout("")))
}

// genUsage generates usage info
func genUsage() *usage.Info {
	info := usage.NewInfo("", "spec-file")

	info.AddOption(OPT_LIST, "List required build dependencies")
	info.AddOption(OPT_ACTUAL, "Install the latest versions of all packages")
	info.AddOption(OPT_CLEAN, "Clean package manager cache before install")
	info.AddOption(OPT_QUIET, "Quiet mode {s}(don't output anything){!}")
	info.AddOption(OPT_DEFINE, "Define a macro for spec file parsing {s-}(mergeble){!}", "macro")
	info.AddOption(OPT_EXCLUDE, "Exclude packages by name or glob {s-}(mergeble){!}", "package")
	info.AddOption(OPT_ENABLEREPO, "Enable additional repositories {s-}(mergeble){!}", "repo")
	info.AddOption(OPT_DISABLEREPO, "Disable repositories {s-}(mergeble){!}", "repo")
	info.AddOption(OPT_NO_COLOR, "Disable colors in output")
	info.AddOption(OPT_HELP, "Show this help message")
	info.AddOption(OPT_VER, "Show version")

	info.AddExample(
		"-L app.spec",
		"List all required build dependencies from the spec",
	)

	info.AddExample(
		"app.spec -ER epel-testing -ER kaos-testing",
		"Enable some repos and install required packages",
	)

	info.AddExample(
		"-A -C app.spec",
		"Install the latest version of required packages",
	)

	info.AddExample(
		"-D '_ver:3.1.8' -D '_git:aa7bbf8' app.spec",
		"Define some macro values and install required packages",
	)

	return info
}

// genAbout generates info about version
func genAbout(gitRev string) *usage.About {
	about := &usage.About{
		App:     APP,
		Version: VER,
		Desc:    DESC,
		Year:    2009,
		Owner:   "ESSENTIAL KAOS",
		License: "Apache License, Version 2.0 <https://www.apache.org/licenses/LICENSE-2.0>",
	}

	if gitRev != "" {
		about.Build = "git:" + gitRev
		about.UpdateChecker = usage.UpdateChecker{"essentialkaos/spec-builddep", update.GitHubChecker}
	}

	return about
}

// ////////////////////////////////////////////////////////////////////////////////// //
