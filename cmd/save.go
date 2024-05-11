/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"strings"
	"github.com/spf13/cobra"
	// "github.com/fatih/color"
	"os/exec"
	// "os"
	"time"
	"bytes"
)

// saveCmd represents the save command
var saveCmd = &cobra.Command{
	Use:   "save",
	Short: "Save state of docker development environment",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		check_docker := exec.Command("docker", "network", "ls")
		var stderr_check_docker bytes.Buffer
		check_docker.Stderr = &stderr_check_docker

		err := check_docker.Run()
		if err != nil{
			initLog("save.go " + fmt.Sprintf("%v", err))
			red.Print("==> [Error] ")
			white.Print("Docker not running\n")
		}else{
			if name == ""{
				red.Print("==> [Error] ")
				white.Print("Missing or invalid flag\n")
			}else{
				Save(name)
			}
		}
		
	},
}

func Save(image string){
		backup_create := exec.Command("mkdir", "backup")
		var stderr_backup_create bytes.Buffer
		backup_create.Stderr = &stderr_backup_create

		_ = backup_create.Run()

		initLog("save.go " + string(stderr_backup_create.Bytes()))

		dt := time.Now()
		ct := dt.Format("01-02-2006 15:04:05 MST")
		ct = strings.ReplaceAll(ct, " ", "")
		ct = strings.ReplaceAll(ct, "-", "")
		ct = strings.ReplaceAll(ct, ":", "")

		save_name := fmt.Sprintf("./backup/%v-%v.tar", image, ct)
		dockersave := exec.Command("docker", "image", "save", "-o", save_name, image)
		var stderr_dockersave bytes.Buffer
		dockersave.Stderr = &stderr_dockersave

		blue.Print("==> [In Progress] ")
		white.Print("Saving...\n")
		errsave := dockersave.Run()

		imagename := strings.Split(save_name, "-")
		container_id := GetContainerID(imagename[0])

		if errsave != nil{
			initLog("save.go " + string(stderr_dockersave.Bytes()))
			red.Print("==> [Error] ")
			white.Print(fmt.Sprintf("Cannot find docker image named or has id %v\n", image))
		}else{
			if container_id == "Not Found"{
				red.Print("==> [Error] ")
				white.Print(fmt.Sprintf("Container not running"))
			}else{
				dockercommit := exec.Command("docker", "commit", container_id, imagename[0])
				var stderr_dockercommit bytes.Buffer
				dockercommit.Stderr = &stderr_dockercommit

				_ = dockercommit.Run()
				_ = dockersave.Run()

				initLog("save.go " + string(stderr_dockercommit.Bytes()))
				initLog("save.go " + string(stderr_dockersave.Bytes()))
			}
			green.Print("==> [Success] ")
			white.Print(fmt.Sprintf("Backup saved at ./backup/\n"))
		}

}

func init() {
	dockerCmd.AddCommand(saveCmd)

	saveCmd.PersistentFlags().String("name", "", "Name of docker development environment image")

}
