package main

import (
	"fmt"
	// "context"
	"os/exec"
	// "strings"
	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
)

func main(){
	// credentials, _ := os.ReadFile("./credential.log")
	// credential := strings.Split(string(credentials), "\n")
	// fmt.Println(credential[0])
	cmd := exec.Command("docker", "network", "ls")
	_, err := cmd.CombinedOutput()

	fmt.Println(err)
}