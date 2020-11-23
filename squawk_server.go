package main

import (
	"flag"
	"fmt"
	"log"
    "net/http"
    "os/exec"
)

var version string

func main() {
	flag.StringVar(&version, "version", "0", "the git commit of this build")
	flag.Parse()
	http.HandleFunc("/", handler)
    http.HandleFunc("/version", versionHandler)
    http.HandleFunc("/ansibletest", ansibleTest)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func ansibleTest(w http.ResponseWriter, r *http.Request) {
    refresh_ansible()
    exec_ansible()
}
func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, modal)
}

func refresh_ansible() {
    cmd := exec.Command("rm", "-rf", "squawk-ixhack2020-ansible")
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
    cmd := exec.Command("git", "clone", "https://github.com/joshua-lo-ix/squawk-ixhack2020-ansible.git")
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
}

func exec_ansible() {
    cmd := exec.Command("ansible-playbook", "-i", "squawk-ixhack2020-ansible/hosts", "squawk-ixhack2020-ansible/squawk-playbook.yml")
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
}
