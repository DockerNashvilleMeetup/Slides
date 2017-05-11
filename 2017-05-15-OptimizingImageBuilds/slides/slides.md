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

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Past

```
FROM    nginx:1.13.0-alpine

WORKDIR /root/single-ui-fan
RUN     apk add --update nodejs && npm install npm@latest -g

COPY    . .

RUN     npm set progress=false && npm config set depth 0 && \
        npm install && npm run build:prod && \
        cp -R ./dist/* /usr/share/nginx/html && \
        cd / && rm -rf /root/single-ui-fan

WORKDIR /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## Dockerfile: The Future

```
ARG     BASE_VERSION=1.13.0-alpine
FROM    nginx:${BASE_VERSION} AS base

FROM    base AS npm
WORKDIR /root/single-ui
RUN     apk add --update nodejs && npm install npm@latest -g

FROM    npm AS dependencies
COPY    . .
RUN     npm set progress=false && npm config set depth 0
RUN     npm install
RUN     npm run build:prod

FROM    base AS release
COPY    --from=dependencies /root/single-ui/dist /usr/share/nginx/html
COPY    entrypoint.sh /usr/local/bin/

ENTRYPOINT [ "entrypoint.sh" ]
```

---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

# Links and More

* https://docs.docker.com/engine/userguide/storagedriver/imagesandcontainers/