// Copyright (C) 2021-2023 Leonid Maslakov.

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
	"crypto/md5"
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type jsTmpl struct {
	Translate func(string, ...interface{}) template.HTML
	Theme     func(string) string
}

func (hand *handler) styleCSSHand(c *gin.Context) {
	c.Header("Content-Type", "text/css; charset=utf-8")
	return data.styleCSS.Execute(rw, jsTmpl{
		Translate: hand.l10n.findLocale(req).translate,
		Theme:     data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).theme,
	})
}

func (hand *handler) mainJSHand(c *gin.Context) {
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	rw.Write(*data.mainJS)
	return nil
}

func (hand *handler) codeJSHand(c *gin.Context) {
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	return data.codeJS.Execute(rw, jsTmpl{Translate: hand.l10n.findLocale(req).translate})
}

func (hand *handler) historyJSHand(c *gin.Context) {
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	return data.historyJS.Execute(rw, jsTmpl{
		Translate: hand.l10n.findLocale(req).translate,
		Theme:     data.themes.findTheme(req, hand.cfg.UI.DefaultTheme).theme,
	})
}

func (hand *handler) pasteJSHand(c *gin.Context) {
	c.Header("Content-Type", "application/javascript; charset=utf-8")
	return data.pasteJS.Execute(rw, jsTmpl{Translate: hand.l10n.findLocale(req).translate})
}

func init() {
	resp := "\u0045\u0072\u0072\u006f\u0072\u002e\u0020\u0059\u006f\u0075\u0020\u006d\u0061"
	resp += "\u0079\u0020\u0062\u0065\u0020\u0076\u0069\u006f\u006c\u0061\u0074\u0069\u006e"
	resp += "\u0067\u0020\u0074\u0068\u0065\u0020\u0041\u0047\u0050\u004c\u0020\u0076\u0033"
	resp += "\u0020\u006c\u0069\u0063\u0065\u006e\u0073\u0065\u0021"

	tmp, err := embFS.ReadFile("data/base.tmpl")
	if err != nil {
		println("error:", err.Error())
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "<a href=\"/about\">{{ call .Translate `base.About` }}</a>") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/about.tmpl")
	if err != nil {
		println("\u0065\u0072\u0072\u006f\u0072\u003a", err.Error())
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "<p>{{call .Translate `about.LenpasteAuthors` `/about/authors`}}</p>") {
		println(resp)
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "/about/source_code") {
		println(resp)
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "/about/license") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/authors.tmpl")
	if err != nil {
		println("\u0065\u0072\u0072\u006f\u0072\u003a", err.Error())
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "<li>Leonid Maslakov (aka lcomrade) &lt<a href=\"mailto:root@lcomrade.su\">root@lcomrade.su</a>&gt - Core Developer.</li>") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/source_code.tmpl")
	if err != nil {
		println("\u0065\u0072\u0072\u006f\u0072\u003a", err.Error())
		os.Exit(1)
	}

	if !strings.Contains(string(tmp), "https://git.lcomrade.su/root/lenpaste") {
		println(resp)
		os.Exit(1)
	}

	tmp, err = embFS.ReadFile("data/license.tmpl")
	if err != nil {
		println("\u0065\u0072\u0072\u006f\u0072\u003a", err.Error())
		os.Exit(1)
	}

	if fmt.Sprintf("%x", md5.Sum(tmp)) != "a1d6dd7f4b7470be5197381b85ee4fb5" {
		println(resp)
		os.Exit(1)
	}
}
