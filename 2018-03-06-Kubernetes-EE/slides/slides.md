class: top
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain

<div style="width: 25%; margin-left: -35px; margin-top:40px; float: left;">
    <img style="width: 100%; outline: 2px solid black" src="./slides/images/portrait-ws.png">
    <div style="color: #FFFFFF; background-color: #475258; padding: 10px; border-radius: 5px;
        border: 2px solid; border-color: #000000; margin-top: 5px; font-size: 18px;">
        Kevin Crawley<br />
        Engineering Manager<br />
        Franklin American<br />
    </div>
</div>
<div style="position: relative; float:right; width: 70%; align: right;">
    <h2>about me</h2>
    <ul>
        <li>Middle Tennessee Native</li>
        <li>Docker user since 2014; putting Docker systems in production since 2015</li>
        <li>Organizer for Nashville Docker, Go, and Serverless Meetups</li>
        <li>I enjoy cycling, competitive shooting, and board games</li>
    </ul>
</div>

<!-- slide 2 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## kubernetes and docker!?

* overview
* installation
* security
* what works?
* why pay for this?
* kube for mac/win

<!-- slide 3 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## overview 

<div class="center">
<img style="text-align: center" width="720px" src="slides/images/ee-kube.png">
</div>

<!-- slide 4 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## installation

* install docker-ee :: `sudo apt-get install docker-ee`
* install the ucp bundle :: 

```
docker container run --rm -it --name ucp \
  -v /var/run/docker.sock:/var/run/docker.sock \
  docker/ucp:3.0.0-beta2 install \
  --host-address <node-ip-address> \
  --interactive
```
* join worker nodes to the manager - you now have swarm/kube in box

<!-- slide 5 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## security

* docker enterprise leverages their collections for kubernetes RBAC
* this probably means nothing to you
* caveat: stuff is broken because of this; helm, tiller, etc
* they're supposed to be fixing this
* getting traefik to work was weird -- but i got it to work

<!-- slide 6 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## what works

* i was able to deploy a fairly "normal" swarm template as a kube service
* secrets are weird (i haven't touched kubernetes in like 2 years)

```
kubectl create secret generic atsea-payment-token --from-file=file=atsea-payment-token
```
* i was also able to deploy some other demo kubernetes services without issue

<!-- slide 7 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## what doesn't work

* RBAC is still in-progress, stuff like Gitlab Kubernetes, Helm, etc. don't work because of the default permissions. (i'm not entirely sure on how RBAC works or could properly articulate WHY)
* Interlock 2.0 (Swarm-native load balancer) doesn't work for Kubernetes
* Deploying Swarm Services with volumes as Kubernetes services doesn't work. This isn't a huge suprise
* Swarms vxlan/overlay network doesn't talk to Calico (the default network which is installed with Kube)

<!-- slide 8 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## why pay for this?

* In most cases Docker just works out-of-the-box
* Their security model is pretty slick, once you get used to it
* Docker Trust Registry / Notary is a fantastic product
* You need a vendor you can trust / yell at

<!-- slide 9 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

<!-- slide 10 -->

---
class: top
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain


<div style="width: 25%; margin-left: -35px; margin-top:40px; float: left;">
    <img style="width: 100%; outline: 2px solid black" src="./slides/images/portrait-ws.png">
    <div style="color: #FFFFFF; background-color: #475258; padding: 10px; border-radius: 5px;
        border: 2px solid; border-color: #000000; margin-top: 5px; font-size: 18px;">
        Kevin Crawley<br />
        Engineering Manager<br />
        Franklin American<br />
    </div>
</div>
<div style="position: relative; float:right; width: 70%; align: right;">
    <h1>Thank You!</h1>
    Questions? Comments?
</div>
