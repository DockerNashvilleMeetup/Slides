class: top, rainbow-bg
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain

# About Me

* Organizer for Nashville Docker and Nashville Go Meetups
* Middle Tennessee Native
* Senior Innovation Engineer @ Franklin American Mortgage Company
* Docker user since 2014; putting Docker systems in production since 2015
* I enjoy cycling, competitive shooting, and board games

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Agenda

* What are layers? How do they work?
* Tips for building minimal images
* Dockerfiles: The Past and The Future
* Docker Security Scan
* Q&A + Docker Nashville Meetups

---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## What are Layers?

<table>
<tbody>
<tr><td style="vertical-align: top;">
<ul>
<li>Layers are composed and then digested</li>
<li>Each digested layer is given a signature</li>
<li>Signatures govern integrity and security</li>
<li>More layers, more problems</li>
</ul>
</td>
<td>
<img style="margin-top: -100px;" src="./slides/images/docker-image-comp.png">
</td>
</tr>
</tbody>
</table>

???
* images reference one or more layers which eventually contribute to a containers filesystem
* building images locally will produce intermediate "cache" layers which will appear in a `docker history`
* these cache layers are used when doing subsequent builds but will not be sent to a registry via push; so if you run a build on one machine and then pull on another machine and run the build there - it will not cache any of those layers


---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Reduce your layers with this one weird trick
(okay, maybe a few tricks)

* Use common base images when possible
* Reduce the amount of data written to the container layer
* Chain your RUN statements
* Defer cache misses until the last mile
* Use multi-stage builds!

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Past

```
FROM registry.hub.docker.com/library/nginx:1.13.0-alpine

RUN apk add --update nodejs && npm install npm@latest -g
WORKDIR /root/single-ui-fan

COPY maven/. .

RUN npm set progress=false && npm config set depth 0 && \
    npm install && npm run build:prod && \
    cp -R ./dist/* /usr/share/nginx/html && \
    cd / && rm -rf /root/single-ui-fan

WORKDIR /usr/share/nginx/html

COPY entrypoint*.sh /usr/local/bin/
RUN chmod a+x /usr/local/bin/*.sh

ENTRYPOINT [ "entrypoint.sh" ]
```

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Future

```
FROM nginx:1.13.0-alpine AS base

FROM base AS npm
RUN apk add --update nodejs && npm install npm@latest -g
WORKDIR /root/single-ui

FROM npm AS dependencies
COPY . .
RUN npm set progress=false && npm config set depth 0
RUN npm install
RUN npm run build:stage

FROM base AS release
COPY --from=dependencies /root/single-ui/dist /usr/share/nginx/html
```