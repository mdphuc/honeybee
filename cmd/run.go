/*
Copyright © 2023 Phuc Mai 0xmdphuc@gmail.com

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
	"os"
	"strings"
	"os/exec"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run Set up the environment",

}

func init() {
	rootCmd.AddCommand(runCmd)
}


func GetDirectoryName() string{
	dir, _ := os.Getwd()
	dir_name_ := strings.Split(dir, "/")
	dir_name := dir_name_[len(dir_name_)-1]
	return dir_name
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

