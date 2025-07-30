package benchmark

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/antlabs/gurl/internal/testserver"
)

// TestgurlCommandLineBasic 测试gurl命令行基本功能
func TestgurlCommandLineBasic(t *testing.T) {
	// 启动测试服务器
	server := testserver.NewTestServer(8090)
	go func() {
		server.Start()
	}()
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(300 * time.Millisecond)

	// 构建gurl二进制文件
	buildCmd := exec.Command("go", "build", "-o", "gurl-test", "./cmd/gurl")
	buildCmd.Dir = "../../"
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("构建gurl失败: %v\n输出: %s", err, string(output))
	}
	defer exec.Command("rm", "../../gurl-test").Run() // 清理

	tests := []struct {
		name     string
		method   string
		endpoint string
		args     []string
	}{
		{
			"GET请求测试",
			"GET",
			"/api/get",
			[]string{"../../gurl-test", "-c", "2", "-d", "2s", "-t", "1"},
		},
		{
			"POST请求测试",
			"POST",
			"/api/post",
			[]string{"../../gurl-test", "-c", "2", "-d", "2s", "-t", "1", "--parse-curl"},
		},
		{
			"PUT请求测试",
			"PUT",
			"/api/put",
			[]string{"../../gurl-test", "-c", "2", "-d", "2s", "-t", "1", "--parse-curl"},
		},
		{
			"DELETE请求测试",
			"DELETE",
			"/api/delete",
			[]string{"../../gurl-test", "-c", "2", "-d", "2s", "-t", "1", "--parse-curl"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var cmd *exec.Cmd
			url := server.GetURL() + test.endpoint

			if test.method == "GET" {
				// 对于GET请求，直接使用URL
				args := append(test.args, url)
				cmd = exec.Command(args[0], args[1:]...)
			} else {
				// 对于其他方法，使用curl命令格式
				curlCmd := fmt.Sprintf(`curl -X %s "%s"`, test.method, url)
				if test.method == "POST" || test.method == "PUT" {
					curlCmd += ` -H "Content-Type: application/json" -d '{"test": "data"}'`
				}
				args := append(test.args, curlCmd)
				cmd = exec.Command(args[0], args[1:]...)
			}

			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			cmd = exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%s失败: %v\n输出: %s", test.name, err, string(output))
			}

			outputStr := string(output)

			// 验证输出包含基本统计信息
			if !strings.Contains(outputStr, "Requests/sec") {
				t.Errorf("%s输出应该包含 'Requests/sec'\n输出: %s", test.name, outputStr)
			}

			if !strings.Contains(outputStr, "Transfer/sec") {
				t.Errorf("%s输出应该包含 'Transfer/sec'\n输出: %s", test.name, outputStr)
			}

			t.Logf("%s成功完成\n输出:\n%s", test.name, outputStr)
		})
	}
}

// TestgurlCommandLineAdvanced 测试gurl命令行高级功能
func TestgurlCommandLineAdvanced(t *testing.T) {
	// 启动测试服务器
	server := testserver.NewTestServer(8091)
	go func() {
		server.Start()
	}()
	defer server.Stop()

	// 等待服务器启动
	time.Sleep(300 * time.Millisecond)

	// 构建gurl二进制文件
	buildCmd := exec.Command("go", "build", "-o", "gurl-test", "./cmd/gurl")
	buildCmd.Dir = "../../"
	output, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("构建gurl失败: %v\n输出: %s", err, string(output))
	}
	defer exec.Command("rm", "../../gurl-test").Run() // 清理

	tests := []struct {
		name string
		args []string
		url  string
	}{
		{
			"延迟端点测试",
			[]string{"../../gurl-test", "-c", "2", "-d", "3s", "-t", "1", "--latency"},
			server.GetURL() + "/api/delay?ms=50",
		},
		{
			"状态码测试",
			[]string{"../../gurl-test", "-c", "2", "-d", "2s", "-t", "1"},
			server.GetURL() + "/api/status/404",
		},
		{
			"高并发测试",
			[]string{"../../gurl-test", "-c", "10", "-d", "2s", "-t", "4"},
			server.GetURL() + "/api/echo",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			args := append(test.args, test.url)
			cmd := exec.Command(args[0], args[1:]...)

			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			cmd = exec.CommandContext(ctx, cmd.Args[0], cmd.Args[1:]...)

			output, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("%s失败: %v\n输出: %s", test.name, err, string(output))
			}

			outputStr := string(output)

			// 验证输出包含基本统计信息
			if !strings.Contains(outputStr, "Requests/sec") {
				t.Errorf("%s输出应该包含 'Requests/sec'\n输出: %s", test.name, outputStr)
			}

			t.Logf("%s成功完成\n输出:\n%s", test.name, outputStr)
		})
	}
}
