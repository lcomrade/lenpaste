// Copyright 2022 Leonid Maslakov. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package lineend

import (
	"strings"
)

// GetLineEnd allows you to get the end of a line used in the text.
// It can return \r\n, \r, \n or an empty string.
func GetLineEnd(text string) string {
	dos := strings.Count(text, "\r\n")
	oldMac := strings.Count(text, "\r")
	unix := strings.Count(text, "\n")

	if dos == 0 && oldMac == 0 && unix == 0 {
		return ""
	}

	if dos >= oldMac && dos >= unix {
		return "\r\n"
	}

	if oldMac >= dos && oldMac >= unix {
		return "\r"
	}

	if unix >= dos && unix >= oldMac {
		return "\n"
	}

	return ""
}

// DosToOldMac DosToOldMac сonverts end of line from CRLF to CR.
func DosToOldMac(text string) string {
	return strings.Replace(text, "\n", "", -1)
}

// DosToUnix сonverts end of line from CRLF to LF.
func DosToUnix(text string) string {
	return strings.Replace(text, "\r", "", -1)
}

// OldMacToDos сonverts end of line from CR to CRLF.
func OldMacToDos(text string) string {
	return strings.Replace(text, "\r", "\r\n", -1)
}

// OldMacToUnix сonverts end of line from CR to LF.
func OldMacToUnix(text string) string {
	return strings.Replace(text, "\r", "\n", -1)
}

// UnixToDos сonverts end of line from LF to CRLF.
func UnixToDos(text string) string {
	return strings.Replace(text, "\n", "\r\n", -1)
}

// UnixToOldMac сonverts end of line from LF to CR.
func UnixToOldMac(text string) string {
	return strings.Replace(text, "\n", "\r", -1)
}

// UnknownToDos converts unknown line end to CRLF.
func UnknownToDos(text string) string {
	switch GetLineEnd(text) {
	case "\r":
		return OldMacToDos(text)

	case "\n":
		return UnixToDos(text)
	}

	return text
}

// UnknownToOldMac converts unknown line end to CR.
func UnknownToOldMac(text string) string {
	switch GetLineEnd(text) {
	case "\r\n":
		return DosToOldMac(text)

	case "\n":
		return UnixToOldMac(text)
	}

	return text
}

// UnknownToUnix converts unknown line end to LF.
func UnknownToUnix(text string) string {
	switch GetLineEnd(text) {
	case "\r\n":
		return DosToUnix(text)

	case "\r":
		return OldMacToUnix(text)
	}

	return text
}
