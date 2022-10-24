# Rate limiting
Rate limits protect your server from abuse.
The default settings are optimal.
**Don't change default settings if you don't understand what you're doing.**

Each individual IP address has its own rate limit.
Also the rate limits of different actions do not depend on each other.

| Acts on             | Measured in              | Default | Docker variable                | CLI flag              |
|---------------------|--------------------------|---------|--------------------------------|-----------------------|
| New paste creation. | New pastes per 5 minute. | `15`    | `LENPASTE_NEW_PASTES_PER_5MIN` | `-new-pastes-per-5min`|
