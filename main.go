package main

import "seasonjs/espack/internal/hooks"

func main() {
	hooks.
		NewHookContext().
		InitHooks().
		InstallPlugin().
		StartDevServer().
		StartESBuild().
		HoldAll()
}
