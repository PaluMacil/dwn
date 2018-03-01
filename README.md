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

```toml
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

# Nginx

If serving with Nginx, you might need to have a configuration something like this:

```nginx
server {
        listen 80;
        listen [::]:80 ipv6only=on;
        server_name danwolf.net www.danwolf.net;
        rewrite ^/(.*) https://danwolf.net/$1 permanent;
}

server {
        listen 443 ssl http2;
        listen [::]:443 ipv6only=on ssl http2;
        charset utf-8;
        client_max_body_size 75M;
        server_name danwolf.net www.danwolf.net;
        ssl_certificate /etc/ssl/certs/danwolf.net.crt;
        ssl_certificate_key /etc/ssl/private/danwolf.net.key;
        error_log /var/log/nginx/error.log warn;
        add_header "X-UA-Compatible" "IE=Edge";
        server_tokens off;

        location / {
                proxy_cache off;
                proxy_pass http://127.0.0.1:3035;
        }
}
```

After updating nginx, you might need to reload the config with `sudo service nginx reload`.

# License

This project is available under an MIT license. Dependencies might carry other licenses.