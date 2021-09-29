package main

import "seasonjs/espack/internal/hooks"

// PrintRunHelp espack run nil 或者espack run 不识别的标志
func PrintRunHelp() {

}

// RunDev espack run dev
func RunDev() {
	hooks.
		NewHookContext().
		InitHooks().
		InstallPlugin().
		StartDevServer().
		StartESBuild().
		HoldAll()
}

// RunESPacK espack run ...
func RunESPacK(subArgs []string) {
	if len(subArgs) < 1 {
		PrintRunHelp()
		return
	}
	arg := subArgs[0]
	switch {
	case arg == "dev":
		RunDev()
	case arg == "build":
		RunDev()
	default:
		PrintRunHelp()
	}

}
