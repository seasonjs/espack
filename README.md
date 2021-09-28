# espack

纯净的web构建工具

native web bundle tools

> 目的是在无node,无webpack,无npm环境下进行前端构建
>
> aim to work with no node no webpack no npm env build faster

# env

```env
set ES_PACK_DEV_ENV=dev
```

## 初衷

无运行时，无需配置(低配置)，轻量快速高效，最佳实践 ... 大概 ：）

## TODO

[√] Toy webpack

[v] react 成功运行 ~~

[_] 参照 [coredns 的插件逻辑](https://coredns.io/2016/12/19/writing-plugins-for-coredns) 实现插件机制

伪代码逻辑：

```go
package xxxPlugin

import (
	"fmt"
	"github.com/seasonjs/espack/pkg/api/plugin"
)

func init() { plugin.Register("espack_plugin_corejs", NewPlugin) }
func NewPlugin(opt interface{}) *interface{} {
	fmt.Println("plugin is regist in espack")
	return nil
}

```

```json

{
  "entry": {
    "main": "index.jsx"
  },
  "output": "dist",
  "plugins": [
    {
      "name": "espack_plugin_corejs",
      "option": {
        "xxx": "xxx",
        "others": "this is others option"
      }
    }
  ]
}
```

[_] 逐步解决代码中的TODO注释

[_] 调研npm代替方案 js.mod -> go mod like?

```
|
|__ [√] 调研npm元数据api获取
|
|__ [√] 调研unpack方式获取，不大可行，es module 太少（可能会成未来的答案,但是现在我选npm）
|
|__ [√] 调研go解压tgz格式文件
|
|__ [√] 调研go.mod 解析生成原理
|
|__ [_] 具体实施 大概
    |__ [√]获取元数据与
    |__ [√]解决循环依赖
    |__ [√]根据js.mod依赖生成js.sum
    |__ [_]根据package.json 生成js.mod
    |__ [_]根据解析好的数据下载tarball
```

[_] 调研esbuild ast语法析出方案

[_] 补充测试，增加代码覆盖率

[_] 命名更符合go风格（x）,更加c-like(√)

## 不支持

目前来看不支持sass，grpc等一切需要使用node原生环境的包 ，除非你使用了node或者npm ：）