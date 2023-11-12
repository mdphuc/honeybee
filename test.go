package main

import (
	"os/exec"
	"fmt"
	// "context"
	"os"
	// "strings"
	// "github.com/docker/docker/api/types"
	// "github.com/docker/docker/client"
)

func main(){
	// file, _ := os.OpenFile("./docker.ps1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// write1 := fmt.Sprintf("Invoke-Command {docker exec -it biever /bin/bash}")
	// _,_ = file.WriteString(write1)
	// _ = os.Chdir("/mnt/c/Windows/System32/WindowsPowerShell/v1.0")
	dockerexec := exec.Command("docker", "exec", "-it", "8cc3f81e97ba", "sh", "-c", "mkdir /biever")
	dockerexec.Stdin = os.Stdin
	dockerexec.Stdout = os.Stdout
	dockerexec.Stderr = os.Stderr

	err := dockerexec.Run()
	fmt.Println(err)
	// cli, err1 := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	// fmt.Println(err1)

	// // execID , err2 := cli.ContainerExecCreate(context.Background(), "882088d89f57", types.ExecConfig{User : "root", AttachStdin : true, AttachStderr: true, AttachStdout: true, Tty: true, Cmd : []string{"/bin/bash"}})

	// exec_attach, err := cli.ContainerAttach(context.TODO(), "882088d89f57", types.ContainerAttachOptions{Stream : true, Stdin: true})


	// // fmt.Println(execID)
	// // fmt.Println(err2)

	// fmt.Println(exec_attach)
	// fmt.Println(err)
}