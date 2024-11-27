package main

import (
	"fmt"
	"os"

	"github.com/kpawlik/compare_msg"
)

func main(){
	err := compare_msg.Execute()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
	return
	
}