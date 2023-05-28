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

package handler

import (
	"git.lcomrade.su/root/lenpaste/internal/model"

	"github.com/gin-gonic/gin"
)

func (hand *handler) robotsTxtHand(c *gin.Context) {
	// Generate robots.txt
	robotsTxt := "User-agent: *\nDisallow: /\n"

	if !hand.cfg.Public.RobotsDisallow {
		proto := getProtocol(c)
		host := getHost(c)

		robotsTxt = "User-agent: *\nAllow: /\nSitemap: " + proto + "://" + host + "/sitemap.xml\n"
	}

	// Write response
	c.Data(200, "text/plain; charset=utf-8", []byte(robotsTxt))
}

func (hand *handler) sitemapHand(c *gin.Context) {
	if hand.cfg.Public.RobotsDisallow {
		hand.writeErrorWeb(c, model.ErrNotFound)
		return
	}

	// Get protocol and host
	proto := getProtocol(c)
	host := getHost(c)

	// Generate sitemap.xml
	sitemapXML := `<?xml version="1.0" encoding="UTF-8"?>`
	sitemapXML = sitemapXML + "\n" + `<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n"
	sitemapXML = sitemapXML + "<url><loc>" + proto + "://" + host + "/" + "</loc></url>\n"
	sitemapXML = sitemapXML + "<url><loc>" + proto + "://" + host + "/about" + "</loc></url>\n"
	sitemapXML = sitemapXML + "<url><loc>" + proto + "://" + host + "/docs/apiv1" + "</loc></url>\n"
	sitemapXML = sitemapXML + "<url><loc>" + proto + "://" + host + "/docs/api_libs" + "</loc></url>\n"
	sitemapXML = sitemapXML + "</urlset>\n"

	// Write response
	c.Data(200, "text/plain; charset=utf-8", []byte(sitemapXML))
}
