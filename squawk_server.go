package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/slack-go/slack"
	"github.com/google/uuid"
)

var version string
var signingsecret string
var slacktoken string
<<<<<<< Updated upstream
var jobUuid uuid.UUID
=======
var database *sql.DB
>>>>>>> Stashed changes

func main() {

	flag.StringVar(&version, "version", "0", "the git commit of this build")
	flag.StringVar(&signingsecret, "signingsecret", "0", "Signing secret for slack")
	flag.StringVar(&slacktoken, "slacktoken", "0", "slack tokent for slack")

	flag.Parse()

	os.Remove("./localsqllite.db")
	database, err := sql.Open("sqlite3", "./localsqllite.db")
	if err != nil {
		log.Printf("%q: \n", err)
		os.Exit(1)
	}
	statement, err := database.Prepare("CREATE TABLE IF NOT EXISTS jobs (id INTEGER PRIMARY KEY, targetservers TEXT, ixconfs TEXT, results TEXT)")
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
		os.Exit(1)
	}

	_, err = statement.Exec()
	if err != nil {
		log.Printf("%q: %s\n", err, statement)
		os.Exit(1)
	}
	exec_ansible("foo", "bar")

	http.HandleFunc("/", handler)
	http.HandleFunc("/version", versionHandler)
	http.HandleFunc("/ansibletest", ansibleTest)
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/squawk", handleSlash)
	//http.HandleFunc("/slash", fastSlash)
	http.HandleFunc("/modal", ansibleTest)

	http.HandleFunc("/jobs", handleJobs)

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
	message(fmt.Sprintf("Job ID: `%v`  |  :bird: Request received!", jobUuid))
}

func refresh_ansible() {
	clean_ansible()
	get_ansible()
}

func clean_ansible() {
	out, err := exec.Command("rm", "-rf", "squawk-ixhack2020-ansible").Output()
	if err != nil {
		message(fmt.Sprintf("Job ID: `%v`  |  Failed to clean Ansible material ```%s```", jobUuid, string(out)))
		log.Fatal(err)
	}
}

func get_ansible() {
	out, err := exec.Command("git", "clone", "https://github.com/joshua-lo-ix/squawk-ixhack2020-ansible.git").Output()
	if err != nil {
		message(fmt.Sprintf("Job ID: `%v`  |  Failed to clone Ansible material ```%s```", jobUuid, string(out)))
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

	message(fmt.Sprintf("Job ID: `%v`  |  :airplane_departure: Running command: ```ansible-playbook %s```", jobUuid, strings.Join(args, " ")))
	out, err := exec.Command("ansible-playbook", args...).Output()

	if err != nil {
		message(fmt.Sprintf("Job ID: `%v`  |  :octagonal_sign: Error running command: ```%v```", jobUuid, err))
		message(fmt.Sprintf(":mag: Playbook output: ```%s```", string(out)))
		return
	}

	database, _ := sql.Open("sqlite3", "./localsqllite.db")
	//spew.Dump(database)

	tx, _ := database.Begin()
	statement, _ := tx.Prepare("INSERT INTO jobs (targetservers, ixconfs, results) VALUES (?, ?, ?)")
	defer statement.Close()
	statement.Exec(targetServers, ixConfs, string(out))
	tx.Commit()

	outString := string(out)
	outSplit := regexp.MustCompile("PLAY RECAP \\*+").Split(outString, -1)
	if len(outSplit) < 2 {
		message(fmt.Sprintf("Job ID: `%v`  |  :rotating_light: Error parsing PLAY RECAP: ``` %s```", jobUuid, outString))
		return
	}

	message(fmt.Sprintf("Job ID: `%v`  |  :airplane_arriving: Successfully completed! ```%s```", jobUuid, outSplit[1]))
}

func message(msg string) {

	api := slack.New(slacktoken)
	_, _, _ = api.PostMessage("C01FCNNDC4B",
		slack.MsgOptionText(msg, false),
		slack.MsgOptionAttachments())
}
