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
	"os/exec"
	"github.com/spf13/cobra"
	"os"
	"strings"
	"io/ioutil"
	"bytes"
)

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Restore a backup",
	Run: func(cmd *cobra.Command, args []string) {
		name, _ := cmd.Flags().GetString("name")
		di, _ := cmd.Flags().GetString("di")
		
		check_docker := exec.Command("docker", "network", "ls")
		var stderr_check_docker bytes.Buffer
		check_docker.Stderr = &stderr_check_docker
		err := check_docker.Run()

		if err != nil{
			initLog("restore.go " + string(stderr_check_docker.Bytes()))
			red.Print("==> [Error] ")
			white.Print("Docker not running\n")
		}else{
			check_docker := exec.Command("docker", "network", "ls")
			var stderr_check_docker bytes.Buffer
			check_docker.Stderr = &stderr_check_docker
	
			if err != nil{
				initLog("restore.go " + string(stderr_check_docker.Bytes()))
				red.Print("==> [Error] ")
				white.Print("Docker not running\n")
			}else{
				if name == ""{
					if di == ""{
						red.Print("==> [Error] ")
						white.Print("Missing di flag\n")
					}else{
						_, err := os.Stat(fmt.Sprintf("./backup/%v-original.tar", di))
						if err != nil{
							initLog("restore.go " + fmt.Sprintf("%v", err))
							red.Print("==> [Error] ")
							white.Print("Environment or backup file not exist\n")
						}else{
							var proceed string;
							var f string;
							red.Print("==> [Caution] ")
							white.Print("Backup file name: ")
							fmt.Scanln(&f)
							if f != ""{
								_, err := os.Stat(fmt.Sprintf("./backup/%v", f))
								for err != nil{
									fmt.Println("Invalid input, please try again")
									fmt.Scanln(&f)
								}
								red.Print("==> [Caution] ")
								white.Print("This will restore your environment's current state to previous state. Action cannot be reverted.\nAre you sure to proceed? (y/n)")
								fmt.Scanln(&proceed)
								for proceed != "y" && proceed != "Y" && proceed != "n" && proceed != "N"{
									fmt.Println("Invalid input, please try again")
									fmt.Scanln(&proceed)
								}
								if proceed == "y" || proceed == "Y"{
									Revert(name)
								}			
							}else{
								red.Print("==> [Caution] ")
								white.Print("This will default restore your environment's current state to original state. Action cannot be reverted.\nAre you sure to proceed? (y/n)")
								for proceed != "y" && proceed != "Y" && proceed != "n" && proceed != "N"{
									fmt.Scanln(&proceed)
								}
								if proceed == "y" || proceed == "Y"{
									Revert(fmt.Sprintf("%v-original.tar", di))
								}			
							}
						}
					}
				}else{
					_, err := os.Stat(fmt.Sprintf("./backup/%v", name))
					if err != nil{
						initLog("restore.go " + fmt.Sprintf("%v", err))
						red.Print("==> [Error] ")
						white.Print("Backup file not found\n")
					}else{
						var proceed string;
						red.Print("==> [Caution] ")
						white.Print("This will restore your environment's current state to previous state. Action cannot be reverted.\nAre you sure to proceed? (y/n)")
						fmt.Scanln(&proceed)
						for proceed != "y" && proceed != "Y" && proceed != "n" && proceed != "N"{
							fmt.Println("Invalid input, please try again")
							fmt.Scanln(&proceed)
						}
						if proceed == "y" || proceed == "Y"{
							Revert(name)
						}	
					}
				}		
			}
		}		
	},
}

func Revert(name string){
	save_name := fmt.Sprintf(fmt.Sprintf("./backup/%v", name))
	dockerrevert := exec.Command("docker", "image", "load", "-i", save_name)
	var stderr_dockerrevert bytes.Buffer
	dockerrevert.Stderr = &stderr_dockerrevert

	imagename := strings.Split(name, "-")
	container_id := GetContainerID(imagename[0])

	blue.Print("==> [In Progress] ")
	white.Print("Shutting down the environment...\n")		

	dockerstop := exec.Command("docker", "stop", container_id)
	var stderr_dockerstop bytes.Buffer
	dockerstop.Stderr = &stderr_dockerstop
	_ = dockerstop.Run()

	initLog("restore.go " + string(stderr_dockerstop.Bytes()))

	blue.Print("==> [In Progress] ")
	white.Print("Reverting to chosen backup file...\n")		
	
	_ = dockerrevert.Run()

	initLog("restore.go " + string(stderr_dockerrevert.Bytes()))

	blue.Print("==> [In Progress] ")
	white.Print("Booting up the environment...\n")		

	dat, err1 := os.ReadFile("./run/connect")
	
	_, err2 := ioutil.ReadDir("./run")


	if err2 != nil{
		initLog("restore.go " + fmt.Sprintf("%v", err2))
		red.Print("==> [Error] ")
		white.Print("./run missing")
	}else{
		if err1 != nil{
			initLog("restore.go " + fmt.Sprintf("%v", err1))
			red.Print("==> [Error] ")
			white.Print("Missing file named connect in ./run")
		}else{
			data := strings.Split(string(dat), " ")
			dockerrun := exec.Command(data[0], data[1:]...)
			var stderr_dockerrun bytes.Buffer
			dockerrun.Stderr = &stderr_dockerrun
			_ = dockerrun.Run()
		}
	}
}

func init() {
	dockerCmd.AddCommand(restoreCmd)

	restoreCmd.PersistentFlags().String("name", "", "Name of backup file")
	restoreCmd.PersistentFlags().String("di", "", "Name of docker environment image")

}
