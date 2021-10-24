# Configuration
Location of configuration files:
- Outside of Docker container: `/var/lib/lenpaste/`
- Inside Docker container: `/app/data/`

## config.json
This documentation is too big so it was put in a separate file: [`config.json.md`](https://github.com/lcomrade/lenpaste/blob/main/docs/self-hosting/config.json.md).

## about.txt
Here you can write information about your server.
It will be displayed on the main page.

## rules.txt
In this file you can describe the rules for using your server.
It will be available at `/rules`.

## robots.txt
As the name suggests, this file is a `robots.txt` file.
If this file does not exist, the following contents will be used:

```
User-agent: *
Disallow: /
```
