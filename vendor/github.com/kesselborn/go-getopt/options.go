// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"fmt"
	"os"
	"path/filepath"
)

type Definitions []Option
type Description string

type Options struct {
	Description Description
	Definitions Definitions
}

func (optionsDefinition Options) setEnvAndConfigValues(options map[string]OptionValue, environment map[string]string) (err *GetOptError) {
	significantEnvVars := make(map[string]Option)

	for _, opt := range optionsDefinition.Definitions {
		if value := opt.EnvVar(); value != "" {
			significantEnvVars[value] = opt
		}
	}

	for key, significantEnvVar := range significantEnvVars {
		if value := environment[key]; value != "" {
			options[significantEnvVar.Key()], err = assignValue(significantEnvVar.DefaultValue, value)
			if err != nil {
				break
			}
		}
	}

	return
}

func checkOptionsDefinitionConsistency(optionsDefinition Options) (err *GetOptError) {
	foundOptionalArg := false
	shortOpts := make(map[string]bool, len(optionsDefinition.Definitions))
	longOpts := make(map[string]bool, len(optionsDefinition.Definitions))
	envVars := make(map[string]bool, len(optionsDefinition.Definitions))

	for _, option := range optionsDefinition.Definitions {
		optionString := fmt.Sprintf("%#v", option)
		consistencyErrorPrefix := optionString + " wrong getopt usage: "

		if option.HasLongOpt() {
			longOpt := option.LongOpt()
			if _, present := longOpts[longOpt]; present {
				err = &GetOptError{ConsistencyError, consistencyErrorPrefix + " long opt '" + longOpt + "' already used in other option"}
			} else {
				longOpts[longOpt] = true
			}
		}

		if option.HasShortOpt() {
			shortOpt := option.ShortOpt()
			if _, present := shortOpts[shortOpt]; present {
				err = &GetOptError{ConsistencyError, consistencyErrorPrefix + " short opt '" + shortOpt + "' already used in other option"}
			} else {
				shortOpts[shortOpt] = true
			}
		}

		if option.HasEnvVar() {
			envVar := option.EnvVar()
			if _, present := envVars[envVar]; present {
				err = &GetOptError{ConsistencyError, consistencyErrorPrefix + " environment variable '" + envVar + "' already used in other option"}
			} else {
				envVars[envVar] = true
			}
		}

		switch {
		case option.Flags&IsArg > 0 && option.Flags&Required == 0 && option.Flags&Optional == 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an argument must be explicitly set to be Optional or Required"}
		case option.Flags&IsArg > 0 && option.Flags&Optional > 0:
			foundOptionalArg = true
		case option.Flags&IsArg > 0 && option.Flags&Required > 0 && foundOptionalArg:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "a required argument can't come after an optional argument"}
		case option.Flags&Optional > 0 && option.Flags&Required > 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an option can not be Required and Optional"}
		case option.Flags&Flag > 0 && option.Flags&ExampleIsDefault > 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an option can not be a Flag and have ExampleIsDefault"}
		case option.Flags&Required > 0 && option.Flags&ExampleIsDefault > 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an option can not be Required and have ExampleIsDefault"}
		case option.Flags&NoLongOpt > 0 && !option.HasShortOpt() && option.Flags&IsArg == 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an option must have either NoLongOpt or a ShortOption"}
		case option.Flags&Flag > 0 && option.Flags&IsArg > 0:
			err = &GetOptError{ConsistencyError, consistencyErrorPrefix + "an option can not be a Flag and be an argument (IsArg)"}
		}
	}

	return
}

func (options Options) FindOption(optionString string) (option Option, found bool) {
	for _, cur := range options.Definitions {
		if cur.ShortOpt() == optionString || cur.LongOpt() == optionString {
			option = cur
			found = true
			break
		}
	}

	return option, found
}

func (options Options) IsOptional(optionName string) (isRequired bool) {
	if option, found := options.FindOption(optionName); found && option.Flags&Optional != 0 {
		isRequired = true
	}

	return isRequired
}

func (options Options) IsRequired(optionName string) (isRequired bool) {
	if option, found := options.FindOption(optionName); found && option.Flags&Required != 0 {
		isRequired = true
	}

	return isRequired
}

func (options Options) IsFlag(optionName string) (isFlag bool) {
	if option, found := options.FindOption(optionName); found && option.Flags&Flag != 0 {
		isFlag = true
	}

	return isFlag
}

func (options Options) ConfigOptionKey() (key string) {
	for _, option := range options.Definitions {
		if option.Flags&IsConfigFile > 0 {
			key = option.Key()
			break
		}
	}

	return
}

func (options Options) RequiredArguments() (requiredOptions Options) {
	for _, cur := range options.Definitions {
		if (cur.Flags&Required != 0 && cur.Flags&IsArg != 0) || cur.Flags&IsSubCommand != 0 {
			requiredOptions.Definitions = append(requiredOptions.Definitions, cur)
		}
	}

	return
}

func (options Options) RequiredOptions() (requiredOptions []string) {
	for _, cur := range options.Definitions {
		if cur.Flags&Required != 0 && cur.Flags&IsArg == 0 && cur.Flags&IsPassThrough == 0 {
			requiredOptions = append(requiredOptions, cur.LongOpt())
		}
	}

	return
}

func (options Options) commandDefinition(arg0 string) (output string) {
	output = arg0

	passThroughSeparatorPrinted := false
	for _, option := range options.Definitions {
		if option.Flags&IsPassThrough > 0 && !passThroughSeparatorPrinted {
			output = output + " --"
			passThroughSeparatorPrinted = true
		}

		output = output + " " + option.Usage()
	}

	return
}

func (options Options) UsageCustomArg0(arg0 string) (output string) {
	return "Usage: " + options.commandDefinition(arg0) + "\n\n"
}

func (options Options) Usage() (output string) {
	return options.UsageCustomArg0(filepath.Base(os.Args[0]))
}

func (options Options) Help() (output string) {
	return options.HelpCustomArg0(filepath.Base(os.Args[0]))
}

func (options Options) calculateLongOptTextLenght() (length int) {
	for _, option := range options.Definitions {
		if curLength := len(option.LongOptString()); curLength > length {
			length = curLength
		}
	}

	length = length + 2

	return
}

func (options Options) HelpCustomArg0(arg0 string) (output string) {
	output = options.UsageCustomArg0(arg0)
	if options.Description != "" {
		output = output + string(options.Description) + "\n\n"
	}

	longOptTextLength := options.calculateLongOptTextLenght()

	var argumentsString string
	var optionsString string
	var passThroughString string

	usageOpt, helpOpt := options.usageHelpOptionNames()

	for _, option := range options.Definitions {
		switch {
		case option.Flags&IsSubCommand > 0:
			continue
		case option.Flags&IsPassThrough > 0:
			passThroughString = passThroughString + option.HelpText(longOptTextLength) + "\n"
		case option.Flags&IsArg > 0:
			argumentsString = argumentsString + option.HelpText(longOptTextLength) + "\n"
		case option.LongOpt() != helpOpt:
			optionsString = optionsString + option.HelpText(longOptTextLength) + "\n"
		}
	}

	if optionsString != "" {
		helpHelp := fmt.Sprintf("usage (-%s) / detailed help text (--%s)", usageOpt, helpOpt)

		if option, found := options.FindOption(helpOpt); found {
			helpHelp = option.Description
		}

		usageHelpOption := Option{fmt.Sprintf("%s|%s", helpOpt, usageOpt),
			helpHelp,
			Usage | Help | Flag, ""}
		optionsString = optionsString + usageHelpOption.HelpText(longOptTextLength) + "\n"
		output = output + "Options:\n" + optionsString + "\n"
	}

	if argumentsString != "" {
		output = output + "Arguments:\n" + argumentsString + "\n"
	}

	if passThroughString != "" {
		output = output + "Pass through arguments:\n" + passThroughString + "\n"
	}

	return
}
