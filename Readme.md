# Go-Mirai

[![Go Report](https://goreportcard.com/badge/github.com/Logiase/gomirai?style=flat-square)](https://goreportcard.com/report/github.com/Logiase/gomirai)![GitHub top language](https://img.shields.io/github/languages/top/Logiase/gomirai)![GitHub](https://img.shields.io/github/license/Logiase/gomirai)![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Logiase/gomirai)![GitHub contributors](https://img.shields.io/github/contributors/Logiase/gomirai)

适配[MiraiHttpApi](https://github.com/project-mirai/mirai-api-http)的Go SDK

目前仍处于开发阶段,部分功能仍未测试

## 如何使用

请参照[example](/test/main.go)

所有导出函数、变量、结构均有详细注释

```go

package main

import (
    "github.com/Logiase/gomirai"
)

func main() {
    //...
}
    
```

### 不了解Go？

(安利狂魔) [Go急速入门](https://learn.go.dev/)

### 目前问题

未实现Http api所有功能

## 计划

 - websocket
 - 上报
 - 心跳

## 维护者

[Logiase](https://github.com/Logiase)

## 许可证

[AGPL-3.0](LICENSE) © Logiase
