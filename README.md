# ![SPA](favicon.svg) SPA server
![Go](https://github.com/D-Haven/spa-server/workflows/Go/badge.svg)
![CodeQL](https://github.com/D-Haven/spa-server/workflows/CodeQL/badge.svg)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=D-Haven_spa-server&metric=security_rating)](https://sonarcloud.io/dashboard?id=D-Haven_spa-server)
[![Vulnerabilities](https://sonarcloud.io/api/project_badges/measure?project=D-Haven_spa-server&metric=vulnerabilities)](https://sonarcloud.io/dashboard?id=D-Haven_spa-server)

Simple HTTP/2 server for static single page apps, complete with auto-redirect to the
home page to handle front-end routing.  Designed for a microservices based
architecture, you can build up your deployment in a couple ways.

## Extend the Docker Image

```Dockerfile
FROM spa-server:latest

COPY config.yaml /
COPY server.crt /
COPY server.key /
ADD www /www/
```
