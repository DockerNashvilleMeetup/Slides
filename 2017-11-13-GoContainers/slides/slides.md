class: top
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain

<div style="width: 25%; margin-left: -35px; margin-top:40px; float: left;">
    <img style="width: 100%; outline: 2px solid black" src="./slides/images/portrait-ws.png">
    <div style="color: #FFFFFF; background-color: #475258; padding: 10px; border-radius: 5px;
        border: 2px solid; border-color: #000000; margin-top: 5px; font-size: 18px;">
        Kevin Crawley<br />
        Lead Engineer @ FAMC<br />
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

# Topics

* build a container in go
* namespaces
* cgroups
* security
* exploits

<!-- slide 3 -->
---
class: top
background-image: url('./slides/images/slide-bg-3.png')
background-size: contain

# Namespaces 
restrict what resources a process can see on a host

> unix timesharing system<br/>
> process ids<br/>
> mounts<br/>
> network<br/>
> user ids<br/>
> interprocess comms<br/>

<!-- slide 4 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain


## cgroups - restrict what resources a process can use on the host

- file system interface<br>
- processes are inherited from parent<br>
- can be reassigned to different cgroups

> memory<br/>
> cpu / cores<br/>
> devices<br/>
> io<br/>
> processes<br/>

<!-- slide 5 -->
---
class: top
background-image: url('./slides/images/slide-bg-2.png')
background-size: contain

<h1 style="text-align: center; margin-top: 150px;">` :() { :|: & }; : `</h2>


<!-- slide 6 -->
---
class: top
background-image: url('./slides/images/slide-bg-1.png')
background-size: contain


<div style="width: 25%; margin-left: -35px; margin-top:40px; float: left;">
    <img style="width: 100%; outline: 2px solid black" src="./slides/images/portrait-ws.png">
    <div style="color: #FFFFFF; background-color: #475258; padding: 10px; border-radius: 5px;
        border: 2px solid; border-color: #000000; margin-top: 5px; font-size: 18px;">
        Kevin Crawley<br />
        Lead Engineer @ FAMC<br />
        @notsureifkevin<br />
    </div>
</div>
<div style="position: relative; float:right; width: 70%; align: right;">
    <h1>Thank You!</h1>
    Questions? Comments?
</div>
