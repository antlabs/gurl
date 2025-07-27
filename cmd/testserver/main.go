package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/antlabs/murl/internal/testserver"
)

func main() {
	var port = flag.Int("port", 8080, "服务器端口")
	flag.Parse()

	server := testserver.NewTestServer(*port)

	// 设置信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 启动服务器
	go func() {
		log.Printf("启动测试服务器在端口 %d", *port)
		log.Printf("访问 http://localhost:%d/health 进行健康检查", *port)
		log.Printf("可用的端点:")
		log.Printf("  GET    /api/get")
		log.Printf("  POST   /api/post")
		log.Printf("  PUT    /api/put")
		log.Printf("  PATCH  /api/patch")
		log.Printf("  DELETE /api/delete")
		log.Printf("  *      /api/echo")
		log.Printf("  GET    /api/delay?ms=<毫秒>")
		log.Printf("  GET    /api/status/<状态码>")
		log.Printf("  GET    /health")

		if err := server.Start(); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待信号
	<-sigChan
	fmt.Println("\n收到停止信号，正在关闭服务器...")

	if err := server.Stop(); err != nil {
		log.Printf("服务器停止时出错: %v", err)
	} else {
		log.Println("服务器已成功停止")
	}
}
