// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

type Scopes map[string]SubCommandOptions

type SubSubCommandOptions struct {
	Global Options
	Scopes Scopes
}

func (ssco SubSubCommandOptions) flattenToSubCommandOptions(scope string) (sco SubCommandOptions, err *GetOptError) {
	globalCommand := ssco.Global
	var present bool

	if sco, present = ssco.Scopes[scope]; present == true {
		//print("\n======= " + scope + "========\n"); print(strings.Replace(fmt.Sprintf("%#v", sco.SubCommands),  "getopt", "\n", -1)); print("\n================")
		sco.Global.Definitions = append(globalCommand.Definitions, sco.Global.Definitions...)
	} else {
		err = &GetOptError{UnknownSubCommand, "Unknown scope: " + scope}
	}

	return
}

func (ssco SubSubCommandOptions) flattenToOptions(scope string, subCommand string) (options Options, err *GetOptError) {
	if sco, err := ssco.flattenToSubCommandOptions(scope); err == nil {
		options, err = sco.flattenToOptions(subCommand)
	}

	return
}

func (ssco SubSubCommandOptions) findScope() (scope string, err *GetOptError) {
	options := ssco.Global
	scope = "*"

	_, arguments, _, _ := options.ParseCommandLine()

	if len(arguments) < 1 {
		err = &GetOptError{NoScope, "Couldn't find scope"}
	} else {
		scope = arguments[0]
		if _, present := ssco.Scopes[scope]; present != true {
			err = &GetOptError{UnknownScope, "Given scope '" + scope + "' not defined"}
		}
	}

	return
}

func (ssco SubSubCommandOptions) findScopeAndSubCommand() (scope string, subCommand string, err *GetOptError) {
	if scope, err = ssco.findScope(); err == nil {
		var sco SubCommandOptions

		if sco, err = ssco.flattenToSubCommandOptions(scope); err == nil {
			var arguments []string
			if _, arguments, _, _ = sco.Global.ParseCommandLine(); len(arguments) > 1 {
				subCommand = arguments[1]
				if _, present := sco.SubCommands[subCommand]; present != true {
					err = &GetOptError{UnknownSubCommand, "Given sub command '" + subCommand + "' not defined"}
				}
			} else {
				err = &GetOptError{NoSubCommand, "Couldn't find sub command"}
			}
		}
	}

	return
}

func (ssco SubSubCommandOptions) ParseCommandLine() (scope string, subCommand string, options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {
	var scopeScError *GetOptError
	var flattenedOptions Options

	scope, subCommand, scopeScError = ssco.findScopeAndSubCommand()

	switch {
	case subCommand == "":
		flattenedOptions = ssco.Global
	case scope == "":
		flattenedOptions = ssco.Global
	default:
		flattenedOptions, _ = ssco.flattenToOptions(scope, subCommand)
	}

	options, arguments, passThrough, err = flattenedOptions.ParseCommandLine()

	if len(arguments) > 2 {
		arguments = arguments[2:]
	}

	if scopeScError != nil && err == nil {
		err = scopeScError
	}

	return
}

func (ssco SubSubCommandOptions) Usage() (output string) {
	return ssco.UsageCustomArg0(filepath.Base(os.Args[0]))
}

func (ssco SubSubCommandOptions) UsageCustomArg0(arg0 string) (output string) {
	scope, subCommand, err := ssco.findScopeAndSubCommand()
	givenScope, foundScope := ssco.Scopes[scope]

	if (err != nil && (err.ErrorCode == UnknownScope || err.ErrorCode == NoScope)) || !foundScope {
		output = ssco.Global.UsageCustomArg0(arg0)
	} else {
		givenCommand, foundCommand := givenScope.SubCommands[subCommand]
		arg0 = arg0 + " " + scope

		if (err != nil && (err.ErrorCode == UnknownSubCommand || err.ErrorCode == NoSubCommand)) || !foundCommand {
			output = givenScope.Global.UsageCustomArg0(arg0)
		} else {
			output = givenCommand.UsageCustomArg0(arg0 + " " + subCommand)
		}
	}

	return
}

func (ssco SubSubCommandOptions) formatScopesHelp(formatLength int) (output string) {
	// TODO: centralize format strings
	fmtStr := fmt.Sprintf("    %%-%ds       %%s\n", formatLength)

	keys := make([]string, len(ssco.Scopes))
	i := 0

	for k := range ssco.Scopes {
		keys[i] = k
		i = i + 1
	}
	sort.Strings(keys)

	for _, key := range keys {
		output = output + fmt.Sprintf(fmtStr, key, ssco.Scopes[key].Global.Description)
	}
	output = output + "\n"

	return
}

func (ssco SubSubCommandOptions) Help() (output string) {
	return ssco.HelpCustomArg0(filepath.Base(os.Args[0]))
}

func (ssco SubSubCommandOptions) HelpCustomArg0(arg0 string) (output string) {
	scope, subCommand, err := ssco.findScopeAndSubCommand()
	givenScope, foundScope := ssco.Scopes[scope]

	if (err != nil && (err.ErrorCode == UnknownScope || err.ErrorCode == NoScope)) || !foundScope {
		output = ssco.Global.HelpCustomArg0(arg0)
		output = output + "Available scopes:\n" + ssco.formatScopesHelp(ssco.Global.calculateLongOptTextLenght())
	} else {
		if (err != nil && (err.ErrorCode == UnknownSubCommand || err.ErrorCode == NoSubCommand)) || !foundScope {
			output = givenScope.HelpCustomArg0(arg0 + " " + scope)
		} else {
			output = givenScope.SubCommands[subCommand].HelpCustomArg0(arg0 + " " + scope + " " + subCommand)
		}
	}

	return
}
