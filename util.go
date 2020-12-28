package toolbox

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"

	"toolbox/log"
)

var LocalIp string

func init() {
	LocalIp = getLocalIP()
}

func LocalIP() string {
	return LocalIp
}

func GetLocalIP() string {
	return LocalIp
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// /helloworld.Greeter/SayHello
func StatProtoMetrix(method string) string {
	if LocalIp == "" {
		LocalIp = "unknown"
	}
	return fmt.Sprintf("%s.%s", strings.Replace(method[1:], "/", ".", -1), LocalIp)
}

func StatMetrix(name string) string {

	if LocalIp == "" {
		LocalIp = "unknown"
	}
	return fmt.Sprintf("%s.%s", name, LocalIp)
}

//捕捉panic
func CatchPanic(ctx context.Context, in interface{}) {
	if err := recover(); err != nil {
		var buf [4096]byte
		n := runtime.Stack(buf[:], false)
		log.Errorf(ctx, fmt.Sprintf("panic req:%+v err:%+v", in, err))
		log.Errorf(ctx, "%s", string(buf[:n]))
	}
}

func MkdirIfNotExist(path string) error {
	var err error

	// 存在且err!=nil的情况不考虑
	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	return err
}

func CreateFileIfNotExist(name string, truncate bool) (*os.File, bool, error) {
	var (
		err  error
		file *os.File

		fileExist bool = true
	)

	_, err = os.Stat(name)

	if os.IsNotExist(err) {
		fileExist = false
	}

	if !fileExist {
		file, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
		if err != nil {
			panic(err.Error())
		}
	} else {
		if truncate {
			file, err = os.OpenFile(name, os.O_TRUNC|os.O_WRONLY, 0644)
			if err != nil {
				panic(err.Error())
			}
		}
	}

	return file, fileExist, err
}
