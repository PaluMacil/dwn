# danwolf.net

Website source for danwolf.net which is my personal homepage.

# Roadmap

This project will begin as a simple Go static webserver with barebones content and stick to a code and content layout that will make it relatively easy to migrate to Rhyvu once it is stable.

# Build

```
go build
```

# Run

If Windows, you'll need to add the `.exe` to the end.

```
.\dwn
```

# Install

To install as a service (assuming Ubuntu 16.04 and Systemd) you will need to first run `chmod u+x /path/to/dwn/dwn` and create the file below:

**/etc/systemd/system/dwn.service**:

```
[Unit]
Description=dwn service
StartLimitBurst=5

[Service]
WorkingDirectory=/path/to/dwn/
ExecStart=/path/to/dwn/dwn
Restart=always
RestartSec=3

[Install]
WantedBy=multi-user.target
```

Start it: `sudo systemctl start dwn`

Enable it to run at boot (should create a symlink in `/etc/systemd/system/multi-user.target.wants/`): `sudo systemctl enable dwn`

Stop it: `sudo systemctl stop dwn`

Soft reload Systemd dependencies: `sudo systemctl daemon-reload`

# License

This project is available under an MIT license. Dependencies might carry other licenses.