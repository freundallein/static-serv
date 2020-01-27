# static-serv
[![Build Status](https://travis-ci.org/freundallein/static-serv.svg?branch=master)](https://travis-ci.org/freundallein/static-serv)
[![Go Report Card](https://goreportcard.com/badge/github.com/freundallein/static-serv)](https://goreportcard.com/report/github.com/freundallein/static-serv)

Static files http server.  

Allows to serve application static files separately.  
Or you can just serve static site.
Has simple caching^ resource expiration is 60 seconds, GC each 30 seconds.

## Configuration
Application supports configuration via environment variables:
```
export STATIC_ROOT=/storage
export PREFIX=/static
export PORT=8000
```
## Installation
### With docker  
```
$> docker pull freundallein/staticserv:latest
```
### With source
```
$> git clone git@github.com:freundallein/static-serv.git
$> cd static-serv
$> make build
```

## Usage
Docker-compose

```
version: "3.5"

networks:
  network:
    name: example-network
    driver: bridge

services:
  gateway:
    image: "traefik:latest"
    container_name: "gateway"
    command:
      - "--api.insecure=true"
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"

      - "--serversTransport.forwardingTimeouts.dialTimeout=5s"

      - "--entrypoints.web.address=:80"
      - "--entryPoints.web.transport.respondingTimeouts.readTimeout=5"
      - "--entryPoints.web.transport.respondingTimeouts.writeTimeout=5"
      - "--entryPoints.web.transport.respondingTimeouts.idleTimeout=5"
    networks: 
      - network
    ports:
      - 127.0.0.1:80:80
      - 127.0.0.1:8080:8080
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro

  staticserv:
    image: freundallein/staticserv:latest
    container_name: staticserv
    restart: always
    environment: 
      - PORT=8000
      - STATIC_ROOT=/storage
      - PREFIX=/static
    networks: 
      - network
    expose: 
      - 8000
    labels:
      - "traefik.enable=true"

      - "traefik.http.routers.static.rule=Host(`127.0.0.1`)&&PathPrefix(`/static`)"
      - "traefik.http.routers.static.entrypoints=web"

      - "traefik.http.services.static.loadbalancer.healthcheck.interval=1s"
      - "traefik.http.services.static.loadbalancer.healthcheck.path=/healthz"
      - "traefik.http.services.static.loadbalancer.healthcheck.timeout=5s"
    volumes:
      - static-files:/storage

  whoami:
    image: containous/whoami:latest
    container_name: whoami
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.whoami.rule=Host(`127.0.0.1`)"
      - "traefik.http.routers.whoami.entrypoints=web"
    networks: 
      - network

  bootstrap:
    image: freundallein/bootstrap:4.4.1
    container_name: bootstrap
    restart: always
    networks: 
      - network
    volumes:
      - static-files:/data

volumes:
    static-files:

```

All requests to `127.0.0.1/static/*` will be served by `staticserv` container, otherwise by `whoami`.

You can mount docker volume or host here:
```
volumes:
  - static-files:/storage
```
