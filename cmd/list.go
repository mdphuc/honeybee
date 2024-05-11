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
	"bytes"
	"github.com/spf13/cobra"
	"os"
	"io/ioutil"
	"os/exec"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list backup",
	Run: func(cmd *cobra.Command, args []string) {
		name,_ := cmd.Flags().GetString("name")
		detailed := cmd.Flags().Lookup("detail").Changed

		check_docker := exec.Command("docker", "network", "ls")

		var stderr bytes.Buffer
		check_docker.Stderr = &stderr
		
		err := check_docker.Run()

		if err != nil{
			red.Print("==> [Error] ")
			white.Print("Docker not running\n")
			initLog("list.go " + string(stderr.Bytes()))
		}else{
			if name == ""{
				files, err := ioutil.ReadDir("./backup")
				if detailed == false{
					if err != nil{
						red.Print("==> [Error] ")
						white.Print("No backup folder found\n")
					}else{
						for _, file := range files {
							fmt.Println(file.Name())
						}
					}
				}else{
					if err != nil{
						red.Print("==> [Error] ")
						white.Print("No backup folder found\n")
					}else{
						for _, file := range files {
							fileInfo, err_stat := os.Stat(fmt.Sprintf("./backup/%v",file.Name()))
							
							initLog("list.go " + fmt.Sprintf("%v", err_stat))

							modificationTime := fileInfo.ModTime() 
							fileSize := fileInfo.Size() 
							fmt.Println(file.Name())
							fmt.Println("Last modified time of the file: ", modificationTime)
							fmt.Println("File size: ", fileSize)
							fmt.Println("\n")
						}
					}
				}
			}else{
				fileInfo, err := os.Stat(fmt.Sprintf("./backup/%v", name))
				if err != nil{
					red.Print("==> [Error] ")
					white.Print("Problem reading file stat\n")
					initLog("list.go " + fmt.Sprintf("%v", err))
				}else{
					modificationTime := fileInfo.ModTime() 
					fileSize := fileInfo.Size() 
					fmt.Println(name)
					fmt.Println("Last modified time of the file: ", modificationTime)
					fmt.Println("File size: ", fileSize)
					fmt.Println("\n")
				}
			}
	
		}	
	},
}

func init() {
	dockerCmd.AddCommand(listCmd)

	listCmd.PersistentFlags().String("name", "", "Name of backup file")
	dockerCmd.PersistentFlags().BoolP("detail", "", true, "Detailed backup file")

}
