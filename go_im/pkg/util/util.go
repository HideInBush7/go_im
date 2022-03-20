package util

import (
	"net"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"unsafe"
)

// 获取本机内网IP
func GetInternalIp() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	for _, addr := range addrs {
		// 非环回地址
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ``
}

// 获取当前程序运行的文件夹
// 因为使用相对路径的话在不同目录下调用可执行文件会找不到
func GetExecDir() string {
	dir, err := filepath.Abs(path.Dir(os.Args[0]))
	if err != nil {
		panic(err)
	}
	// windows下分隔符不一致
	return strings.ReplaceAll(dir, "\\", "/")
}

// 通过结构体返回所有字段tag中的db, 加参数without->去除without中的字符串
func GetDbFieldNames(in interface{}) (result []string) {
	v := reflect.ValueOf(in)

	// 指针需要转换
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		panic(`receive 'struct only, but '` + v.Kind().String())
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		tag := typ.Field(i).Tag.Get(`db`)
		result = append(result, tag)
	}
	return
}

func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			cap int
		}{s, len(s)},
	))
}

func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
