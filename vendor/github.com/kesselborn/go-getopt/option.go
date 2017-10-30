// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import "strings"

const (
	Required = 1 << iota
	Optional
	Flag
	NoLongOpt
	ExampleIsDefault
	IsArg
	Argument
	Usage
	Help
	IsPassThrough
	IsConfigFile
	NoEnvHelp
	IsSubCommand
)

type Option struct {
	OptionDefinition string
	Description      string
	Flags            int
	DefaultValue     interface{}
}

func (option Option) eq(other Option) bool {
	return option.OptionDefinition == other.OptionDefinition &&
		option.Description == other.Description &&
		option.Flags == other.Flags &&
		option.DefaultValue == other.DefaultValue
}

func (option Option) neq(other Option) bool {
	return !option.eq(other)
}

func (option Option) Key() (key string) {
	return strings.Split(option.OptionDefinition, "|")[0]
}

func (option Option) LongOpt() (longOpt string) {
	if option.Flags&NoLongOpt == 0 {
		longOpt = option.Key()
	}

	return longOpt
}

func (option Option) HasLongOpt() (result bool) {
	return option.LongOpt() != ""
}

func (option Option) ShortOpt() (shortOpt string) {
	token := strings.Split(option.OptionDefinition, "|")

	if len(token) > 1 {
		shortOpt = token[1]
	}

	return shortOpt
}

func (option Option) HasShortOpt() (result bool) {
	return option.ShortOpt() != ""
}

func (option Option) EnvVar() (envVar string) {
	token := strings.Split(option.OptionDefinition, "|")

	if len(token) > 2 {
		envVar = token[2]
	}

	return envVar
}

func (option Option) HasEnvVar() (result bool) {
	return option.EnvVar() != ""
}
