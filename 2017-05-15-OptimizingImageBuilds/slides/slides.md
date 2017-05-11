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
* Dockerfiles: 
* Docker Security Scan
* The Future

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
