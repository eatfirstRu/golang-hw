package main

import (
	"fmt"
	"os"
)

func main() {
	env, err := ReadDir(os.Args[1])
	if err != nil {
		fmt.Printf("error during read dir: %v", err)
		os.Exit(103)
		return
	}

	if len(os.Args) < 5 {
		fmt.Printf("uncorrect count args\n")
		os.Exit(104)
		return
	}

	rslt := RunCmd(os.Args[2:], env)
	os.Exit(rslt)
}
