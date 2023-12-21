# Changelog
Semantic versioning is used (https://semver.org/).


## v1.3.1
- Fixed a problem with building Lenpaste from source code.
- Revised documentation.
- Minor improvements were made.

## v1.3
- UI: Added custom themes support. Added light theme.
- UI: Added translations into Bengali and German (thanks Pardesi_Cat and Hiajen).
- UI: Check boxes and spoilers now have a custom design.
- Admin: Added support for `X-Real-IP` header for reverse proxy.
- Admin: Added Server response header (for example: `Lenpaste/1.3`).
- Fix: many bugs and errors.
- Dev: Improved quality of `Dockerfile` and `entrypoint.sh`

## v1.2
- UI: Add history tab.
- UI: Add copy to clipboard button.
- Admin: Rate-limits on paste creation (`LENPASTE_NEW_PASTES_PER_5MIN` or `-new-pastes-per-5min`).
- Admin: Add terms of use support (`/data/terms` or `-server-terms`).
- Admin: Add default paste life time for WEB interface (`LENPASTE_UI_DEFAULT_LIFETIME` or `-ui-default-lifetime`).
- Admin: Private servers - password request to create paste (`/data/lenpasswd` or `-lenpasswd-file`).
- Fix: **Critical security fix!**
- Fix: not saving cookies.
- Fix: display language name in WEB.
- Fix: compatibility with WebKit (Gnome WEB).
- Dev: Drop Go 1.15 support. Update dependencies.


## v1.1.1
- Fixed: Incorrect operation of the maximum paste life parameter.
- Updated README.


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
