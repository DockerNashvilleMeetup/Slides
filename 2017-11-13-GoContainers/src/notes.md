## Slide 1

build a container in go

* namespaces
  - linux kernel features which allows the user to control the visibility of the system around a process
  - introduced around 2007, embraced by google
  - lxc (linux containers) started a movement, docker automated it
* cgroups
  - more features which allow you to restrict the resources a process may consume
* security
* exploits

## Slide 2

# DO IT LIVE

initial copy of go program

* brief on `exec.Command()` and that it's creating a forked process
* `must()` is just a catchall error handler. try/catch for those java/php devs among us.

```
package main

import (
	"fmt"
	"os"
	"os/exec"
)

// go run main.go run <cmd> <args>
func main() {
}

func run() {
	cmd := exec.Command()
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}

func must(err error) {
	if err != nil {
		fmt.Printf("has error: %v \n", err)
	}
}
```

# CLI + fork/exec
```
func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		panic("nope")
	}
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(cmd.Run())
}
```
* add switch to `main()`
* add output to `run()` for debugging
* add arguments to `exec.Command()`
* run `go run main.go run echo hello` to demonstrate it works. we aren't in a container yet, we're just forking and executing a process



# CLONE_NEWUTS
```
	...
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}
```

* add `SysProcAttr` to cmd
* add `Cloneflags` -> `syscall.CLONE_NEWUTS`
  - this is the namespace for `hostname`, insert funny joke about naming things
* demonstrate setting the hostname within the container
  - `go run main.go run /bin/bash`
  - `hostname container` -> inside container
  - `hostname` -> outside container
* its possible to set the hostname during the application runtime
* `must(syscall.Sethostname([]byte("container")))`
* but this wont work because its executed prior to the fork/exec and we're just setting the hostname in the parent process


# Just keep forking!
```
func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
		...
}

func run() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	must(cmd.Run())
}

func child() {
	fmt.Printf("Running %v \n", os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	must(syscall.Sethostname([]byte("container")))
	must(cmd.Run())
}
```
* copy/paste our run function, rename it to child() and adjust arguments
  - `run()` should execute `/proc/self/exe` and append the original arguments to a `child` argument
  - `child()` may now set the hostname within the namespace and execute our command
* we will setup a namespace before forking a new process
* we will setup a process id namespace to isolate processes spawned from the parent process
  - demonstrate that `ps` still shows parent process
  - demonstrate that `ls /proc` still shows host processes since we're still on the host rootfs

# Chroot Jail
```
func child() {
	...
	must(syscall.Sethostname([]byte("container")))
	must(syscall.Chroot("/home/kcrawley/go/src/gocontainer/ubuntufs"))
	must(os.Chdir("/"))
}
```
* we've setup an empty ubuntu filesystem for our container, demonstrate `ls ubuntufs` and `ls /` to show the indicators
* we will tell the kernel to put our process into a new rootfs, through chroot
* when chroot'ing we are abandoning our current pointer in the filesystem, so we must chdir to the new rootfs
* attempt to show our process list `ps`

# Mount Proc
```
	must(syscall.Mount("proc", "proc", "proc", 0, ""))
	must(cmd.Run())
	must(syscall.Unmount("proc", 0))
```
* point out that we haven't setup a mount namespace, yet.
* now we can run `ps`, show off the slim proc folder too, `ls /proc`

# Mount Namespace
```
// add mount for /mnt/container_mnt
func child() {
	...
	must(syscall.Mount("cnt_mnt", "/mnt/container_mnt", "tmpfs", 0, ""))
	...
	must(syscall.Unmount("/mnt/container_mnt"))
}
```

* demonstrate that files created in the container under `/mnt/container_mnt` are visable on the host in `/home/kcrawley/go/src/gocontainer/ubuntufs/mnt/container_mnt`
* we can use a mount namespace to keep changes isolated to the container

```
// add syscall.CLONE_NEWNS
func run() {
	...
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
}

// enforce privacy for mounts in the jail
func child() {
	...
	must(syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""))
}
```

* demonstrate that files are no longer visible on the host

# Security

* run `sleep 1000` inside container
* identify `pidof sleep` on host
* `ls /proc/CPID` - wats this?
* `cat /proc/CPID/mounts`
* `mount | grep cnt_mnt`

## There is MOAR

* set a secret in the container, `export SECRET=horsecarbatterystaple`
* run `sleep 1000` in the container
* host: `pidof sleep`
* `cat /proc/CPID/environ | tr '\0' '\n' | grep SECRET`
* heard of 12 factor app?
* they say you should put secrets in the environment
* i just showed you why you shouldn't

# Slide 3

# CGROUPS

* `mount | grep cgroup`
* `cd /sys/fs/cgroup/memory && ls`
  - explain that you can assign processes to cgroups through this filesystem
* `cat cgroup.procs`
---
* add `cg()` func (from `cg.go`) to `main.go` and call `cg()` at the top of `child()`
* explain what `cg()` is doing
* run container
* on host:
  - `cd /sys/fs/cgroup/pids`
  - `ls` : woop, see `kc` there
  - `cd kc`
* on container:
  - `sleep 1000`
* on host:
  - `pidof sleep`
  - `cat cgroup.procs` and see that we have our pid
  - `cat pids.max` orly?

# Slide 4

# FORK BOMB

- define a function called colon
- inside that function, call colon, and pipe the results of colon to colon, and run it in the background
- then we're going to invoke it.
- profit???
- just kiddin, we've constrained a fork bomb