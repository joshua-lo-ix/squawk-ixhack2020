package main
import (
    "fmt"
    "log"
    "net/http"
    "os/exec"
)

func main() {
    cmd := exec.Command("ansible-playbook", "-i", "squawk-ixhack2020/ansible/hosts", "squawk-ixhack2020/ansible/squawk-playbook.yml")
    if err := cmd.Run(); err != nil {
        log.Fatal(err)
    }
}

func refresh_ansible() {
    cmd := exec.Command("git", "clone", "https://github.com/joshua-lo-ix/squawk-ixhack2020.git")
}
