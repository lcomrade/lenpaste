// Copyright (C) 2021-2022 Leonid Maslakov.

// This file is part of Lenpaste.

// Lenpaste is free software: you can redistribute it
// and/or modify it under the terms of the
// GNU Affero Public License as published by the
// Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.

// Lenpaste is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY
// or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero Public License for more details.

// You should have received a copy of the GNU Affero Public License along with Lenpaste.
// If not, see <https://www.gnu.org/licenses/>.

package web

import (
	"bytes"
	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

func tryHighlight(source string, lexer string) string {
	// Determine lexer
	l := lexers.Get(lexer)
	if l == nil {
		return source
	}

	l = chroma.Coalesce(l)

	// Determine formatter
	f := html.New(
		html.Standalone(false),
		html.WithClasses(false),
		html.TabWidth(4),
		html.WithLineNumbers(true),
		html.WrapLongLines(true),
	)

	s := styles.Get("monokai")

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return source
	}

	// Format
	var buf bytes.Buffer

	err = f.Format(&buf, s, it)
	if err != nil {
		return source
	}

	return buf.String()
}
