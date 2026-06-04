package main

import (
	"fmt"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[0])
	if err != nil {
		fmt.Printf("error during read dir: %v", err)
		os.Exit(103)
		return
	}

	if len(os.Args) < 4 {
		fmt.Printf("uncorrect count args\n")
		os.Exit(104)
		return
	}

	rslt := RunCmd(os.Args[1:], env)
	os.Exit(rslt)
}
