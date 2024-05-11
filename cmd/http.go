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
	"bufio"
	"github.com/spf13/cobra"
	"strings"
	"os"
	"slices"
	"strconv"
	// "log"
	"bytes"
	"io/ioutil"
	"os/exec"
	"time"
)

// httpCmd represents the http command
var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "Web UI for beaver",
	Long: "Caution: Only work with Ubuntu/Debian deployed machine",
	Run: func(cmd *cobra.Command, args []string) {
		file, _ := os.Open("./config/default.conf")
		defer file.Close()
		
		var port_check, terminal_check, exec_check, password_check, access_check, ip_check = true, true, true, true, true, true
		var PORT = 0
		var IP = ""
		
		check_docker := exec.Command("docker", "network", "ls")

		var stderr_check_docker bytes.Buffer
		check_docker.Stderr = &stderr_check_docker

		var docker_bool = true

		err := check_docker.Run()
		// fmt.Sprintf("%v",err)
		if err != nil{
			red.Print("==> [Error] ")
			white.Print(fmt.Sprintf("Docker not running\n"))		
			initLog("http.go " + "Docker not running")
			docker_bool = false
		}

		if docker_bool{
			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines) 
			line_number := 0
			for scanner.Scan(){
				line_number += 1
				line := scanner.Text()

				line_rune := []rune(line)
				check_comment := CheckComment(string(line_rune[0:1]))

				if strings.Split(line, " ")[0] != ""{
					if !check_comment{
						keywords := strings.Split(line, "=")
						check_keyword := CheckKeyword(keywords[0])
						if !check_keyword {
							red.Print("==> [Error] ")
							white.Print(fmt.Sprintf("Syntax error at line %v. Keyword not found: %v\n", line_number, keywords[0]))		
							initLog("http.go " + fmt.Sprintf("Syntax error at line %v. Keyword not found: %v", line_number, keywords[0]))
						}else{
							if strings.Split(keywords[1], " ")[0] == ""{
								red.Print("==> [Error] ")
								white.Print(fmt.Sprintf("Syntax error at line %v. There's space after = sign\n", line_number))		
								initLog("http.go " + fmt.Sprintf("Syntax error at line %v. There's space after = sign", line_number))
								if strings.Split(keywords[0], " ")[1] == ""{
									red.Print("==> [Error] ")
									white.Print(fmt.Sprintf("Syntax error at line %v. There's space after = sign\n", line_number))		
									initLog("http.go " + fmt.Sprintf("Syntax error at line %v. There's space before = sign", line_number))
								}
							}else{
								if keywords[0] == "Port"{
									i_p, err := strconv.Atoi(keywords[1])
									if err != nil{
										red.Print("==> [Error] ")
										white.Print(fmt.Sprintf("Port must be positive integer < 65535\n"))			
										initLog("http.go " + "Port must be positive integer < 65535")
										port_check = false
									}else if i_p < 1 || i_p > 65535{
										red.Print("==> [Error] ")
										white.Print(fmt.Sprintf("Port must be positive integer < 65535\n"))			
										initLog("http.go " + "Port must be positive integer < 65535")
										port_check = false
									}else{
										PORT = i_p
									}
								}else{
									PORT = 5500
								}
								if keywords[0] == "AllowTerminalOnWeb"{
									if keywords[1] != "TRUE" && keywords[1] != "FALSE"{
										red.Print("==> [Error] ")
										white.Print(fmt.Sprintf("AllowTerminalOnWeb must be TRUE or FALSE\n"))			
										initLog("http.go " + "AllowTerminalOnWeb must be TRUE or FALSE")
										terminal_check = false
									}
								}
								if keywords[0] == "Exec"{
									_, err := ioutil.ReadDir(keywords[1])
									if err != nil{
										red.Print("==> [Error] ")
										white.Print(fmt.Sprintf("Invalid path\n"))			
										initLog("http.go " + fmt.Sprintf("Invalid path"))
										exec_check = false
									}
								}else{
									red.Print("==> [Error] ")
									white.Print(fmt.Sprintf("Exec path not found\n"))			
									initLog("http.go " + "Exec path not found")
									exec_check = false
								}
								if keywords[0] == "PasswordProtected"{
									if keywords[1] != "TRUE" && keywords[1] != "FALSE"{
										red.Print("==> [Error] ")
										white.Print(fmt.Sprintf("PasswordProtected must be TRUE or FALSE\n"))			
										initLog("http.go " + "PasswordProtected must be TRUE or FALSE")
										password_check = false
									}
								}
								if keywords[0] == "AllowAccess"{
									check_allow_access := CheckValidIP(keywords[1])
									if !check_allow_access{
										access_check = false
									}					
								}
								if keywords[0] == "IP"{
									check_IP := CheckValidIP(keywords[1])
									if !check_IP{
										ip_check = false
									}else{
										IP = keywords[0]
									}
								}else{
									IP = "127.0.0.1"
								}
								
							}

						}
					// fmt.Println(check_keyword)
					}	
				}
			}
		}
		if port_check && terminal_check && exec_check && password_check && access_check && ip_check{
			bind := fmt.Sprintf("%v:%v", IP, PORT)
			webui := exec.Command("gunicorn", "--chdir", "./web", "--bind", bind, "-w", "3", "wsgi:app")
			webui.Stderr = os.Stderr
			webui.Stdout = os.Stdout
			webui.Stdin = os.Stdin

			var stderr_webui bytes.Buffer
			webui.Stderr = &stderr_webui

			initLog("http.go " + fmt.Sprintf("%v", stderr_webui))
			initLog("http.go " + "httplaunch")
			
			blue.Print("==> [In Progress] ")
			white.Print("Spinning up WebUI...\n")

			time.Sleep(2 * time.Second)

			green.Print("==> [Success] ")
			white.Print(fmt.Sprintf("Finished spinning up environment http://%v:%v...", IP, PORT))

			err := webui.Run()	

			if err != nil{
				red.Print("==> [Error] ")
				white.Print(fmt.Sprintf("WebUI cannot launched. See ./log/app.log for more detail\n"))			
			}
		} 
	},
}

func CheckValidIP(ips string) bool {
	ips_list := strings.Split(ips, " ")
	for _, ip := range ips_list{
		if ip != ""{
			ip_frags := strings.Split(ip, ".")
			if len(ip_frags) != 4{
				red.Print("==> [Error] ")
				white.Print("Invalid IP\n")
				initLog("http.go " + "Invalid IP")
				return false
			}else{
				for i, ip_frag := range ip_frags{
					ip_f_i, err := strconv.Atoi(ip_frag)
					if err != nil{
						red.Print("==> [Error] ")
						white.Print("Invalid IP\n")
						initLog("http.go " + "Invalid IP")
						return true
					}else{
						if i == 0{
							if ip_f_i == 0{
								red.Print("==> [Error] ")
								white.Print("Invalid IP\n")
								initLog("http.go " + "Invalid IP")
								return false	
							}
							if !(ip_f_i <= 255 && ip_f_i > 0){
								red.Print("==> [Error] ")
								white.Print("Invalid IP\n")
								initLog("http.go " + "Invalid IP")
								return false		
							}
						}else{
							if !(ip_f_i <= 255 && ip_f_i >= 0){
								red.Print("==> [Error] ")
								white.Print("Invalid IP\n")
								initLog("http.go " + "Invalid IP")
								return false		
							}
						}
					}
				}
			}
		}
	}
	return true
}

func CheckKeyword(keyword string) bool {
	keywords := []string{"Port", "AllowTerminalOnWeb", "Exec", "PasswordProtected", "AllowAccess", "IP"}
	if slices.Contains(keywords, keyword){
		return true
	}
	return false
}

func CheckComment(hashtag string) bool {
	if hashtag == "#"{
		return true
	}
	return false
}

func init() {
	runCmd.AddCommand(httpCmd)
}
