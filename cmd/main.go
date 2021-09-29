package main

import (
	"fmt"
	"os"
)

func main() {
	osArgs := os.Args[1:]
	if len(osArgs) < 1 {
		PrintHelp()
		os.Exit(0)
	}
	//子命令
	cmdArg := osArgs[0]
	//子命令参数
	subArgs := osArgs[1:]
	switch {
	case cmdArg == "-h", cmdArg == "-help", cmdArg == "--help", cmdArg == "help":
		PrintHelp()
		os.Exit(0)
	case cmdArg == "run":
		RunESPacK(subArgs)
		os.Exit(0)
	case cmdArg == "--version", cmdArg == "-v", cmdArg == "-version", cmdArg == "version":
		fmt.Printf("espack 版本为: %s\n", espackVersion)
		os.Exit(0)
	case cmdArg == "get", cmdArg == "install", cmdArg == "add", cmdArg == "i":
		EsPackGet(subArgs)
	default:
		PrintHelp()
	}

}
