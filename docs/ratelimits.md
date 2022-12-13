# Rate limiting
Rate limits protect your server from abuse.
The default settings are optimal.
**Don't change default settings if you don't understand what you're doing.**

Each individual IP address has its own rate limit.
Also the rate limits of different actions do not depend on each other.

| Acts on             | Measured in               | Default | Docker variable                 | CLI flag                |
|---------------------|---------------------------|---------|---------------------------------|-------------------------|
| New paste creation. | New pastes per 5 minute.  | `15`    | `LENPASTE_NEW_PASTES_PER_5MIN`  | `-new-pastes-per-5min`  |
| New paste creation. | New pastes per 15 minute. | `15`    | `LENPASTE_NEW_PASTES_PER_15MIN` | `-new-pastes-per-15min` |
| New paste creation. | New pastes per 1 hour.    | `40`    | `LENPASTE_NEW_PASTES_PER_1HOUR` | `-new-pastes-per-1hour` |
| Paste page get.     | Get pastes per 5 minute.  | `50`    | `LENPASTE_GET_PASTES_PER_5MIN`  | `-get-pastes-per-5min`  |
| Paste page get.     | Get pastes per 15 minute. | `100`   | `LENPASTE_GET_PASTES_PER_15MIN` | `-get-pastes-per-15min` |
| Paste page get.     | Get pastes per 1 hour.    | `500`   | `LENPASTE_GET_PASTES_PER_1HOUR` | `-get-pastes-per-1hour` |
