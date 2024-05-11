package main

import (
	"log"
	"os"
	"fmt"
	"os/exec"
)

func main(){
	mkdir_run := exec.Command("mkdir", "../log")

	_ = mkdir_run.Run()
	err1 := os.Args[1]

	logFile, err_log := os.OpenFile("../log/app.log", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	log.SetOutput(logFile)
	log.SetFlags(log.Lshortfile | log.LstdFlags | log.LUTC) 

	log.Println("UTC: " + fmt.Sprintf("%v",err1))  

	if err_log != nil{
		log.Println("UTC: " + fmt.Sprintf("%v",err_log))
	}
}
