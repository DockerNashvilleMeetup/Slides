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
        <li>Organizer for Nashville Docker and Nashville Go Meetups</li>
        <li>I enjoy cycling, competitive shooting, and board games</li>
    </ul>
</div>

<!-- slide 2 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## monitoring docker

* installation
* prometheus
* collectors
* pros/cons
* vendors

<!-- slide 3 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## installation 
Feel free to follow along if you have a linux vm or linux host.<br /> (this might work on osx, too)

1. clone: git@github.com:DockerNashvilleMeetup/GrafanaDemo.git
2. initialize swarm: `docker swarm init`
3. start up: `./run.sh`

<!-- slide 4 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain


## prometheus

open source tsdb (time-series database)

* scrapes telemetry data from ancillary services
* provides a functional expression language that lets the user select and aggregate time series data in real time
* alerting rules define conditions based on expressions and can send notifications about alerts to an external service (pagerduty, slack, etc)

<!-- slide 5 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## prometheus (collectors)

config/prometheus-config.yml
```
  - job_name: "gateway"
    scrape_interval: 5s
    dns_sd_configs:
    - names: ['tasks.gateway']
      port: 8080
      type: A
      refresh_interval: 5s
```

<!-- slide 6 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## prometheus (alerting)

config/prometheus-alert-rules.yml
```
  - alert: node_disk_fill_rate_6h
    expr: predict_linear(node_filesystem_free{mountpoint="/"}[1h], 6*3600) 
          * on(instance) group_left(node_name) node_meta < 0
    for: 1h
    labels:
      severity: critical
    annotations:
      summary: "Disk fill alert for Swarm node '{{ $labels.node_name }}'"
      description: "Swarm node {{ $labels.node_name }} disk is going to fill up in 6h."
```

<!-- slide 7 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## collectors


```
$ docker run -t --rm \
  --network monitoring byrnedo/alpine-curl \ 
  http://tasks.gateway:8080/metrics

process_resident_memory_bytes 8.671232e+06
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.51897387334e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 1.499136e+07
```

<!-- slide 8 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

## pros/cons

<div style="width: 50%; float: left;">
    <h3>Pros</h3>
    <ul>
        <li>Free. Nothing quite like free</li>
        <li>Easy to deploy</li>
        <li>Large community</li>
        <li>Huge dashboard library</li>
    </ul>
</div>
<div style="float:left; width: 50%;">
    <h3>Cons</h3>
    <ul>
        <li>Hope you have plenty of CPUs</li>
        <li>Time consuming setting up custom dashboards</li>
        <li>Got trust issues?</li>
    </ul>
</div>


<!-- slide 9 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

## vendors

* Instana
* Sysdig
* Datadog

<!-- slide 9 -->

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
