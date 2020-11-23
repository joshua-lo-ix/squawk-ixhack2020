package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

var version string
var signingsecret string
var slacktoken string

func main() {

	flag.StringVar(&version, "version", "0", "the git commit of this build")
	flag.StringVar(&signingsecret, "signingsecret", "0", "Signing secret for slack")
	flag.StringVar(&slacktoken, "slacktoken", "0", "slack tokent for slack")

	flag.Parse()
	http.HandleFunc("/", handler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/ansibletest", ansibleTest)
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/modal", handleModal)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func ansibleTest(w http.ResponseWriter, r *http.Request) {
	initial_req_ack()
	w.WriteHeader(http.StatusOK)

	refresh_ansible()
	exec_ansible()
}
func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}
func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, modal)
}

func initial_req_ack() {
	msg := "Request received! :+1:"
	message(msg)
}

func refresh_ansible() {
	clean_ansible()
	get_ansible()
}

func clean_ansible() {
	cmd := exec.Command("rm", "-rf", "squawk-ixhack2020-ansible")
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func get_ansible() {
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

func message(msg string) {
	api := slack.New("xoxb-1318525037713-1546374229952-IVKijTnx5boTgt8PjjvTd7wA")
	_, _, err := api.PostMessage("C01FCNNDC4B",
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
}
