/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"os"
	"os/exec"
	"bytes"
)

// resetCmd represents the reset command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the environment",
	Run: func(cmd *cobra.Command, args []string) {
		docker := cmd.Flags().Lookup("docker").Changed
		remote_machine := cmd.Flags().Lookup("remote_machine").Changed
		
		check_docker := exec.Command("docker", "network", "ls")
		var stderr_check_docker bytes.Buffer
		check_docker.Stderr = &stderr_check_docker
		err := check_docker.Run()

		if err != nil{
			initLog("reset.go " + string(stderr_check_docker.Bytes()))
			red.Print("==> [Error] ")
			white.Print("Docker not running\n")
		}else{
			if docker == false && remote_machine == false{
				red.Print("==> [Error] ")
				white.Print("Missing flag\n")	
			}else if docker == true && remote_machine == true{
				red.Print("==> [Error] ")
				white.Print("Invalid combination of flag\n") 
			}else{
				dir_name := GetDirectoryName()
				container_id := GetContainerID(dir_name)
				if container_id == "Not Found"{
					red.Print("==> [Error] ")
					white.Print("The environment not found\n") 	
				}else{
					if docker == true{
						dockerstop := exec.Command("docker", "stop", container_id)
						var stderr_dockerstop bytes.Buffer
						dockerstop.Stderr = &stderr_dockerstop

						dockerprune := exec.Command("docker", "system", "prune")
						var stderr_dockerprune bytes.Buffer
						dockerprune.Stderr = &stderr_dockerprune

						dockerrmi := exec.Command("docker","rmi",dir_name)
						var stderr_dockerrmi bytes.Buffer
						dockerrmi.Stderr = &stderr_dockerrmi

						blue.Print("==> [In Progress] ")
						white.Print("Deleting the environment...\n")		
						_ = dockerstop.Run()
						_ = dockerrmi.Run()
						_ = dockerprune.Run()
						_ = dockerprune.Run()
						
						initLog("reset.go " + string(stderr_dockerstop.Bytes()))
						initLog("reset.go " + string(stderr_dockerrmi.Bytes()))
						initLog("reset.go " + string(stderr_dockerprune.Bytes()))
						initLog("reset.go " + string(stderr_dockerprune.Bytes()))

						green.Print("==> [Success] ")
						white.Print(fmt.Sprintf("Deletion complete\n"))
	
					}else{
						credentials, err := os.ReadFile("./credential.log. Please add <username>@<ip> to the first line of credential.log\n")
						if err == nil{
							credential := strings.Split(string(credentials), "\n")
							reset_command := fmt.Sprintf("ssh -i /root/.ssh/id_rsa %v 'rm -rf /%v'", credential[0], dir_name)
							dockerrm := exec.Command("docker", "exec", "-it", container_id, "sh", "-c", reset_command)
							var stderr_dockerrm bytes.Buffer
							dockerrm.Stderr = &stderr_dockerrm

							dockerstop := exec.Command("docker", "stop", container_id)
							var stderr_dockerstop bytes.Buffer
							dockerstop.Stderr = &stderr_dockerstop

							dockerprune := exec.Command("docker", "system", "prune")
							var stderr_dockerprune bytes.Buffer
							dockerprune.Stderr = &stderr_dockerprune

							dockerrmi := exec.Command("docker","rmi",dir_name)
							var stderr_dockerrmi bytes.Buffer
							dockerrmi.Stderr = &stderr_dockerrmi
								
							blue.Print("==> [In Progress] ")
							white.Print("Resetting the environment...\n")			
							_ = dockerrm.Run()
							_ = dockerstop.Run()
							_ = dockerrmi.Run()
							_ = dockerprune.Run()
							_ = dockerprune.Run()
							
							initLog("reset.go " + string(stderr_dockerrm.Bytes()))
							initLog("reset.go " + string(stderr_dockerstop.Bytes()))
							initLog("reset.go " + string(stderr_dockerrmi.Bytes()))
							initLog("reset.go " + string(stderr_dockerprune.Bytes()))
							initLog("reset.go " + string(stderr_dockerprune.Bytes()))
	
	
							green.Print("==> [Success] ")
							white.Print(fmt.Sprintf("Reset complete\n"))
	
		
						}else{
							red.Print("==> [Error] ")
							white.Print("credential.log not found. \n") 			
						}
					}
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(resetCmd)

	resetCmd.PersistentFlags().BoolP("docker","", true, "Username for remote machine")
	resetCmd.PersistentFlags().BoolP("remote_machine", "", true, "Build proxy server")

}

