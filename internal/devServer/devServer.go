package devServer

import (
	"github.com/evanw/esbuild/pkg/api"
	"mime"
	"net"
	"net/http"
	"path/filepath"
	"seasonjs/espack/internal/devServer/pkg/liteS"
	"seasonjs/espack/internal/utils"
	"strings"
)

//TODO： Server-sent events 解决HMR 问题
//不能使用esbuild的serve,因为它的serve不会提供额外的插件输出能力

type INMemory uint8

const (
	IsINMemoryTrue INMemory = iota
	IsINMemoryFalse
)

type ctx struct {
	iNMemory INMemory
	r        *liteS.Engine
	res      *map[string][]byte
}

func NewDevServer() *ctx {
	//if !utils.Env.Dev() {
	//	liteS.SetMode(liteS.ReleaseMode)
	//}
	// 默认内存读取 TODO:支持从文件夹读取
	return &ctx{IsINMemoryTrue, liteS.Default(), nil}
}

// Add 将build好的资源转换为Map格式
func (c *ctx) Add(outputFiles *[]api.OutputFile) *ctx {
	res := make(map[string][]byte)
	if len(*outputFiles) <= 0 {
		c.res = nil
	}
	for _, file := range *outputFiles {
		path, _ := utils.FS.ConvertPath(file.Path)
		res[path] = file.Contents
	}
	c.res = &res
	return c
}

// Run liteS服务器启动
func (c *ctx) Run() {
	r := c.r
	res := make(map[string][]byte)
	go func() {
		//TODO:proxy,websocket,使用内置特性
		r.GET("/*action", func(g *liteS.Context) {
			p := strings.ToLower(g.Request.URL.Path)
			//转换为映射的路径key TODO: 替换成从配置中读取
			p, _ = utils.FS.ConvertPath("./dist/" + p)
			// 是否是文件夹或者根路径
			if p == "" || utils.FS.IsDir(p) {
				p = filepath.Join(p, "/index.html")
			}
			ct := g.ContentType()
			if len(ct) <= 0 {
				ct = mime.TypeByExtension(filepath.Ext(p))
				g.Header("Content-Type", ct)
			}
			if c.res != nil {
				res = *c.res
			}
			if len(res) <= 0 {
				g.String(http.StatusOK, "构建进行中...")
				return
			}

			body, ok := res[p]
			if ok {
				g.String(http.StatusOK, string(body))
				return
			}

			g.String(http.StatusNotFound, "资源未找到...")
		})
		ln, _ := net.Listen("tcp", "127.0.0.1:8080")
		r.RunListener(ln)
	}()
}
