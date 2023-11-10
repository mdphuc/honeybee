package main

import 	"github.com/fatih/color"

func main(){
	err := color.New(color.FgWhite, color.Bold)

	err.Print("Hello")
	err.Print("yo")
}