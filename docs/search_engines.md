# Add Lenpaste to Search Engines
Lenpaste can be added to any engine that supports SiteMap and `robots.txt`.
For example: Google, Yandex, Bing.

Make sure that the `-robots-disallow` cmd flag is not present.
Or if you use Docker, set `LENPASTE_ROBOTS_DISALLOW=false`.
If you have done everything correctly, you should see the following:
- File `<SERVER>/robots.txt` must contain the directive `Allow: /`.
- Also must be present file `<SERVER>/sitemap.xml`.

After the above steps, you can add the site to search engines.

PS: It is also advisable to add your server to [Lenmonitor](https://monitor.lcomrade.su/add?srv=lenpaste) - Lenpaste server aggregator.
