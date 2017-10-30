// Copyright (c) 2011, SoundCloud Ltd., Daniel Bornkessel
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// Source code and contact info at http://github.com/kesselborn/go-getopt

package getopt

import "strings"

func parseShortOpt(option string) (opt string, val string, found bool) {
	if len(option) > 1 && option[0] == '-' && option[1] >= 'A' && option[1] <= 'z' {
		found = true
		opt = option[1:2]
		if len(option) > 2 {
			val = option[2:]
		}

	}

	return opt, val, found
}

func parseLongOpt(option string) (opt string, val string, found bool) {
	if len(option) > 3 && option[0:2] == "--" {
		found = true

		optTokens := strings.Split(option[2:], "=")

		opt = optTokens[0]

		if len(optTokens) > 1 {
			val = optTokens[1]
		}
	}

	return opt, val, found
}

func isValue(option string) bool {
	_, _, isShortOpt := parseShortOpt(option)
	_, _, isLongOpt := parseLongOpt(option)

	return !isShortOpt && !isLongOpt && !argumentsEnd(option)
}

func argumentsEnd(option string) bool {
	return option == OPTIONS_SEPARATOR
}
