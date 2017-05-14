class: top
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain

<div style="width: 25%; margin-left: -35px; margin-top:40px; float: left;">
    <img style="width: 100%; outline: 2px solid black" src="./slides/images/portrait-ws.png">
    <div style="color: #FFFFFF; background-color: #475258; padding: 10px; border-radius: 5px;
        border: 2px solid; border-color: #000000; margin-top: 5px;">
        Kevin Crawley<br />
        Sr. Engineer @ FAMC<br />
        @notsureifkevin<br />
    </div>
</div>
<div style="position: relative; float:right; width: 70%; align: right;">
    <h1>About Me</h1>
    <ul>
        <li>Middle Tennessee Native</li>
        <li>Docker user since 2014; putting Docker systems in production since 2015</li>
        <li>Organizer for Nashville Docker and Nashville Go Meetups</li>
        <li>I enjoy cycling, competitive shooting, and board games</li>
    </ul>
</div>

<!-- slide 2 -->
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

## Pro-tip

* These slides can be viewed online at http://nashdocker.io/live

<!-- slide 3 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## What are Layers?

<table>
<tbody>
    <tr>
    <td style="vertical-align: top;">
        <ul>
            <li>Layers are composed and then digested</li>
            <li>Each digested layer is given a signature</li>
            <li>Signatures govern integrity and security</li>
            <li>More layers, more problems</li>
        </ul>
    </td>
    <td><img style="margin-top: -100px;" src="./slides/images/docker-image-comp.png"></td>
    </tr>
</tbody>
</table>

???
* images reference one or more layers which eventually contribute to a containers filesystem
* building images locally will produce intermediate "cache" layers which will appear in a `docker history`
* these cache layers are used when doing subsequent builds but will not be sent to a registry via push; so if you run a build on one machine and then pull on another machine and run the build there - it will not cache any of those layers
* discuss the implications of having more stuff in your builds (security, size, complexity)
  * security: more dependencies/software in your image == larger attack surface, we'll discuss more later
  * size: larger images mean larger disks, more bandwidth, more build time
  * complexity: this can be attacked boths ways, we'll discuss more later

<!-- slide 4 -->
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


???
* common base images are useful if you're extending base images with software which might be common across your services. (new relic, jolokia, supervisord, tini)
* 2+3 are your best defense. if you need to install a bunch of dependencies to compile/build, but don't need it to run your application try and clean up before the end of the RUN

<!-- slide 5 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Past

```Dockerfile
FROM    ubuntu:latest

RUN     apt-get update -y && apt-get install -y curl gnupg nginx
RUN     curl -sL https://deb.nodesource.com/setup_6.x | bash - 
RUN     apt-get install nodejs
RUN     npm install npm@latest -g

WORKDIR /root/single-ui-fan
COPY    single-ui-fan/. .
RUN     npm install
RUN     npm run build:prod
RUN     cp -R ./dist/* /usr/share/nginx/html

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

???
* explain that entrypoint.sh does some configuration (environment) and then starts nginx (webserver)
* this is a dockerfile. this is something i might have made in 2014/15. it's not bad, but it's not good either.
* ask the audience if they see any thing we can do better
  * highlight that we'll explore this file later

<!-- slide 6 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Result: 555 MB

```
image            tag                 size
----------------------------------------------------------------------
ubuntu           latest ?            117 MB
    \_ ubuntu, nginx

single-ui-fan    0.1.0-ubuntu        555 MB
    \_ ubuntu, nginx, nodejs+npm, app
```

???
* base image alone is 117mb, and we're installing nginx, nodejs, npm, node_modules, and our app
* the result: 555mb. doesn't seem terrible. i've seen worse.
* if you're pushing/pulling/deploying several times a day and aren't cleaning up you'll eventually run out of space.

<!-- slide 7 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Have you heard of Alpine?

<div class="center">
<img style="
    background-color: white; 
    max-width: 600px; padding: 10px;
    box-shadow: 0 0 10px 0 black" src="slides/images/alpine.png">
</div>

```
nginx            1.13.0-alpine      15.5 MB
    \_ alpine, nginx
```
???
* alpine linux is a minimal linux distribution
* minimal: includes slimmed down binary packages
* built on musc libc (rather than glibc)
  * we won't go into details of what that means, but it's significant enough to mention

<!-- slide 8 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## When should I use a full base OS?

<div style="float: left">
<ul>
<li>Compliance</li>
<li>Security</li>
<li>Ease of Use</li>
</div>
<div style="float: right">
<img style="" src="./slides/images/os-collage.png">
</div>
???
* larger organizations have requirements that small startups may not
* ask the audience what they think might be some of those requirements?
  * monitoring/logging packages
  * security scanning/package repos
  * these may be configured/built for a specific distro (probably not alpine)
* ask what could be some "ease of use" considerations?
  * easier to dev/debug containers which have common networking tools
  * familiar tooling yum/apt
* point out that it's possible to dev first using a common distro then minimizing with alpine

<!-- slide 9 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: Today

```Dockerfile
FROM    nginx:1.13.0-alpine

RUN     apk add --update nodejs && npm install npm@4.5.0 -g && \
        rm -rf /var/cache/apk/*

WORKDIR /root/single-ui-fan
COPY    single-ui-fan/package.json .
RUN     npm set progress=false && npm config set depth 0 && \
        npm install --production

COPY    single-ui-fan/. .
RUN     npm run build:prod && cp -R dist/. /usr/share/nginx/html && \
        cd / && rm -rf /root/single-ui-fan

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```
???
* highlight the use of:
  * alpine
  * specific base image
    * note that one could decide to use nodejs base image, and that's valid. however, you'll soon see why I did not.
  * specific versions (base and npm)
  * cleaning up in each RUN (both SRC and apk cache)

<!-- slide 10 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Face Off

<div class="center">
    <img style="max-width: 700px;" src="./slides/images/faceoff.jpg">
</div>
???
LOL NIC CAGE

<!-- slide 11 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Result: 291 MB

```
single-ui-fan   0.1.0-alpine        291 MB
    \_ alpine, nginx, nodejs+npm, app
```
<span class="center">vs.</span>
```
single-ui-fan   0.1.0-debian        555 MB
    \_ debian, nginx, nodejs+npm, app
```
<span class="center">can we do better than **48%**?</span>

???
* point out that we're shipping our application with a bunch of stuff we don't need
  * ask the audience what are some of the things we don't need.
    * nodejs/npm
    * node_modules

<!-- slide 12 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Multi Stage Builds

<div class="center">
<img src="./slides/images/science.gif">
</div>

<!-- slide 13 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Future

```Dockerfile
ARG     BASE_VERSION=1.13.0-alpine
FROM    nginx:${BASE_VERSION} AS base

FROM    base AS build
WORKDIR /root/single-ui-fan
RUN     apk add --update nodejs && npm install npm@4.5.0 -g
COPY    single-ui-fan/package.json .
RUN     npm set progress=false && npm config set depth 0 && \
        npm install
COPY    single-ui-fan/. .
RUN     npm run build:prod

FROM    base AS release
COPY    --from=build /root/single-ui-fan/dist /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

???
* yes, we can do better.
* introducing multi-stage builds
* ask the audience what are some of the differences
    * multiple FROMs doing in my Dockerfile?
    * regression? RUN statement is no longer compact
    * copy --from? wat?

<!-- slide 14 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Result: 20.4 MB

```
single-ui-fan   0.1.0-alpine        20.4 MB
    \_ alpine, nginx, app
```
<span class="center">vs.</span>
```
single-ui-fan   0.1.0-alpine        291 MB
    \_ alpine, nginx, nodejs+npm, app

single-ui-fan   0.1.0-debian        555 MB
    \_ debian, nginx, nodejs+npm, app
```
<span class="center">How does **96%** sound to you?</span>

???
* our final image only contains nginx and our app "distribution"
* earlier i mentioned that you could build based on nodejs, this is the reason why i did not.

<!-- slide 15 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Yeah...
<div class="center">
    <img width="720px" src="./slides/images/flying-laptop.gif">
</div>

<!-- slide 16 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

class: center
# Caching
<img style="max-height: 400px;" src="./slides/images/chaching.png"><br />
<div style="
font-size: 40pt;
font-weight: bold;
margin-top: -270px;
text-shadow: 0px 0px 6px #000000;
">
$<span id="exchrate"></span> USD
</div>

???
* cache is determined by the contents of each Dockerfile line and the layer which was digested
* cache means common steps like "apt-get" or "make" can be skipped on subsequent builds
* cache is neither push nor pulled 

<!-- slide 17 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: Let's start **FROM** the beginning
```Dockerfile
*FROM    ubuntu:latest

RUN     apt-get update -y && apt-get install -y curl gnupg nginx
RUN     curl -sL https://deb.nodesource.com/setup_6.x | bash - 
RUN     apt-get install nodejs
*RUN     npm install npm@latest -g

WORKDIR /root/single-ui-fan
COPY    single-ui-fan/. .
RUN     npm install
RUN     npm run build:prod
RUN     cp -R ./dist/* /usr/share/nginx/html

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

???
* explain that we're going to transform this image into the present "optimal" version
* using latest is bad
  * when debian dropped npm in their node package a Dockerfile with latest and `nodejs` via apt would have exploded
* consider a smaller base image; also, consider using official base images for things like php, mysql, postgres or nginx. they often have container specific tooling and optimizations.

<!-- slide 18 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: **RUN** for your lives

```Dockerfile
FROM    nginx:1.13.0-alpine

*RUN     apk add --update nodejs
*RUN     npm install npm@4.5.0 -g

WORKDIR /root/single-ui-fan
COPY    single-ui-fan/. .
*RUN     npm install
*RUN     npm run build:prod
*RUN     cp -R ./dist/* /usr/share/nginx/html

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

???
* we've combined our RUN statements to execute in a single layer
* before ending the RUN we try and clean up (if possible)
  * we can't delete/remove NPM/NODEJS because we need it to build
  * removing it in the subsequent run is pointless because it's already committed in a previous layer
  * we're removing the apk cache (recent versions this is unneccessary with apk add --no-cache)
  * we're removing the source code (and node_modules) directory

<!-- slide 19 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: Got **Cache**?

```Dockerfile
FROM    nginx:1.13.0-alpine

RUN     apk add --update nodejs && npm install npm@4.5.0 -g && \
        rm -rf /var/cache/apk/*

WORKDIR /root/single-ui-fan
*COPY    single-ui-fan/. .
*RUN     npm set progress=false && npm config set depth 0 && \
*       npm install --production && \
*       npm run build:prod && cp -R ./dist/* /usr/share/nginx/html && \
*       cd / && rm -rf /root/single-ui-fan

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```
???
* explain that the contents of your COPY will determine cache hits or misses
* data which changes frequently will invalidate the cache from that point forward
* people can/will do a lot of things to speed up their builds post-cache miss
  * cached repositories
  * copy local file cache
  * copy entire vendor/dependency folders (is this bad?)

<!-- slide 20 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: Are We There Yet?

```Dockerfile
FROM    nginx:1.13.0-alpine

RUN     apk add --update nodejs && npm install npm@4.5.0 -g && \
        rm -rf /var/cache/apk/*

WORKDIR /root/single-ui-fan
COPY    single-ui-fan/package.json .
RUN     npm set progress=false && npm config set depth 0 && \
        npm install --production

COPY    single-ui-fan/. .
RUN     npm run build:prod && cp -R dist/. /usr/share/nginx/html && \
        cd / && rm -rf /root/single-ui-fan

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```
???
* this is our "present image"; lets recap
  * we've picked a smaller and more specific base image
  * we've chained our run statements and tidied up after them
  * we've pinned our image and dependencies to a specific version
  * we've applied language specific optimizations to our build

<!-- slide 21 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: Back **FROM** The Future

```Dockerfile
ARG     BASE_VERSION=1.13.0-alpine
*FROM    nginx:${BASE_VERSION} AS base

*FROM    base AS build
WORKDIR /root/single-ui-fan
RUN     apk add --update nodejs && npm install npm@4.5.0 -g
COPY    single-ui-fan/package.json .
RUN     npm set progress=false && npm config set depth 0 && \
        npm install
COPY    single-ui-fan/. .
RUN     npm run build:prod

*FROM    base AS release
COPY    --from=build /root/single-ui-fan/dist /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```
???
* you can now use build args in FROM
* each FROM is an image independent of each other
* they can be aliased and referenced by another
* ask the audience what this means?
  * the last FROM is what is digested and pushed
  * if you had dockerfiles for specific stages of builds (im sorry), you no longer have to
  * why would I no longer have to chain RUNs in intermediate steps

<!-- slide 22 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: COPY what?

```Dockerfile
ARG     BASE_VERSION=1.13.0-alpine
FROM    nginx:${BASE_VERSION} AS base

FROM    base AS build
WORKDIR /root/single-ui-fan
RUN     apk add --update nodejs && npm install npm@4.5.0 -g
COPY    single-ui-fan/package.json .
RUN     npm set progress=false && npm config set depth 0 && \
        npm install
COPY    single-ui-fan/. .
RUN     npm run build:prod

FROM    base AS release
*COPY    --from=build /root/single-ui-fan/dist /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```
???
* copy directly from aliased images
  * this references a location on an images filesystem
  * think of it like artifacts. but the entire "image" contents is available

<!-- slide 23 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Application Specific Builds

* Consider what you're building
* Compile on your build server and then copy _into_ the container
* Try new things, compare, weigh your options


<!-- slide 24 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain
# Go Build

### Build the Binary

```shell
$  go build -o yourapp .
$  docker build -t yourapp .
```

### Dockerfile:
```Dockerfile
FROM alpine:3.5
COPY yourapp /usr/local/bin
ENTRYPOINT ["yourapp"]
```

<!-- slide 25 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain
# Java Build

### Build your JAR/WAR

```shell
$  mvn clean package
$  docker build -t yourapp .
```

### Dockerfile

```Dockerfile
FROM openjdk:8u121-jre-alpine
COPY target/yourapp*.jar ./app.jar
ENTRYPOINT [ "java", "-Djava.security.egd=file:/dev/urandom", "-jar", "app.jar" ]
```

<!-- slide 26 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# NodeJS
* Copy your package.json and run `npm install` before building your app. The tradeoff is worth it.
* Utilize .dockerignore and ignore unimportant things (like `npm-debug.log` and `.git`)

Example
```Dockerfile
COPY package.json .
RUN npm install --production
COPY . .
RUN npm run build
```

<!-- slide 27 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Security

* Docker Security Scan is the cats pajamas
* CVE scanning
* Matches binary signatures in each layer
  * OS packages
  * Component level (JAR, CPAN, PIP, etc)
* Available on private Docker Cloud ($) and Docker Trusted Registry EE ($$$$)

<!-- slide 28 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

class: center
<img style="max-width: 820px" src="./slides/images/dtr.png">

<!-- slide 29 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Image Signing

* Key signatures can be leveraged for authenticity and validity

<div class="center" style="margin-top: 10px">
<img src="./slides/images/trust.png">
</div>

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Recap

todo

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain




# Links and More

todo
* https://docs.docker.com/engine/userguide/storagedriver/imagesandcontainers/