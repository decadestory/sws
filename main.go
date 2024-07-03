package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"
)

func main() {
	staticDir, _ := os.Getwd()
	fmt.Println(staticDir)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rp := r.Intn(22000) + 10000

	port := ":" + strconv.Itoa(rp)
	ip := GetIp()

	fs := http.FileServer(http.Dir(staticDir))
	http.Handle("/", http.StripPrefix("/", fs))

	OpenBrowser("http://" + ip + port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func GetIp() string {
	// 获取本机所有网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return ""
	}

	// 遍历接口信息，找到非回环接口的 IPv4 地址
	for _, iface := range interfaces {
		// 排除回环接口和无效接口
		if iface.Flags&net.FlagLoopback == 0 && iface.Flags&net.FlagUp != 0 {
			// 获取接口的 IP 地址
			addrs, err := iface.Addrs()
			if err != nil {
				fmt.Println("Failed to get IP addresses:", err)
				continue
			}

			// 遍历 IP 地址，找到 IPv4 地址
			for _, addr := range addrs {
				ip, ok := addr.(*net.IPNet)
				if ok && !ip.IP.IsLoopback() && ip.IP.To4() != nil {
					return ip.IP.String()
				}
			}
		}
	}

	return ""
}
