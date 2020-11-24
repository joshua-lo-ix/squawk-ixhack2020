package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"

	"github.com/slack-go/slack"
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
	http.HandleFunc("/squawk", handleSlash)
	//http.HandleFunc("/slash", fastSlash)
	http.HandleFunc("/modal", ansibleTest)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func ansibleTest(w http.ResponseWriter, r *http.Request) {

	targetServers, ixConfs := handleModal(w, r)

	go func() {
		initial_req_ack()
		refresh_ansible()
		exec_ansible(targetServers, ixConfs)
	}()

	w.WriteHeader(http.StatusOK)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, version)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, modal)
}

func fastSlash(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	fmt.Fprintf(w, serverModal)
}

func initial_req_ack() {
	message("Request received! :+1:")
}

func refresh_ansible() {
	clean_ansible()
	get_ansible()
}

func clean_ansible() {
	out, err := exec.Command("rm", "-rf", "squawk-ixhack2020-ansible").Output()
	if err != nil {
		message(fmt.Sprintf("```%s```", string(out)))
		log.Fatal(err)
	}
}

func get_ansible() {
	out, err := exec.Command("git", "clone", "https://github.com/joshua-lo-ix/squawk-ixhack2020-ansible.git").Output()
	if err != nil {
		message(fmt.Sprintf("```%s```", string(out)))
		log.Fatal(err)
	}
}

func exec_ansible(targetServers string, ixConfs string) {
	args := []string{"-i", "squawk-ixhack2020-ansible/hosts"}

	if targetServers != "all" {
		args = append(args, "--limit", targetServers)
	}

	if ixConfs != "all" {
		args = append(args, "--extra-vars", fmt.Sprintf(`ix_confs_selective_files=["%s"]`, ixConfs))
	}

	out, err := exec.Command("ansible-playbook", args...).Output()
	outString := string(out)
	outString = regexp.MustCompile("PLAY RECAP \\*+$").Split(outString, -1)[1]

	message(fmt.Sprintf("```%s```", string(outString)))
	if err != nil {
		log.Fatal(err)
	}
}

func message(msg string) {
	api := slack.New(slacktoken)
	_, _, _ = api.PostMessage("C01FCNNDC4B",
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
}
