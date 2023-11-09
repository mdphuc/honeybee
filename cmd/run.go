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

var supported_distro = []string{"Ubuntu", "Debian", "CentOS", "Fedora", "openSUSE", "ArchLinux"}

var supported_pkg = []string{"apt", "dnf", "yum", "zypper", "pacman"}

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
			fmt.Println("Supported operating system for Docker environment")
			fmt.Println(`
apt:
	Ubuntu
	Debian

dnf:
	Fedora

yum:
	CentOS

zypper:
	openSUSE

pacman:
	ArchLinux
			`)
		}else{
			color.Red("[Error] Invalid combination of flag")
		}
	}else{
		if environment != "Windows" && environment != "Linux"{
			color.Red("[Error] Missing or invalid flag")
		}else if environment == "Linux"{
			CheckFlag(pkgmanager, distro, library, "Linux")
		}else if environment == "Windows"{
			CheckFlag(pkgmanager, distro, library, "Windows")
		}
	}

}

func RemoteMachine(environment string, pkgmanager string, library string, ip string, user string){

}

// func RunCommand(cmd){
// 	err := cmd.Run()
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	fmt.Println(args)
// }

func GetContainerID(dir_name string) string {
	cmd := exec.Command("docker", "container", "ls", "--all", "--quiet", "--filter", fmt.Sprintf("name=%v", dir_name))
	output, _ := cmd.CombinedOutput()
	container_name := strings.TrimSpace(string(output)[0:len(string(output))-1])
	return container_name
}

func CheckValidDistro(distro string) bool{
	for _, o := range supported_distro{
		if distro == o{
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
	if pkgmanager == "apt" && distro == "Ubuntu"{
		return true
	}
	if pkgmanager == "apt" && distro == "Debian"{
		return true
	}
	if pkgmanager == "dnf" && distro == "Fedora"{
		return true
	}
	if pkgmanager == "yum" && distro == "CentOS"{
		return true
	}
	if pkgmanager == "zypper" && distro == "openSUSE"{
		return true
	}
	if pkgmanager == "pacman" && distro == "Arch"{
		return true
	}
	return false
}

func writeDockerfile(pkgmanager string, distro string, library string, target string){
	reset := []byte("")
	_ = os.WriteFile("./Dockerfile", reset, 0644)
	//Open Docker file
	file, _ := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	// Define string to write into Dockerfile
	write1 := fmt.Sprintf("FROM %v:latest\n\n", strings.ToLower(distro))
	write2 := fmt.Sprintf("RUN %v update 2>/dev/null\n", strings.ToLower(pkgmanager))
	write3 := fmt.Sprintf("RUN %v upgrade 2>/dev/null\n", strings.ToLower(pkgmanager))

	//Write into Dockerfile
	_,_ = file.WriteString(write1)
	_,_ = file.WriteString(write2)
	_,_ = file.WriteString(write3)

	if library != ""{
		library_arr := strings.Split(library, ",")
		for _, o := range library_arr{
			write := fmt.Sprintf("RUN %v install %v 2>/dev/null\n", strings.ToLower(pkgmanager), o)
			_,_ = file.WriteString(write)
		}
	}

	if target == "Windows"{
		if pkgmanager == "apt"{
			write4 := "\nRUN apt install -y wget apt-transport-https software-properties-common 2>/dev/null"
			write5 := "RUN wget -q 'https://packages.microsoft.com/config/ubuntu/22.04/packages-microsoft-prod.deb' 2>/dev/null"
			write6 := "RUN dpkg -i packages-microsoft-prod.deb 2>/dev/null"
			write7 := "rm packages-microsoft-prod.deb 2>/dev/null"
			write8 := "RUN apt update 2>/dev/null"
			write9 := "RUN apt install -y powershell 2>/dev/null"

			_,_ = file.WriteString(write4)
			_,_ = file.WriteString(write5)
			_,_ = file.WriteString(write6)
			_,_ = file.WriteString(write7)
			_,_ = file.WriteString(write8)
			_,_ = file.WriteString(write9)

		}else if pkgmanager == "dnf"{
			write4 := "\nRUN rpm --import https://packages.microsoft.com/keys/microsoft.asc 2>/dev/null"
			Write5 := "RUN dnf install curl -y 2>/dev/null"
			Write6 := "RUN curl https://packages.microsoft.com/config/rhel/7/prod.repo | tee /etc/yum.repos.d/microsoft.repo 2>/dev/null"
			Write7 := "RUN dnf makecache 2>/dev/null"
			Write8 := "RUN dnf install powershell -y 2>/dev/null"

			_,_ = file.WriteString(write4)
			_,_ = file.WriteString(write5)
			_,_ = file.WriteString(write6)
			_,_ = file.WriteString(write7)
			_,_ = file.WriteString(write8)

		}else if pkgmanager == "yum"{
			write4 := "RUN yum install curl -y 2>/dev/null"
			write5 := "\nRUN curl -sSL -O https://packages.microsoft.com/config/rhel/7/packages-microsoft-prod.rpm 2>/dev/null"
			write6 := "RUN rpm -i packages-microsoft-prod.rpm 2>/dev/null"
			write7 := "RUN rm packages-microsoft-prod.rpm 2>/dev/null"
			write8 := "RUN yum update 2>/dev/null"
			write9 := "RUN yum install powershell -y 2>/dev/null"

			_,_ = file.WriteString(write4)
			_,_ = file.WriteString(write5)
			_,_ = file.WriteString(write6)
			_,_ = file.WriteString(write7)
			_,_ = file.WriteString(write8)
			_,_ = file.WriteString(write9)
		}else if pkgmanager == "pacman"{
			write4 := "RUN pacman -S git -y"
			write5 := "RUN git clone https://aur.archlinux.org/powershell.git"
			write6 := "RUN cd powershell-bin"
			write7 := "RUN makepkg -si"

			_,_ = file.WriteString(write4)
			_,_ = file.WriteString(write5)
			_,_ = file.WriteString(write6)
			_,_ = file.WriteString(write7)
		}
	}
}

func CheckFlag(pkgmanager string, distro string, library string, target string){
	if distro == "" && pkgmanager == ""{
		writeDockerfile("apt", "Ubuntu", library, target)
	}else if distro != "" && pkgmanager == ""{
		validity := CheckValidDistro(distro)
		if validity == true{
			if distro == "Ubuntu" || distro == "Debian"{
				writeDockerfile("apt", distro, library, target)
			}else if distro == "CentOS"{
				writeDockerfile("yum", distro, library, target)
			}else if distro == "Fedora"{
				writeDockerfile("dnf", distro, library, target)
			}else if distro == "openSUSE"{
				writeDockerfile("zypper", distro, library, target)
			}else if distro == "ArchLinux"{
				writeDockerfile("pacman", distro, library, target)
			}
		}else{
			color.Red("[Error] Distro not supported")
		}
	}else if distro == "" && pkgmanager != ""{
		validity := CheckValidPkg(pkgmanager)
		if validity == true{
			if pkgmanager == "apt"{
				writeDockerfile(pkgmanager, "Ubuntu", library, target)
			}else if pkgmanager == "dnf"{
				writeDockerfile(pkgmanager, "Fedora", library, target)
			}else if pkgmanager == "yum"{
				writeDockerfile(pkgmanager, "CentOS", library, target)
			}else if pkgmanager == "zypper"{
				writeDockerfile(pkgmanager, "openSUSE", library, target)
			}else{
				writeDockerfile(pkgmanager, "ArchLinux", library, target)
			}
		}else{
			color.Red("[Error] Package Manager not supported")
		}
	}else{
		checkmatch := CheckMatch(pkgmanager, distro)
		if checkmatch == false{
			color.Red("[Error] Selected Package Manager is not supported by selected distro")
		}else{
			writeDockerfile(pkgmanager, distro, library, target)
		}
		validity_pkg := CheckValidPkg(pkgmanager)
		validity_distro := CheckValidDistro(distro)
		if validity_pkg == false{
			color.Red("[Error] Package Manager not supported")
		} 
		if validity_distro == false{
			color.Red("[Error] Distro not supported")
		}
	}
}