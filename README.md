## Project
It's a Golang app that expose HTTP API and Dashboard to manage DNS Records for RFC 2136 compatible DNS Server.
Tested with KnotDNS server for personal use.

![screen1](https://github.com/dgoujard/rfc2136_web_api/blob/master/screenshot/screen1.png?raw=true)

![screen2](https://github.com/dgoujard/rfc2136_web_api/blob/master/screenshot/screen2.png?raw=true)


## Features
- Rest API (List/Create/Delete records) (/api endpoint)
- Swagger API Doc (/swagger endpoint, [Doc](https://petstore.swagger.io/?url=https://raw.githubusercontent.com/dgoujard/rfc2136_web_api/master/public/swagger/swagger.json))
- Web Dashboard in React (CRUD for A/AAAA/Cname/NS/MX records only) (/ endpoint)
- Dockerfile to compile the project (docker dir)

## TODO List
- Improve dashboard record type implementation
- Add fake endpoint to mimic dyndns and other cloud DNS services

## Developer

#### Env parameters
- DEBUG=true
- SERVER_PORT=8080
- SERVER_TIMEOUT_READ=5s
- SERVER_TIMEOUT_WRITE=10s
- SERVER_TIMEOUT_IDLE=15s
- SERVER_LOGIN=admin (Dashboard login)
- SERVER_PASSWORD=password (Dashboard password)

- RFC2136_HOST=ns1.dgoujard.network (DNS Server with RFC2136 enabled)
- RFC2136_PORT=53
- RFC2136_ZONE=dgoujard.network (zones you can remotly manage, multipe zones with coma separator)
- RFC2136_TSIG_SECRET=TODO
- RFC2136_TSIG_SECRET_ALG=hmac-sha256
- RFC2136_TSIG_KEYNAME=mykeyname.

#### To generate swagger file
- swagger generate spec -o public/swagger/swagger.json --scan-models
