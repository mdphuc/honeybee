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
	"github.com/fatih/color"
	"os/exec"
	"strings"
	// "slices"
	// "encoding/base64"
	"os"
)

var validargs = []string{
	"target",
	"environment",
	"pkgmanager",
	"library",
	"ip",
	"user",
	"distro",
	"os",
}

var supported_distro = []string{"ubuntu", "debian", "fedora", "opensuse"}

var supported_pkg = []string{"apt", "dnf", "yum", "zypper", "pacman"}

var red = color.New(color.FgRed, color.Bold)
var white = color.New(color.FgWhite, color.Bold)
var green = color.New(color.FgGreen, color.Bold)
var blue = color.New(color.FgBlue, color.Bold)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run remote development environment",
	Run: func(cmd *cobra.Command, args []string) {

		target, _ := cmd.Flags().GetString("target")
		environment, _ := cmd.Flags().GetString("environment")
		pkgmanager, _ := cmd.Flags().GetString("pkgmanager")
		library, _ := cmd.Flags().GetString("library")
		ip,_ := cmd.Flags().GetString("ip")
		user,_ := cmd.Flags().GetString("user")
		distro,_ := cmd.Flags().GetString("distro")
		os_ := cmd.Flags().Lookup("os").Changed

		if target != "docker" && target != "remote_machine"{
			color.Red("[Error] Invalid target")
		}else if target == "docker"{
			Docker(environment, pkgmanager, library, os_, distro)
		}else{
			RemoteMachine(environment, pkgmanager, library, ip, user)
		}
	},
	ValidArgs: validargs,
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.PersistentFlags().String("target", "", "Target mode (docker, remote machine)")
	runCmd.PersistentFlags().String("environment", "Linux", "Operating system for development environment (Windows, Linux)")
	runCmd.PersistentFlags().String("pkgmanager", "", "Package manager (Use for Linux)")
	runCmd.PersistentFlags().String("distro", "", "Linux distro")
	runCmd.PersistentFlags().String("library", "", "Library to pre-install (separate by comma)")
	runCmd.PersistentFlags().String("ip", "", "IP for remote machine")
	runCmd.PersistentFlags().String("user", "", "Username for remote machine")
	runCmd.PersistentFlags().BoolP("os", "", true, "Supported OS for Docker environment and their package manager")
}

func Docker(environment string, pkgmanager string, library string, os_ bool, distro string){
	if os_ == true {
		if environment != " " && pkgmanager != " " && distro != " "{
			white.Println("Supported operating system for Docker environment")
			white.Println(`
apt:
	Ubuntu
	Debian

dnf:
	Fedora

yum:
	Fedora

zypper:
	openSUSE

pacman:
	ArchLinux
			`)
		}else{
			red.Print("==> [Error] ")
			white.Print("Invalid combination of flag\n")
		}
	}else{
		if strings.ToLower(environment) != "windows" && strings.ToLower(environment) != "linux"{
			red.Print("==> [Error] ")
			white.Print("Missing or invalid flag\n")
		}else if strings.ToLower(environment) == "linux"{
			check_flag := CheckFlag(pkgmanager, distro, library, "Linux")
			if check_flag == true{
				dir_name := GetDirectoryName()
				dir_mount := fmt.Sprintf(".:/%v", dir_name)
				workdir := fmt.Sprintf("/%v", dir_name)

				reset := []byte("")
				_ = os.WriteFile("./docker.ps1", reset, 0644)
			
				file, _ := os.OpenFile("./docker.ps1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
				

				dockerbuild := exec.Command("docker", "build", ".", "-t", string(dir_name))
				dockerrun := exec.Command("docker", "run", "--name", dir_name, "-v", dir_mount, "--workdir", workdir, "-itd", dir_name)

				write1 := fmt.Sprintf("Invoke-Command {docker exec -it %v /bin/bash}", dir_name)

				_,_ = file.WriteString(write1)

				blue.Print("==> [In Progress] ")
				white.Print("Setting up environment...\n")
				_, _ := dockerbuild.CombinedOutput()
				blue.Print("==> [In Progress] ")
				white.Print("Starting the environment...\n")
				_ , _ := dockerrun.CombinedOutput()

				container_id := GetContainerID(dir_name)

				green.Print("==> [Success] ")
				white.Print(fmt.Sprintf("The environment is running as a docker container with id %v\n", container_id))

				dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
				dockerexec.Stdin = os.Stdin
				dockerexec.Stdout = os.Stdout
				dockerexec.Stderr = os.Stderr
			
				_ := dockerexec.Run()
			}
			
		}else if strings.ToLower(environment) == "windows"{
			check_flag := CheckFlag(pkgmanager, distro, library, "Windows")
			if check_flag == true {
				dir_name := GetDirectoryName()
				dir_mount := fmt.Sprintf(".:/%v", dir_name)
				workdir := fmt.Sprintf("/%v", dir_name)

				reset := []byte("")
				_ = os.WriteFile("./docker.ps1", reset, 0644)

				file, _ := os.OpenFile("./docker.ps1", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

				dockerbuild := exec.Command("docker", "build", ".", "-t", string(dir_name))
				dockerrun := exec.Command("docker", "run", "--name", dir_name, "-v", dir_mount, "--workdir", workdir, "-itd", dir_name)

				write1 := fmt.Sprintf("Invoke-Command {docker exec -it %v /bin/bash}", dir_name)

				_,_ = file.WriteString(write1)
				

				blue.Print("==> [In Progress] ")
				white.Print("Setting up environment...\n")
				dockerbuild_output, _ := dockerbuild.CombinedOutput()
				blue.Print("==> [In Progress] ")
				white.Print("Starting the environment...\n")
				dockerrun_output , _ := dockerrun.CombinedOutput()

				container_id := GetContainerID(dir_name)

				green.Print("==> [Success] ")
				white.Print(fmt.Sprintf("The environment is running as a docker container with id %v\n", container_id))

				dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
				dockerexec.Stdin = os.Stdin
				dockerexec.Stdout = os.Stdout
				dockerexec.Stderr = os.Stderr
			
				_ := dockerexec.Run()
			}
		}
	}
}

func RemoteMachine(environment string, pkgmanager string, library string, ip string, user string){

}

func GetContainerID(dir_name string) string {
	cmd := exec.Command("docker", "container", "ls", "--all", "--quiet", "--filter", fmt.Sprintf("name=%v", dir_name))
	output, _ := cmd.CombinedOutput()
	container_id := strings.TrimSpace(string(output)[0:len(string(output))-1])
	return container_id
}

func GetDirectoryName() string{
	dir, _ := os.Getwd()
	dir_name_ := strings.Split(dir, "/")
	dir_name := dir_name_[len(dir_name_)-1]
	return dir_name
}

func CheckValidDistro(distro string) bool{
	for _, o := range supported_distro{
		if strings.ToLower(distro) == o{
			return true
		}
	} 
	return false
}

func CheckValidPkg(pkgmanager string) bool{
	for _, o := range supported_pkg{
		if pkgmanager == o{
			return true
		}
	} 
	return false
}

func CheckMatch(pkgmanager string, distro string) bool {
	if pkgmanager == "apt" && strings.ToLower(distro) == "ubuntu"{
		return true
	}
	if pkgmanager == "apt" && strings.ToLower(distro) == "debian"{
		return true
	}
	if pkgmanager == "dnf" && strings.ToLower(distro) == "fedora"{
		return true
	}
	if pkgmanager == "yum" && strings.ToLower(distro) == "fedora"{
		return true
	}
	if pkgmanager == "zypper" && strings.ToLower(distro) == "opensuse"{
		return true
	}
	return false
}

func apt(distro string, target string, libraries string){
	reset := []byte("")
	_ = os.WriteFile("./Dockerfile", reset, 0644)
	//Open Docker file
	file, _ := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	 
	write1 := fmt.Sprintf("FROM %v:latest\n\n", strings.ToLower(distro))
	write2 := "RUN apt update 2>/dev/null\n"
	write3 := "RUN apt upgrade -y 2>/dev/null\n"

	_,_ = file.WriteString(write1)
	_,_ = file.WriteString(write2)
	_,_ = file.WriteString(write3)


	if strings.ToLower(target) == "windows"{
		write4 := "\nRUN apt install -y wget apt-transport-https software-properties-common 2>/dev/null\n"
		write5 := "RUN wget -q 'https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb' 2>/dev/null\n"
		write6 := "RUN dpkg -i packages-microsoft-prod.deb 2>/dev/null\n"
		write7 := "rm packages-microsoft-prod.deb 2>/dev/null\n"
		write8 := "RUN apt update 2>/dev/null\n"
		write9 := "RUN apt install -y powershell 2>/dev/null\n"

		_,_ = file.WriteString(write4)
		_,_ = file.WriteString(write5)
		_,_ = file.WriteString(write6)
		_,_ = file.WriteString(write7)
		_,_ = file.WriteString(write8)
		_,_ = file.WriteString(write9)
	}
	
	if libraries != ""{
		library := strings.Split(libraries, ",")
		for _, o := range library{
			writelibrary := fmt.Sprintf("RUN apt install %v 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

}

func dnf(distro string, target string, libraries string){
	reset := []byte("")
	_ = os.WriteFile("./Dockerfile", reset, 0644)
	//Open Docker file
	file, _ := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	write1 := "FROM fedora:latest\n\n"
	write2 := "RUN dnf update -y 2>/dev/null\n"
	write3 := "RUN dnf upgrade -y 2>/dev/null\n"

	_,_ = file.WriteString(write1)
	_,_ = file.WriteString(write2)
	_,_ = file.WriteString(write3)

	if strings.ToLower(target) == "windows"{
		write4 := "\nRUN rpm --import https://packages.microsoft.com/keys/microsoft.asc 2>/dev/null\n"
		write5 := "RUN dnf install curl -y 2>/dev/null\n"
		write6 := "RUN curl https://packages.microsoft.com/config/rhel/7/prod.repo | tee /etc/yum.repos.d/microsoft.repo 2>/dev/null\n"
		write7 := "RUN dnf makecache 2>/dev/null\n"
		write8 := "RUN dnf install powershell -y 2>/dev/null\n"

		 
		_,_ = file.WriteString(write4)
		_,_ = file.WriteString(write5)
		_,_ = file.WriteString(write6)
		_,_ = file.WriteString(write7)
		_,_ = file.WriteString(write8)
	}

	if libraries != ""{
		library := strings.Split(libraries, ",")
		for _, o := range library{
			writelibrary := fmt.Sprintf("RUN dnf install %v 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

}

func yum(distro string, target string, libraries string){
	reset := []byte("")
	_ = os.WriteFile("./Dockerfile", reset, 0644)
	//Open Docker file
	file, _ := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	write1 := "FROM fedora:latest\n\n"
	write2 := "RUN yum update -y 2>/dev/null\n"
	write3 := "RUN yum upgrade 2>/dev/null\n"

	_,_ = file.WriteString(write1)
	_,_ = file.WriteString(write2)
	_,_ = file.WriteString(write3)

	if strings.ToLower(target) == "windows"{
		write4 := "\nRUN rpm --import https://packages.microsoft.com/keys/microsoft.asc 2>/dev/null\n"
		write5 := "RUN yum install curl -y 2>/dev/null\n"
		write6 := "RUN curl https://packages.microsoft.com/config/rhel/7/prod.repo | tee /etc/yum.repos.d/microsoft.repo 2>/dev/null\n"
		write7 := "RUN yum makecache 2>/dev/null\n"
		write8 := "RUN yum install powershell -y 2>/dev/null\n"

		_,_ = file.WriteString(write4)
		_,_ = file.WriteString(write5)
		_,_ = file.WriteString(write6)
		_,_ = file.WriteString(write7)
		_,_ = file.WriteString(write8)
	}	

	if libraries != ""{
		library := strings.Split(libraries, ",")
		for _, o := range library{
			writelibrary := fmt.Sprintf("RUN yum install %v 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

}

func zypper(distro string, target string, libraries string){
	reset := []byte("")
	_ = os.WriteFile("./Dockerfile", reset, 0644)
	//Open Docker file
	file, _ := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	 
	write1 := "FROM opensuse/leap:latest\n\n"
	write2 := "RUN zypper update -y 2>/dev/null\n"
	write3 := "RUN zypper install gzip curl tar libicu60_2 libopenssl1_0_0 -y 2>/dev/null\n"

	_,_ = file.WriteString(write1)
	_,_ = file.WriteString(write2)
	_,_ = file.WriteString(write3)

	if strings.ToLower(target) == "windows"{
		write4 := "RUN curl -L https://github.com/PowerShell/PowerShell/releases/download/v6.1.3/powershell-6.1.3-linux-x64.tar.gz -o /tmp/powershell.tar.gz 2>/dev/null\n"
		write5 := "RUN mkdir -p /opt/microsoft/powershell\n"
		write6 := "RUN gzip -d /tmp/powershell.tar.gz\n"
		write7 := "RUN tar -xf /tmp/powershell.tar -C /opt/microsoft/powershell\n"
		write8 := "RUN ln -s /opt/microsoft/powershell/pwsh /usr/bin/pwsh\n"
		write9 := "RUN chmod +x /usr/bin/pwsh\n"

		_,_ = file.WriteString(write4)
		_,_ = file.WriteString(write5)
		_,_ = file.WriteString(write6)
		_,_ = file.WriteString(write7)
		_,_ = file.WriteString(write8)
		_,_ = file.WriteString(write9)
	}

	if libraries != ""{
		library := strings.Split(libraries, ",")
		for _, o := range library{
			writelibrary := fmt.Sprintf("RUN zypper install %v 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

}

func CheckFlag(pkgmanager string, distro string, library string, target string) bool{
	if distro == "" && pkgmanager == ""{
		apt(distro, target, library)
		return true
	}else if distro != "" && pkgmanager == ""{
		validity := CheckValidDistro(distro)
		if validity == true{
			if strings.ToLower(distro) == "ubuntu" || strings.ToLower(distro) == "debian"{
				apt(distro, target, library)
				return true
			}else if strings.ToLower(distro) == "fedora"{
				dnf(distro, target, library)
				return true
			}else if strings.ToLower(distro) == "opensuse"{
				zypper(distro, target, library)
				return true
			}
		}else{
			red.Print("==> [Error] ")
			white.Print("Distro not supported\n")
		}
	}else if distro == "" && pkgmanager != ""{
		validity := CheckValidPkg(pkgmanager)
		if validity == true{
			if pkgmanager == "apt"{
				apt("Ubuntu", target, library)
				return true
			}else if pkgmanager == "dnf"{
				dnf("Fedora", target, library)
				return true
			}else if pkgmanager == "yum"{
				yum("Fedora", target, library)
				return true
			}else if pkgmanager == "zypper"{
				zypper("openSUSE", target, library)
				return true
			}
		}else{
			red.Print("==> [Error] ")
			white.Print("Package Manager not supported\n")
		}
	}else{
		checkmatch := CheckMatch(pkgmanager, distro)
		if checkmatch == false{
			red.Print("==> [Error] ")
			white.Print("Selected Package Manager is not supported by selected distro\n")
		}
		validity_pkg := CheckValidPkg(pkgmanager)
		validity_distro := CheckValidDistro(distro)
		if validity_pkg == false{
			red.Print("==> [Error] ")
			white.Print("Package Manager not supported\n")
		} 
		if validity_distro == false{
			red.Print("==> [Error] ")
			white.Print("Distro not supported\n")
		}
		if validity_distro == true && validity_pkg == true && checkmatch == true{
			if strings.ToLower(distro) == "ubuntu" || strings.ToLower(distro) == "debian"{
				apt(distro, target, library)
				return true
			}else if strings.ToLower(distro) == "fedora"{
				dnf(distro, target, library)
				return true
			}else if strings.ToLower(distro) == "opensuse"{
				zypper(distro, target, library)
				return true
			}
		}
	}
	return false
}
