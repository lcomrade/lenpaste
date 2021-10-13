# Installation
At the moment, **only** the **Docker container** is supported.


## Installation
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


## Docker
### Tags
`stable` - Based on the latest stable release.

`latest` - Based on the last comit in the main branch.
This release is for developers and beta testers only!

### Changing the port
**You cannot change the port in the config (`./data/config.json`)!**

To change the port, you need to change the startup parameters of the Docker container:
1. Open the file `/etc/systemd/system/lenpaste.service`.
2. Find the parameter `ExecStart=`.
3. Replace `-p 8000:8000/tcp` with `-p PORT_HERE:8000/tcp`.
4. Reload daemons `systemctl daemon-reload`.
5. Restart service `systemctl restart lenpaste`.
