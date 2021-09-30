package main

import "seasonjs/espack/internal/mod"

func EsPackGet(subArgs []string) {
	if len(subArgs) == 0 {
		mod.NewMod().AnalyzeDependencies().DownLoadDependencies()
	}
}
