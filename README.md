# biever v1.0

<img src="https://github.com/mdphuc/beaver/assets/41264640/bd857d8c-104b-4c95-9c70-bdb6531c406e" style="width:250px; height:auto">

#
The project, utilizing Golang, is designed as a middle place to run and test code without having to worried about unmatched environemnt, missing library, or cannot install library or reach out to certain sources because of firewall or proxy problems. It supports docker and a remote machine

## Table of Content
1) [Installation](#installation)
2) [Usage](#usage)
3) [Requirement](#requirement)
4) [User Guide](#user-guide)
5) [Contributing Policy](#contributing)
6) [License](#license)

## Installation
- ```git clone https://github.com/mdphuc/beaver; cd ./beaver; go get -u -v -f all; go build; go install```

  _or_
- If you have recent go compiler install: ```go install github.com/mdphuc/beaver```

  _or_
- Download prebuilt file: https://github.com/mdphuc/beaver/releases/download/v1.0/beaver

## Usage
- beaver 
```
beaver v1.0
    
Set up remote development environment in isolated
enivronment like docker or proxy server to use
remote machine or cloud machine as development environment

Usage:
  beaver [command]

Available Commands:
  help        Help about any command
  reset       Reset the environment
  run         Run Set up the environment

Flags:
  -h, --help      help for biever
  -v, --version   version for biever

Use "beaver [command] --help" for more information about a command.
```
- beaver run
```
Run Set up the environment

Usage:
  beaver run [command]

Available Commands:
  docker        Set up the environment in Docker Container
  remoteMachine Set up proxy server to use remote machine as development environment

Flags:
  -h, --help   help for run

Use "beaver run [command] --help" for more information about a command.
```
- beaver run docker
```
Set up the environment in Docker Container

Usage:
  beaver run docker [flags]

Flags:
      --build string         Build docker machine using docker file
      --distro string        Linux distro
      --environment string   Base of the Environment
  -h, --help                 help for docker
      --library string       Library to pre-install (separate by comma)
      --os                   Supported OS for Docker environment and their package manager (default true)
      --pkgmanager string    Package manager
```
- beaver run remoteMachine
```
Set up proxy server to use remote machine as development environment

Usage:
  beaver run remoteMachine [flags]

Flags:
      --compose       Build proxy server (default true)
      --connect       Set up connection between proxy server and remote machine (default true)
  -h, --help          help for remoteMachine
      --ip string     IP address of the remote environment
      --user string   Username for remote machine
```
- beaver reset
```
Reset the environment

Usage:
  beaver reset [flags]

Flags:
      --docker           Username for remote machine (default true)
  -h, --help             help for reset
      --remote_machine   Build proxy server (default true)
```

## Requirement
- Go v1.21+

## User Guide
- When run `beaver run remoteMachine`:
```
1) beaver run remoteMachine --compose --ip=<ip of remote machine> --user=<username of remote machine>
2) beaver run remoteMachine --connect
```

## License
[Apache-2.0](https://choosealicense.com/licenses/apache-2.0/)

