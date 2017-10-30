// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import "os"

const (
	InvalidOption = iota
	MissingValue
	InvalidValue
	MissingOption
	OptionValueError
	ConsistencyError
	ConfigFileNotFound
	ConfigParsed
	MissingArgument
	NoSubCommand
	NoScope
	UnknownSubCommand
	UnknownScope
)

const OPTIONS_SEPARATOR = "--"

type GetOptError struct {
	ErrorCode int
	message   string
}

func (err *GetOptError) Error() (message string) {
	return err.message
}

func (optionsDefinition Options) usageHelpOptionNames() (shortOpt string, longOpt string) {
	shortOpt = "h"
	longOpt = "help"

	for _, option := range optionsDefinition.Definitions {
		if option.Flags&Usage > 0 {
			shortOpt = option.ShortOpt()
		}
		if option.Flags&Help > 0 {
			longOpt = option.LongOpt()
		}
	}

	return
}

// todo: method signature sucks
func (optionsDefinition Options) checkForHelpOrUsage(args []string, usageString string, helpString string) (wantsHelp bool, wantsUsage bool) {
	for _, arg := range args {
		switch arg {
		case usageString:
			wantsUsage = true
		case helpString:
			wantsHelp = true
		case OPTIONS_SEPARATOR:
			goto allOptsParsed
		}
	}
allOptsParsed:

	return
}

func (optionsDefinition Options) ParseCommandLine() (options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {
	return optionsDefinition.parseCommandLineImpl(os.Args[1:], mapifyEnvironment(os.Environ()), 0)
}

func (optionsDefinition Options) parseCommandLineImpl(args []string, environment map[string]string, flags int) (options map[string]OptionValue, arguments []string, passThrough []string, err *GetOptError) {

	if err = checkOptionsDefinitionConsistency(optionsDefinition); err == nil {
		options = make(map[string]OptionValue)
		arguments = make([]string, 0)

		for _, option := range optionsDefinition.Definitions {
			switch {
			case option.Flags&Flag != 0: // all flags are false by default
				options[option.Key()], err = assignValue(false, "false")
			case option.Flags&ExampleIsDefault != 0: // set default
				var newOptionValue OptionValue
				newOptionValue, err = assign(option.DefaultValue)
				newOptionValue.Set = false
				options[option.Key()] = newOptionValue
			}
		}

		usageString, helpString := optionsDefinition.usageHelpOptionNames()

		wantsHelp, wantsUsage := optionsDefinition.checkForHelpOrUsage(args, "-"+usageString, "--"+helpString)

		if err == nil {
			err = optionsDefinition.setEnvAndConfigValues(options, environment)

			for i := 0; i < len(args) && err == nil; i++ {

				var opt, val string
				var found bool

				token := args[i]

				if argumentsEnd(token) {
					passThrough = args[i+1:]
					break
				}

				if isValue(token) {
					arguments = append(arguments, token)
					continue
				}

				opt, val, found = parseShortOpt(token)

				if found {
					buffer := token

					for found && optionsDefinition.IsFlag(opt) && len(buffer) > 2 {
						// concatenated options ... continue parsing
						currentOption, _ := optionsDefinition.FindOption(opt)
						key := currentOption.Key()

						options[key], err = assignValue(currentOption.DefaultValue, "true")

						// make it look as if we have a normal option with a '-' prefix
						buffer = "-" + buffer[2:]
						opt, val, found = parseShortOpt(buffer)
					}

				} else {
					opt, val, found = parseLongOpt(token)
				}

				currentOption, found := optionsDefinition.FindOption(opt)
				key := currentOption.Key()

				if !found {
					err = &GetOptError{InvalidOption, "invalid option '" + token + "'"}
					break
				}

				if optionsDefinition.IsFlag(opt) {
					options[key], err = assignValue(true, "true")
				} else {
					if val == "" {
						if len(args) > i+1 && isValue(args[i+1]) {
							i = i + 1
							val = args[i]
						} else {
							err = &GetOptError{MissingValue, "Option '" + token + "' needs a value"}
							break
						}
					}

					if !isValue(val) {
						err = &GetOptError{InvalidValue, "Option '" + token + "' got invalid value: '" + val + "'"}
						break
					}

					options[key], err = assignValue(currentOption.DefaultValue, val)
				}

			}
		}

		if configKey := optionsDefinition.ConfigOptionKey(); configKey != "" && flags&ConfigParsed == 0 {
			if option, found := options[configKey]; found {
				if environment, e := processConfigFile(option.String, environment); e == nil {
					return optionsDefinition.parseCommandLineImpl(args, environment, flags|ConfigParsed)
				} else if option.Set == true { // if config file had a default value, don't freak out
					err = e
				}
			}
		}

		if err == nil {
			for _, requiredOption := range optionsDefinition.RequiredOptions() {
				if options[requiredOption].Set == false {
					err = &GetOptError{MissingOption, "Option '" + requiredOption + "' is missing"}
					break
				}
			}

			requiredArguments := optionsDefinition.RequiredArguments()

			if numOfRequiredArguments := len(requiredArguments.Definitions); numOfRequiredArguments > len(arguments) {
				firstMissingArgumentName := requiredArguments.Definitions[len(arguments)].Key()
				err = &GetOptError{MissingArgument, "Missing required argument <" + firstMissingArgumentName + ">"}
			}
		}
		if wantsHelp {
			options[helpString], err = assignValue("", "help")
		}
		if wantsUsage {
			options[helpString], err = assignValue("", "usage")
		}
	}

	return
}
