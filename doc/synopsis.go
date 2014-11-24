// antha/doc/synopsis.go: Part of the Antha language
// Copyright (C) 2014 The Antha authors. All rights reserved.
// 
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
// 
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
// 
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
// 
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o 
// Synthace Ltd. The London Bioscience Innovation Centre
// 1 Royal College St, London NW1 0NH UK


package doc

import (
	"strings"
	"unicode"
)

// firstSentenceLen returns the length of the first sentence in s.
// The sentence ends after the first period followed by space and
// not preceded by exactly one uppercase letter.
//
func firstSentenceLen(s string) int {
	var ppp, pp, p rune
	for i, q := range s {
		if q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}
		if q == ' ' && p == '.' && (!unicode.IsUpper(pp) || unicode.IsUpper(ppp)) {
			return i
		}
		if p == '。' || p == '．' {
			return i
		}
		ppp, pp, p = pp, p, q
	}
	return len(s)
}

const (
	keepNL = 1 << iota
)

// clean replaces each sequence of space, \n, \r, or \t characters
// with a single space and removes any trailing and leading spaces.
// If the keepNL flag is set, newline characters are passed through
// instead of being change to spaces.
func clean(s string, flags int) string {
	var b []byte
	p := byte(' ')
	for i := 0; i < len(s); i++ {
		q := s[i]
		if (flags&keepNL) == 0 && q == '\n' || q == '\r' || q == '\t' {
			q = ' '
		}
		if q != ' ' || p != ' ' {
			b = append(b, q)
			p = q
		}
	}
	// remove trailing blank, if any
	if n := len(b); n > 0 && p == ' ' {
		b = b[0 : n-1]
	}
	return string(b)
}

// Synopsis returns a cleaned version of the first sentence in s.
// That sentence ends after the first period followed by space and
// not preceded by exactly one uppercase letter. The result string
// has no \n, \r, or \t characters and uses only single spaces between
// words. If s starts with any of the IllegalPrefixes, the result
// is the empty string.
//
func Synopsis(s string) string {
	s = clean(s[0:firstSentenceLen(s)], 0)
	for _, prefix := range IllegalPrefixes {
		if strings.HasPrefix(strings.ToLower(s), prefix) {
			return ""
		}
	}
	return s
}

var IllegalPrefixes = []string{
	"copyright",
	"all rights",
	"author",
}