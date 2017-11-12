func cg() {
	cgroups := "/sys/fs/cgroup"
	pids := filepath.Join(cgroups, "pids")
	os.Mkdir(filepath.Join(pids, "kc"), 0755)
	must(ioutil.WriteFile(filepath.Join(pids, "kc/pids.max"), []byte("20"), 0700))
	//remove the cgroup after container exits
	must(ioutil.WriteFile(filepath.Join(pids, "kc/notify_on_release"), []byte("1"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids, "kc/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())), 0700))
}
