package main

import(
	"os/exec"
	"os"
	// "fmt"


)

func main(){
	// a := ".23.45."
	// b := strings.Split(a, ".")
	// fmt.Println(len(b))
	// for i, _ := range b{
	// 	fmt.Println(i)
	// }

	// c := "0"
	// d, _ := strconv.Atoi(c)
	// fmt.Println(d + 9)
	// yo := CheckValidAllowAccess(" 223.56.56.7")
	// fmt.Println(yo)

	// notexist1 := true
	// ok1 := true

	// notexist2 := true
	// ok2 := true
	// if (notexist1 && ok1) && (notexist2 && ok2){
	// 	fmt.Println("hello")
	// }
	webui := exec.Command("gunicorn", "--chdir", "./web", "-w", "3", "wsgi:app")
	webui.Stderr = os.Stderr
	webui.Stdout = os.Stdout
	webui.Stdin = os.Stdin

	_ = webui.Run()
}


