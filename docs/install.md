# Installation
At the moment, **only** the **Docker container** is supported.

## About Docker tags
### stable
The stable version is based on the latest stable release.

### latest
Latest is generated based on the last comit in the main branch.
This release is for developers and beta testers only!


## Installation Manual
1. Installing Docker
```
apt install -y docker docker.io
```

2. Adding systemd unit
```
rm -f /etc/systemd/system/lenpaste.service
wget -O '/etc/systemd/system/lenpaste.service' 'https://raw.githubusercontent.com/lcomrade/lenpaste/main/init/lenpaste.stable.service'
systemctl daemon-reload
```

3. Enabling and starting systemd service
```
systemctl enable lenpaste
systemctl start lenpaste
```

4. All data and config files are now located in `/var/lib/lenpaste/`.
