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

	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"os/exec"
	"strings"
	// "slices"
	b64 "encoding/base64"
	"os"
)

var red = color.New(color.FgRed, color.Bold)
var white = color.New(color.FgWhite, color.Bold)
var green = color.New(color.FgGreen, color.Bold)
var blue = color.New(color.FgBlue, color.Bold)


// remoteMachineCmd represents the remoteMachine command
var remoteMachineCmd = &cobra.Command{
	Use:   "remoteMachine",
	Short: "Set up proxy server to use remote machine as development environment",
	Run: func(cmd *cobra.Command, args [string]){
		ip,_ := cmd.Flags().GetString("ip")
		user,_ := cmd.Flags().GetString("user")
		compose := cmd.Flags().Lookup("compose").Changed
		connect := cmd.Flags().Lookup("connect").Changed

		RemoteMachine(ip, user, compose, connect)
	}
}

func init() {
	runCmd.AddCommand(remoteMachineCmd)

	remoteMachineCmd.PersistentFlags().String("ip", "", "IP address of the remote environment")
	runCmd.PersistentFlags().String("user", "u", "", "Username for remote machine")
	runCmd.PersistentFlags().BoolP("compose", "", true, "Build proxy server")
	runCmd.PersistentFlags().BoolP("connect", "c" , "", true, "Set up connection between proxy server and remote machine")
	
}


func RemoteMachine(ip string, username string, compose bool, connect bool){
	if compose == true {
		if reset == false && connect == false{
			dir_name := GetDirectoryName()
			container_id := GetContainerID(dir_name)
			if container_id == "Not Found"{
				if ip != "" && username != ""{
					dir_mount := fmt.Sprintf(".:/%v", dir_name)
					workdir := fmt.Sprintf("/%v", dir_name)
					monitor_service_b64 := "W1VuaXRdCkRlc2NyaXB0aW9uPU1vbml0b3IgRmlsZQoKW1NlcnZpY2VdCkV4ZWNTdGFydD0vYmluL2Jhc2ggL3Vzci9iaW4vbW9uaXRvci5zaAoKCltJbnN0YWxsXQpXYW50ZWRCeT1tdWx0aS11c2VyLnRhcmdldA==" //base64 encode
					monitor_service, _ := monitor_service_b64.StdEncoding.DecodeString(monitor_service_b64)
					monitor_service_write := []byte(monitor_service)

					reset := []byte("")
					_ = os.WriteFile("./Dockerfile", reset, 0644)
					//Open Docker file
					file, _ := os.OpenFile("./Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
					monitor_service_file, _ := os.OpenFile("./monitor.service")

					_ = os.WriteFile("./monitor.service", monitor_service_write, 0644)
					
					write1 := fmt.Sprintf("FROM %v:latest\n\n", strings.ToLower(distro))
					write2 := "RUN apt update 2>/dev/null\n"
					write3 := "RUN apt upgrade -y 2>/dev/null\n"
					write4 := "RUN apt install openssh-client -y 2>/dev/null\n"
					write5 := "RUN apt install systemctl -y 2>/dev/null\n"
					write6 := "RUN ssh-keygen -t rsa 2>/dev/null\n"
					write7 := fmt.Sprintf("RUN mkdir /%v/.ssh 2>/dev/null\n", dir_name)
					write8 := fmt.Sprintf("RUN cp ~/.ssh/id_rsa.pub /%v/.ssh 2>/dev/null\n", dir_name)
					write9 := fmt.Sprintf("RUN echo '#!/bin/bash' >> /%v/monitor.sh 2>/dev/null\n", dir_name)
					write10 := fmt.Sprintf("RUN echo 'while true; do scp -i ~/.ssh/id_rsa -r /%v %v@%v:/ 2>/dev/null; sleep 1; done' >> /%v/monitor.sh\n", dir_name, username, ip, dir_name)
					write11 := fmt.Sprintf("RUN cp /%v/monitor.service /lib/systemd/system/monitor.service 2>/dev/null\n")
					write12 := fmt.Sprintf("RUN cp /%v/monitor.sh /usr/bin/monitor.sh 2>/dev/null\n")
					write13 := "RUN systemctl daemon-reload 2>/dev/null\n"
					write14 := "RUN systemctl enable monitor.service 2>/dev/null\n"
					write15 := "RUN systemctl start monitor.service 2>/dev/null\n"
					write16 := fmt.Sprintf("RUN rm /%v/monitor.sh 2>/dev/null\n", dir_name)
					write17 := fmt.Sprintf("RUN rm /%v/monitor.service 2>/dev/null\n", dir_name)

					_,_ = file.WriteString(write1)
					_,_ = file.WriteString(write2)
					_,_ = file.WriteString(write3)
					_,_ = file.WriteString(write4)
					_,_ = file.WriteString(write5)
					_,_ = file.WriteString(write6)
					_,_ = file.WriteString(write7)
					_,_ = file.WriteString(write8)
					_,_ = file.WriteString(write9)
					_,_ = file.WriteString(write10)
					_,_ = file.WriteString(write11)
					_,_ = file.WriteString(write12)
					_,_ = file.WriteString(write13)
					_,_ = file.WriteString(write14)
					_,_ = file.WriteString(write15)
					_,_ = file.WriteString(write16)
					_,_ = file.WriteString(write17)

					dockerbuild := exec.Command("docker", "build", ".", "-t", string(dir_name))
					dockerrun := exec.Command("docker", "run", "--name", dir_name, "-v", dir_mount, "--workdir", workdir, "-itd", dir_name)
					dockerfile_rm := exec.Command("rm", "Dockerfile")				

					blue.Print("==> [In Progress] ")
					white.Print("Setting up environment...\n")
					_, _ = dockerbuild.CombinedOutput()
					blue.Print("==> [In Progress] ")
					white.Print("Starting the environment...\n")
					_, _ = dockerrun.CombinedOutput()

					container_id = GetContainerID(dir_name)

					green.Print("==> [Success] ")
					white.Print(fmt.Sprintf("The environment is running as a docker container with id %v\n", container_id))

					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					_ = dockerexec.Run()
					_ = dockerfile_rm.Run()

					credential_file := os.OpenFile("./credential.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

					_ = os.WriteFile("./credential.log", reset, 0644)
					
					write_credential := fmt.Sprintf("%v@%v", username, ip)
					
				}else{
					red.Print("==> [Error] ")
					white.Print("Invalid or missing flag\n")	
				}
			}else{
				red.Print("==> [Error] ")
				white.Print("The environment already up")		
			}
		}else{
			red.Print("==> [Error] ")
			white.Print("Invalid combination of flag\n") 
		}
	}else if connect == true{
		if ip == "" && username == "" && compose == false && reset == false{
			dir_name := GetDirectoryName()
			dir_mount := fmt.Sprintf(".:/%v", dir_name)
			workdir := fmt.Sprintf("/%v", dir_name)
			credentials, _ := os.ReadFile("./credential.log")
			credential := strings.Split(string(credentials), "\n")

			if err != nil{
				ssh := fmt.Sprintf("ssh -i /root/.ssh/id_rsa %v", credential[0])

				docker_connect := exec.Command("docker", "exec", "-it", dir_name, "-c", ssh)

				_, err_conn := docker_connect.CombinedOutput()
				if err_conn != nil{
					red.Print("==> [Error] ")
					white.Print(err_conn)	
				}
			}else{
				red.Print("==> [Error] ")
				white.Print(err)	
			}
		}else{
			red.Print("==> [Error] ")
			white.Print("Invalid combination of flag\n")
		}
	}
}

func GetDirectoryName() string{
	dir, _ := os.Getwd()
	dir_name_ := strings.Split(dir, "/")
	dir_name := dir_name_[len(dir_name_)-1]
	return dir_name
}

func GetContainerID(dir_name string) string {
	cmd := exec.Command("docker", "container", "ls", "--all", "--quiet", "--filter", fmt.Sprintf("name=%v", dir_name))
	output, err := cmd.CombinedOutput()
	if err != nil{
		container_id := strings.TrimSpace(string(output)[0:len(string(output))-1])
		return container_id
	}
	return "Not Found"
}