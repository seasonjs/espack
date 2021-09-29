// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in https://github.com/gin-gonic/gin/blob/master/LICENSE

package liteS

import (
	"io"
	"os"
	"runtime"
	"seasonjs/espack/internal/logger"
	"seasonjs/espack/internal/utils"
	"strconv"
	"strings"
)

// DefaultWriter is the default io.Writer used by Gin for debug output and
// middleware output like Logger() or Recovery().
// Note that both Logger and Recovery provides custom ways to configure their
// output io.Writer.
// To support coloring in Windows use:
// 		import "github.com/mattn/go-colorable"
// 		gin.DefaultWriter = colorable.NewColorableStdout()
var DefaultWriter io.Writer = os.Stdout

// DefaultErrorWriter is the default io.Writer used by Gin to debug errors
var DefaultErrorWriter io.Writer = os.Stderr

const ginSupportMinGoVer = 13

// IsDebugging returns true if the framework is running in debug mode.
// Use SetMode(gin.ReleaseMode) to disable debug mode.
func IsDebugging() bool {
	//return ginMode == debugCode
	return utils.Env.Dev()
}

// DebugPrintRouteFunc indicates debug log output format.
var DebugPrintRouteFunc func(httpMethod, absolutePath, handlerName string, nuHandlers int)

//func debugPrintRoute(httpMethod, absolutePath string, handlers HandlersChain) {
//	if IsDebugging() {
//		nuHandlers := len(handlers)
//		handlerName := nameOfFunction(handlers.Last())
//		if DebugPrintRouteFunc == nil {
//			debugPrint("%-6s %-25s --> %s (%d handlers)\n", httpMethod, absolutePath, handlerName, nuHandlers)
//		} else {
//			DebugPrintRouteFunc(httpMethod, absolutePath, handlerName, nuHandlers)
//		}
//	}
//}
//
//func debugPrintLoadTemplate(tmpl *template.Template) {
//	if IsDebugging() {
//		var buf strings.Builder
//		for _, tmpl := range tmpl.Templates() {
//			buf.WriteString("\t- ")
//			buf.WriteString(tmpl.Name())
//			buf.WriteString("\n")
//		}
//		debugPrint("Loaded HTML Templates (%d): \n%s\n", len(tmpl.Templates()), buf.String())
//	}
//}

func debugPrint(format string, values ...interface{}) {
	if IsDebugging() {
		if strings.HasSuffix(format, "\n") {
			strings.TrimRight(format, "\n")
		}
		logger.Info(format, values...)
	}
}

func getMinVer(v string) (uint64, error) {
	first := strings.IndexByte(v, '.')
	last := strings.LastIndexByte(v, '.')
	if first == last {
		return strconv.ParseUint(v[first+1:], 10, 64)
	}
	return strconv.ParseUint(v[first+1:last], 10, 64)
}

func debugPrintWARNINGDefault() {
	if v, e := getMinVer(runtime.Version()); e == nil && v <= ginSupportMinGoVer {
		debugPrint(`liteS 需要 Go 1.13及以上.

`)
	}
	debugPrint(`liteS Engine 实例化完成.`)
}

//
//func debugPrintWARNINGNew() {
//	debugPrint(`[WARNING] Running in "debug" mode. Switch to "release" mode in production.
//- using env:	export GIN_MODE=release
//- using code:	gin.SetMode(gin.ReleaseMode)
//
//`)
//}
//
//func debugPrintWARNINGSetHTMLTemplate() {
//	debugPrint(`[WARNING] Since SetHTMLTemplate() is NOT thread-safe. It should only be called
//at initialization. ie. before any route is registered or the router is listening in a socket:
//
//	router := gin.Default()
//	router.SetHTMLTemplate(template) // << good place
//
//`)
//}

func debugPrintError(err error) {
	if err != nil && IsDebugging() {
		logger.Warn("[liteS] 告警 %v\n", err)
	}
}
