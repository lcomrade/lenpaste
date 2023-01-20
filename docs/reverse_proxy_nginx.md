# Reverse proxy: Nginx
Here is an example of the basic Nginx config (`/etc/nginx/nginx.conf`):
```nginx
events {
	worker_connections 1024;
}

http {
	error_log /var/log/nginx/error.log warn;
	server_tokens off; # Disables emitting nginx version on error pages and in the “Server” response header field.

	ssl_protocols TLSv1.2 TLSv1.3; # TLSv1.2 enables HTTPS on older devices.

	client_max_body_size 1M;
	client_body_timeout 300s;

	proxy_http_version 1.1;
	
	# Required for Lenpaste to work correctly.
	proxy_set_header Host $host;
	proxy_set_header X-Forwarded-For $remote_addr;
	proxy_set_header X-Forwarded-Proto $scheme;
}

# HTTP
server {
	server_name YOUR_DOMAIN;
	listen 80;
	listen [::]:80;

	access_log /var/log/nginx/YOUR_DOMAIN.access.log combined;

	location / {
		proxy_pass http://localhost:8000/;
		#return 301 https://$host$request_uri; - redirect to HTTPS
	}

	# Required for Lets Encrypt
	location /.well-known/acme-challenge/ {
		root /var/www/letsencrypt/;
	}
}

# HTTPS
server {
	server_name YOUR_DOMAIN;
	listen 443 ssl http2;
	listen [::]:443 ssl http2;
	ssl_certificate /etc/letsencrypt/live/YOUR_DOMAIN/fullchain.pem;
	ssl_certificate_key /etc/letsencrypt/live/YOUR_DOMAIN/privkey.pem;

	access_log /var/log/nginx/YOUR_DOMAIN.access.log combined;
	
	location / {
		proxy_pass http://localhost:8000/;
	}
}
```
