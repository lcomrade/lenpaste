# Changelog
Semantic versioning is used (https://semver.org/).


## v1.1
- You can now specify author, author email and author URL for paste.
- Full localization into Russian.
- Added settings menu.
- Paste creation and expiration times are now displayed in the user's time zone.
- Add PostgreSQL DB support.


## v1.0
This is the first stable release of Lenpaste🎉

Compared to the previous unstable versions, everything has been drastically improved:
design, loading speed of the pages, API, work with the database.
Plus added syntax highlighting in the web interface.


## v0.2
Features:
- Paste title
- About server information
- Improved documentation
- Logging and log rotation
- Storage configuration
- Code optimization

Bug fixes:
- Added `./version.json` to Docker image
- Added paste expiration check before opening
- Fixed incorrect error of expired pastes
- API errors now return in JSON


## v0.1
Features:
- Alternative to pastebin.com
- Creating expiration pastes
- Web interface
- API
