package latency

import (
	"fmt"
	"net"
	"testing"
	"time"
)

func TestXxx(t *testing.T) {
	// 可以直接用域名
	dst := "duckduckgo.com:80" // HTTP 端口

	start := time.Now()
	conn, err := net.DialTimeout("tcp", dst, 2*time.Second)
	if err != nil {
		fmt.Println("连接失败:", err)
		return
	}
	defer conn.Close()

	elapsed := time.Since(start)
	fmt.Printf("连接到 %s 延迟: %v\n", dst, elapsed)
}
