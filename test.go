package main

import (
	"os/exec"
	"fmt"
	"os"
	// "strings"
)

func main(){
	// file, _ := os.OpenFile("./docker.ps1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// write1 := fmt.Sprintf("Invoke-Command {docker exec -it biever /bin/bash}")
	// _,_ = file.WriteString(write1)
	_ = os.Chdir("/mnt/c/Windows/System32/WindowsPowerShell/v1.0")
	dockerexec := exec.Command(`./powershell.exe`, "-File" , `/mnt/c/Users/Mai Dinh Phuc/Desktop/MDP/MDP-apps/biever/docker.ps1`)
	err := dockerexec.Start()
	fmt.Println(err)


}