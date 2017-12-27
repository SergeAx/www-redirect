Eternal redirect from any domain to www. subdomain
==================================================

You should use www. subdomain to serve website from any second-level domain.

You also need to redirect traffic from second-level domain there.

Adding `return 301 $scheme://$host$request_uri` to every nginx config
is just plain boring. Learning a new programming language is fun. So
here we go[lang]!

Usage
=====

Add this container to Traefik network:

```docker-compose
version: "2.0"
services:
  www-redirect:
    container_name: www-redirect
    image: sergeax/eternal-www-redirect:latest
    restart: unless-stopped
    labels:
      traefik.enable: 'true'
      traefik.backend: www-redirect
      traefik.frontend.rule: HostRegexp:{*}
      traefik.priority: 999999
      traefik.port: '80'
      traefik.docker.network: traefik_default
    networks:
    - traefik_default

networks:
  traefik_default:
    external: true
```
