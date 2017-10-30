// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import (
	"fmt"
	"strings"
)

func valueIze(originalValue string) (value string) {
	value = "<" + strings.ToLower(originalValue) + ">"

	return
}

func (option Option) HelpText(longOptLength int) (output string) {
	fmtStringLongAndShort := fmt.Sprintf("    -%%-1s, --%%-%ds %%s", longOptLength) // "p", "port=PORT", "the port that should be used"
	fmtStringShort := fmt.Sprintf("    -%%-1s%%-%ds    %%s", longOptLength)         // "p", "PORT", "the port that should be used"
	fmtStringLong := fmt.Sprintf("        --%%-%ds %%s", longOptLength)             // "port=PORT", "the port that should be used"
	fmtArgument := fmt.Sprintf("    %%-%ds       %%s", longOptLength)               // "port=PORT", "the port that should be used"

	if option.Description != "" {
		switch {
		case option.Flags&(IsArg|IsPassThrough) > 0:
			output = fmt.Sprintf(fmtArgument, valueIze(option.Key()), option.DescriptionText())
		case option.HasLongOpt() && option.HasShortOpt():
			output = fmt.Sprintf(fmtStringLongAndShort, option.ShortOptString(), option.LongOptString(), option.DescriptionText())
		case !option.HasLongOpt() && option.HasShortOpt() && option.Flags&Flag > 0:
			output = fmt.Sprintf(fmtStringShort, option.ShortOptString(), "", option.DescriptionText())
		case option.HasShortOpt():
			output = fmt.Sprintf(fmtStringShort, option.ShortOptString(), valueIze(option.Key()), option.DescriptionText())
		case option.HasLongOpt():
			output = fmt.Sprintf(fmtStringLong, option.LongOptString(), option.DescriptionText())
		}
	}

	return output
}

func (option Option) LongOptString() (longOptString string) {
	if option.HasLongOpt() {
		longOptString = option.LongOpt()

		if option.Flags&(Flag|Usage|Help|IsPassThrough) == 0 {
			longOptString = longOptString + "=" + valueIze(option.LongOpt())
		}
	}

	return
}

func (option Option) ShortOptString() (shortOptString string) {
	if option.HasShortOpt() {
		shortOptString = option.ShortOpt()

		if !option.HasLongOpt() && option.Flags&(Flag|Usage|Help) == 0 {
			shortOptString = shortOptString + " " + option.LongOpt()
		}

		if option.Flags&(Flag|NoLongOpt) == Flag|NoLongOpt {
			shortOptString = shortOptString + " "
		}
	}

	return
}

func (option Option) Usage() (usageString string) {
	switch {
	case option.Flags&(IsArg|IsPassThrough|IsSubCommand) > 0:
		usageString = valueIze(option.Key())
	case option.HasShortOpt():
		usageString = "-" + option.ShortOpt()
	default:
		usageString = "--" + option.LongOpt()
	}

	if option.Flags&(Flag|IsArg|IsSubCommand|IsPassThrough|Usage|Help) == 0 {
		var separator string
		if option.HasShortOpt() {
			separator = " "
		} else {
			separator = "="
		}

		usageString = usageString + separator + valueIze(option.Key())
	}

	if option.Flags&Optional > 0 || option.Flags&Help > 0 || option.Flags&Usage > 0 || option.Flags&ExampleIsDefault > 0 {
		usageString = "[" + usageString + "]"
	}

	return
}

func (option Option) DescriptionText() (description string) {
	description = option.Description

	defaultValue := fmt.Sprintf("%v", option.DefaultValue)
	// %v for arrays prints something like [3 4 5] ... let's transform that to 3,4,5:
	defaultValue = strings.Replace(strings.Replace(strings.Replace(defaultValue, "[", "", -1), "]", "", -1), " ", ",", -1)

	if defaultValue != "" && option.DefaultValue != nil {
		switch {
		case option.Flags&(Optional|IsConfigFile) > 0 && option.Flags&ExampleIsDefault > 0:
			description = description + " (default: " + defaultValue + ")"
		case option.Flags&Required > 0 || option.Flags&Optional > 0 || option.Flags&IsArg > 0:
			description = description + " (e.g. " + defaultValue + ")"
		}
	}

	if option.HasEnvVar() && option.Flags&NoEnvHelp == 0 {
		description = description + "; setable via $" + option.EnvVar()
	}

	return
}
