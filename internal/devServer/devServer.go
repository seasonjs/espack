package devServer

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/gin-gonic/gin"
	"mime"
	"net/http"
	"path/filepath"
	"seasonjs/espack/internal/utils"
	"strings"
)

type INMemory uint8

const (
	IsINMemoryTrue INMemory = iota
	IsINMemoryFalse
)

type ctx struct {
	iNMemory INMemory
	r        *gin.Engine
	res      *map[string][]byte
}

func NewDevServer() *ctx {
	if !utils.Env.Dev() {
		gin.SetMode(gin.ReleaseMode)
	}
	// 默认内存读取 TODO:支持从文件夹读取
	return &ctx{IsINMemoryTrue, gin.Default(), nil}
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

// Run gin开发服务启动,TODO 定制log
func (c *ctx) Run() {
	r := c.r
	res := make(map[string][]byte)
	go func() {
		//TODO:proxy,websocket
		r.GET("/*action", func(g *gin.Context) {
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
		r.Run()
	}()
}
