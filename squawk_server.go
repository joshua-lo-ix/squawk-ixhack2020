package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"

	"github.com/slack-go/slack"
	"github.com/google/uuid"
)

var version string
var signingsecret string
var slacktoken string
var jobUuid uuid.UUID

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
		jobUuid = uuid.New()
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
	message(fmt.Sprintf("Request received! :bird: ::    Job ID: `%v`", jobUuid))
}

func refresh_ansible() {
	clean_ansible()
	get_ansible()
}

func clean_ansible() {
	out, err := exec.Command("rm", "-rf", "squawk-ixhack2020-ansible").Output()
	if err != nil {
		message(fmt.Sprintf("Failed to clean Ansible material ::    Job ID: `%v` ```%s```", jobUuid, string(out)))
		log.Fatal(err)
	}
}

func get_ansible() {
	out, err := exec.Command("git", "clone", "https://github.com/joshua-lo-ix/squawk-ixhack2020-ansible.git").Output()
	if err != nil {
		message(fmt.Sprintf("Failed to clone Ansible material ::    Job ID: `%v` ```%s```", jobUuid, string(out)))
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

	args = append(args, "squawk-ixhack2020-ansible/squawk-playbook.yml")

	message(fmt.Sprintf(":airplane_departure: Running command: ```ansible-playbook %s``` ::    Job ID: `%v`", strings.Join(args, " "), jobUuid))
	out, err := exec.Command("ansible-playbook", args...).Output()

	if err != nil {
		message(fmt.Sprintf(":octagonal_sign: Error running command: ```%v``` ::    Job ID: `%v`", err, jobUuid))
		message(fmt.Sprintf(":mag: Playbook output: ```%s```", string(out)))
		return
	}

	outString := string(out)
	outSplit := regexp.MustCompile("PLAY RECAP \\*+").Split(outString, -1)
	if len(outSplit) < 2 {
		message(fmt.Sprintf(":rotating_light: Error parsing PLAY RECAP: ``` %s``` ::    Job ID: `%v`", outString, jobUuid))
		return
	}

	message(fmt.Sprintf(":airplane_arriving: Successfully completed! ::    Job ID: `%v` ```%s```", jobUuid, outSplit[1]))
}

func message(msg string) {
	api := slack.New(slacktoken)
	_, _, _ = api.PostMessage("C01FCNNDC4B",
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
}
