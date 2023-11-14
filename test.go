package main

import (
	"fmt"
	// "context"
	"os/exec"
	"os"
	"strings"
	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
)

func main(){
	// credentials, _ := os.ReadFile("./credential.log")
	// credential := strings.Split(string(credentials), "\n")
	// fmt.Println(credential[0])
	// cmd := exec.Command("docker", "network", "ls")
	// _, err := cmd.CombinedOutput()

	// fmt.Println(err)
	dir_name := GetDirectoryName()

	container_id := GetContainerID(dir_name)
	fmt.Println(container_id)
}

func GetContainerID(dir_name string) string {
	cmd := exec.Command("docker", "container", "ls", "--all", "--quiet", "--filter", fmt.Sprintf("name=%v", dir_name))
	output, _ := cmd.CombinedOutput()
	if string(output) != ""{
		container_id := strings.TrimSpace(string(output)[0:len(string(output))-1])
		return container_id
	}
	return "Not Found"
}

func GetDirectoryName() string{
	dir, _ := os.Getwd()
	dir_name_ := strings.Split(dir, "/")
	dir_name := dir_name_[len(dir_name_)-1]
	return dir_name
}