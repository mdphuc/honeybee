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
	"bytes"
	"github.com/spf13/cobra"
	"github.com/fatih/color"
	"os/exec"
	"strings"
	// "slices"
	// "encoding/base64"
	"os"
	"strconv"

)

type Flag struct{
	flag string
	name string
	iter int
}

var red = color.New(color.FgRed, color.Bold)
var white = color.New(color.FgWhite, color.Bold)
var green = color.New(color.FgGreen, color.Bold)
var blue = color.New(color.FgBlue, color.Bold)

var supported_distro = []string{"ubuntu", "debian", "fedora", "opensuse"}

var supported_pkg = []string{"apt", "dnf", "yum", "zypper"}

// dockerCmd represents the docker command
var dockerCmd = &cobra.Command{
	Use:   "docker",
	Short: "Set up the environment in Docker Container",
	Run: func(cmd *cobra.Command, args []string) {
		environment, _ := cmd.Flags().GetString("environment")
		pkgmanager, _ := cmd.Flags().GetString("pkgmanager")
		library, _ := cmd.Flags().GetString("library")
		distro,_ := cmd.Flags().GetString("distro")
		os_ := cmd.Flags().Lookup("os").Changed
		build,_ := cmd.Flags().GetString("build")
		port,_ := cmd.Flags().GetString("port")
		name,_ := cmd.Flags().GetString("name")
		network,_ := cmd.Flags().GetString("network")
		driver,_ := cmd.Flags().GetString("driver")
		subnet,_ := cmd.Flags().GetString("subnet")
		gateway,_ := cmd.Flags().GetString("gateway")
		ip,_ := cmd.Flags().GetString("ip")

		check_docker := exec.Command("docker", "network", "ls")
		var stderr_check_docker bytes.Buffer
		check_docker.Stderr = &stderr_check_docker

		err := check_docker.Run()
		

		if err != nil {
			red.Print("==> [Error] ")
			white.Print("Docker not running\n")
			initLog("docker.go " + string(stderr_check_docker.Bytes()))
		}else{
			if environment == "" && pkgmanager == "" && library == "" && os_ == false && distro == "" && build == ""{
				red.Print("==> [Error] ")
				white.Print("Missing flag\n")
			}else{
				Network(name, network, driver, subnet, gateway, ip)
				if port != ""{
					checkport := CheckPort(port)
					if checkport == true{
						Docker(environment, pkgmanager, library, os_, distro, build, port, name)
					}
				}else{
					Docker(environment, pkgmanager, library, os_, distro, build, port, name)
				}
			}		
		}
	},
}

func init() {
	runCmd.AddCommand(dockerCmd)

	dockerCmd.PersistentFlags().String("environment", "", "Base of the Environment")
	dockerCmd.PersistentFlags().String("pkgmanager", "", "Package manager")
	dockerCmd.PersistentFlags().String("distro", "", "Linux distro")
	dockerCmd.PersistentFlags().String("library", "", "Library to pre-install (separate by comma)")
 	dockerCmd.PersistentFlags().BoolP("os", "", true, "Supported OS for Docker environment and their package manager")
	dockerCmd.PersistentFlags().String("build", "", "Build docker machine using docker file")
	dockerCmd.PersistentFlags().String("port", "", "Expose port")
	dockerCmd.PersistentFlags().String("name", "", "Name of environment (Name of current directory by default)")
	dockerCmd.PersistentFlags().String("network", "", "Network name (Default: name of docker image)")
	dockerCmd.PersistentFlags().String("driver", "", "Driver (Default: bridge) (bridge, host, none, overlay, ipvlan, macvlan)")
	dockerCmd.PersistentFlags().String("subnet", "", "Subnet")
	dockerCmd.PersistentFlags().String("gateway", "", "Gateway")
	dockerCmd.PersistentFlags().String("ip", "", "IP address")
}

func Network(name string, network string, driver string, subnet string, gateway string, ip string){
	dockernetwork_slice := []Flag{}
	dockernetwork_setup := []string{}
	dockernetwork_name := []string{"network", "driver", "subnet", "gateway", "ip"}
	if network == ""{
		network = name
	}
	dockernetwork_flag := []string{network, driver, subnet, gateway, ip}
	for i,s := range dockernetwork_flag{
		if (s != ""){
			temp Flag;
			temp.flag = dockernetwork_name[i]
			temp.name = s
			temp.iter = i
			dockernetwork_slice = append(dockernetwork_slice, temp)
		}
	}
	for _,ds := range dockernetwork_slice{
		if ds.iter != 0{
			dockernetwork_setup = append(dockernetwork_setup, fmt.Sprintf("--%v=%v", ds.flag, ds.name))
		}else{
			n := ds.name
		}
	}
	dockernetwork_setup = append(dockernetwork_setup, n)
	dockernetwork := exec.Command("docker", "network", "create", dockernetwork_setup[0], dockernetwork_setup[1:]...)

	var stderr_dockernetwork bytes.Buffer
	dockernetwork.Stderr = &stderr_dockernetwork

	err := dockernetwork.Run()

	if err != nil{
		red.Print("==> [Error] ")
		white.Print("Problem when setting up network\n")	
		initLog("docker.go " + fmt.Sprintf("%v",err))
	}
	
}

func CheckPort(port string) bool {
	port_ := strings.Split(port, ",")
	for _, p := range port_{ 
		port_int, err := strconv.Atoi(p)
		if err != nil{
			red.Print("==> [Error] ")
			white.Print("Port must be a positive whole number\n")	
			initLog("docker.go " + fmt.Sprintf("%v",err))
			break	
			return false
		}else{
			if port_int > 65535{
				red.Print("==> [Error] ")
				white.Print("Port out of range\n")	
				break	
				return false
			} 
		}
	}
	return true
}

func Docker(environment string, pkgmanager string, library string, os_ bool, distro string, build string, port string, name string){
	if os_ == true {
		if environment == "" && pkgmanager == "" && distro == "" && build == ""{
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
			`)
		}else{
			red.Print("==> [Error] ")
			white.Print("Invalid combination of flag\n")
		}
	}else if build != ""{
		if environment == "" && pkgmanager == "" && distro == "" && os_ == false{
			dir_name := ""
			if name == ""{	
				dir_name = GetDirectoryName()
			}else{
				dir_name = name
			}
			container_id := GetContainerID(dir_name)
			if container_id == "Not Found"{
				dir_mount := fmt.Sprintf(".:/%v", dir_name)
				workdir := fmt.Sprintf("/%v", dir_name)

				blue.Print("==> [In Progress] ")
				white.Print("Setting up environment...\n")
				dockerbuild := exec.Command("docker", "build", build, "-t", string(dir_name), "--quite")
				var stderr_dockerbuild bytes.Buffer
				dockerbuild.Stderr = &stderr_dockerbuild

				blue.Print("==> [In Progress] ")
				white.Print("Starting the environment...\n")				

				err_build := dockerbuild.Run()

				if err_build == nil{
					data := Port_to_Docker(dir_name, dir_mount, workdir, port)
					dockerrun := exec.Command(data[0], data[1:]...)
					var stderr_dockerrun bytes.Buffer
					dockerrun.Stderr = &stderr_dockerrun

					save_name := fmt.Sprintf("./backup/%v-original.tar", dir_name)
					dockersave := exec.Command("docker", "image", "save", "-o", save_name, dir_name)
					var stderr_dockersave bytes.Buffer
					dockersave.Stderr = &stderr_dockersave
	
					_ = dockersave.Run()
					_ = dockerrun.Run()

					initLog("docker.go " + string(stderr_dockersave.Bytes()))
					initLog("docker.go " + string(stderr_dockerrun.Bytes()))

					container_id = GetContainerID(dir_name)

					green.Print("==> [Success] ")
					white.Print(fmt.Sprintf("The environment is running as a docker container named %v with id %v\n", dir_name, container_id))

					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					err_exec := dockerexec.Run()

					initLog("docker.go " + fmt.Sprintf("%v",err_exec))
				}else{
					red.Print("==> [Error] ")
					white.Print(err_build)	
					initLog("docker.go " + string(stderr_dockerbuild.Bytes()))
				}
			}else{
				red.Print("==> [Error] ")
				white.Print("The environment already up\n")		
				
				dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
				dockerexec.Stdin = os.Stdin
				dockerexec.Stdout = os.Stdout
				dockerexec.Stderr = os.Stderr
				
				err_exec := dockerexec.Run()
				initLog("docker.go " + fmt.Sprintf("%v", err_exec))
			}
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
				dir_name := ""
				if name == ""{	
					dir_name = GetDirectoryName()
				}else{
					dir_name = name
				}
				container_id := GetContainerID(dir_name)
				if container_id == "Not Found"{
					dir_mount := fmt.Sprintf(".:/%v", dir_name)
					workdir := fmt.Sprintf("/%v", dir_name)		

					dockerbuild := exec.Command("docker", "build", ".", "-t", string(dir_name), "--quite")
					var stderr_dockerbuild bytes.Buffer
					dockerbuild.Stderr = &stderr_dockerbuild

					dockerfile_rm := exec.Command("rm", "Dockerfile")		
					var stderr_dockerfile_rm bytes.Buffer
					dockerfile_rm.Stderr = &stderr_dockerfile_rm
					
					data := Port_to_Docker(dir_name, dir_mount, workdir, port)
					dockerrun := exec.Command(data[0], data[1:]...)
					var stderr_dockerrun bytes.Buffer
					dockerrun.Stderr = &stderr_dockerrun

					blue.Print("==> [In Progress] ")
					white.Print("Setting up environment...\n")
					_ = dockerbuild.Run()

					initLog("docker.go " + string(stderr_dockerbuild.Bytes()))

					save_name := fmt.Sprintf(fmt.Sprintf("./backup/%v-original.tar", dir_name))
					dockersave := exec.Command("docker", "image", "save", "-o", save_name, dir_name)
					var stderr_dockersave bytes.Buffer
					dockersave.Stderr = &stderr_dockersave

					_ = dockersave.Run()

					blue.Print("==> [In Progress] ")
					white.Print("Starting the environment...\n")
					_ = dockerrun.Run()

					initLog("docker.go " + string(stderr_dockersave.Bytes()))
					initLog("docker.go " + string(stderr_dockerrun.Bytes()))

					container_id = GetContainerID(dir_name)

					green.Print("==> [Success] ")
					white.Print(fmt.Sprintf("The environment is running as a docker container named %v with id %v\n", dir_name, container_id))

					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					err_exec := dockerexec.Run()
					_ = dockerfile_rm.Run()

					initLog("docker.go " + fmt.Sprintf("%v", err_exec))
					initLog("docker.go " + string(stderr_dockerfile_rm.Bytes()))
				}else{
					red.Print("==> [Error] ")
					white.Print("The environment already up\n")	
					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					err_exec := dockerexec.Run()
					initLog("docker.go " + fmt.Sprintf("%v", err_exec))
				}

			}		
		}else if strings.ToLower(environment) == "windows"{
			check_flag := CheckFlag(pkgmanager, distro, library, "Windows")
			if check_flag == true { 
				dir_name := ""
				if name == ""{	
					dir_name = GetDirectoryName()
				}else{
					dir_name = name
				}
				container_id := GetContainerID(dir_name)
				if container_id == "Not Found"{
					dir_mount := fmt.Sprintf(".:/%v", dir_name)
					workdir := fmt.Sprintf("/%v", dir_name)

					dockerbuild := exec.Command("docker", "build", ".", "-t", string(dir_name))
					var stderr_dockerbuild bytes.Buffer
					dockerbuild.Stderr = &stderr_dockerbuild

					dockerfile_rm := exec.Command("rm", "Dockerfile")				
					var stderr_dockerfile_rm bytes.Buffer
					dockerfile_rm.Stderr = &stderr_dockerfile_rm

					data := Port_to_Docker(dir_name, dir_mount, workdir, port)
					dockerrun := exec.Command(data[0], data[1:]...)
					var stderr_dockerrun bytes.Buffer
					dockerrun.Stderr = &stderr_dockerrun

					blue.Print("==> [In Progress] ")
					white.Print("Setting up environment...\n")
					_ = dockerbuild.Run()
					initLog("docker.go " + string(stderr_dockerbuild.Bytes()))
					save_name := fmt.Sprintf(fmt.Sprintf("./backup/%v-original.tar", dir_name))
					dockersave := exec.Command("docker", "image", "save", "-o", save_name, dir_name)
					var stderr_dockersave bytes.Buffer
					dockersave.Stderr = &stderr_dockersave

					_ = dockersave.Run()
					initLog("docker.go " + string(stderr_dockersave.Bytes()))

					blue.Print("==> [In Progress] ")
					white.Print("Starting the environment...\n")
					_ = dockerrun.Run()

					initLog("docker.go " + string(stderr_dockerrun.Bytes()))

					container_id = GetContainerID(dir_name)

					green.Print("==> [Success] ")
					white.Print(fmt.Sprintf("The environment is running as a docker container named %v with id %v\n", dir_name, container_id))

					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					err_exec := dockerexec.Run()
					_ = dockerfile_rm.Run()

					initLog("docker.go " + fmt.Sprintf("%v", err_exec))
					initLog("docker.go " + string(stderr_dockerfile_rm.Bytes()))
				}else{
					red.Print("==> [Error] ")
					white.Print("The environment already up\n")		
					dockerexec := exec.Command("docker", "exec", "-it" , container_id, "/bin/bash")
					dockerexec.Stdin = os.Stdin
					dockerexec.Stdout = os.Stdout
					dockerexec.Stderr = os.Stderr
				
					err_exec := dockerexec.Run()

					initLog("docker.go " + fmt.Sprintf("%v", err_exec))

				}
			}
		}
	}
}

func Port_to_Docker(dir_name string, dir_mount string, workdir string, port string) []string {
	if port == ""{
		mkdir_run := exec.Command("mkdir", "./run")

		err_mkdir := mkdir_run.Run()
		
		initLog("docker.go " + fmt.Sprintf("%v", err_mkdir))

		reset := []byte("")
		err_writefile := os.WriteFile("./run/connect", reset, 0644)

		initLog("docker.go " + fmt.Sprintf("%v", err_writefile))

		file, err_openfile := os.OpenFile("./run/connect", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

		initLog("docker.go " + fmt.Sprintf("%v", err_openfile))

		write := fmt.Sprintf("docker run --name %v -v %v --workdir %v -itd %v bash", dir_name, dir_mount, workdir, dir_name)

		_,_ = file.WriteString(write)

		dat, err_readfile := os.ReadFile("./run/connect")
		data := strings.Split(string(dat), " ")

		initLog("docker.go " + fmt.Sprintf("%v", err_readfile))

		return data
	}else{
		p := strings.Split(port, ",")
		reset := []byte("")
		err_writefile := os.WriteFile("./run/connect", reset, 0644)
		initLog("docker.go " + fmt.Sprintf("%v", err_writefile))

		file, err_openfile := os.OpenFile("./run/connect", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		initLog("docker.go " + fmt.Sprintf("%v", err_openfile))
		write := fmt.Sprintf("docker run --name %v -v %v --workdir %v",dir_name, dir_mount, workdir)

		for _, p_ := range p{
			write += fmt.Sprintf(" -p %v:%v", p_, p_)
		}

		write += fmt.Sprintf(" -itd %v bash", dir_name)	
		_,_ = file.WriteString(write)

		dat, err_readfile := os.ReadFile("./run/connect")

		initLog("docker.go " + fmt.Sprintf("%v", err_readfile))

		data := strings.Split(string(dat), " ")

		return data
	}	

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
	err_writefile := os.WriteFile("./Dockerfile", reset, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_writefile))
	//Open Docker file
	file, err_openfile := os.OpenFile("./Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	 
	initLog("docker.go " + fmt.Sprintf("%v", err_openfile))

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
			writelibrary := fmt.Sprintf("RUN apt install -y %v 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

	writelibrary_essential := "\nRUN apt install git -y 2>/dev/null"
	_,_ = file.WriteString(writelibrary_essential)

	write_ttyd := "\nRUN apt install -y build-essential cmake git libjson-c-dev libwebsockets-dev\nRUN git clone https://github.com/tsl0922/ttyd.git\nWORKDIR ttyd\nRUN mkdir build\nWORKDIR build\nRUN cmake ..\nRUN make && make install"

	_,_ = file.WriteString(write_ttyd)

}

func dnf(distro string, target string, libraries string){
	reset := []byte("")
	err_writefile := os.WriteFile("./Dockerfile", reset, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_writefile))
	//Open Docker file
	file, err_openfile := os.OpenFile("./Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_openfile))

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
			writelibrary := fmt.Sprintf("RUN dnf install %v -y 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}
	writelibrary_essential := "\nRUN dnf install git -y"
	_,_ = file.WriteString(writelibrary_essential)


}

func yum(distro string, target string, libraries string){
	reset := []byte("")
	err_writefile := os.WriteFile("./Dockerfile", reset, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_writefile))
	//Open Docker file
	file, err_openfile := os.OpenFile("./Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_openfile))

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
			writelibrary := fmt.Sprintf("RUN yum install %v -y 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

	writelibrary_essential := "\nRUN yum install git -y"
	_,_ = file.WriteString(writelibrary_essential)

}

func zypper(distro string, target string, libraries string){
	reset := []byte("")
	err_writefile := os.WriteFile("./Dockerfile", reset, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_writefile))
	//Open Docker file
	file, err_openfile := os.OpenFile("Dockerfile", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	initLog("docker.go " + fmt.Sprintf("%v", err_openfile))
	 
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
			writelibrary := fmt.Sprintf("RUN zypper install %v -y 2>/dev/null\n", o)
			_,_ = file.WriteString(writelibrary)
		}
	}

	writelibrary_essential := "\nRUN zypper install git -y"
	_,_ = file.WriteString(writelibrary_essential)


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
