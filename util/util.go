package util

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
)

func GetLocalAddressOld() string {
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

func GetLocalAddress() string {
	// get available network interfaces for
	// this machine
	interfaces, err := net.Interfaces()

	if err != nil {
		return GetLocalAddressOld()
	}

	ips := make(map[string]string)

	for _, i := range interfaces {
		byNameInterface, err := net.InterfaceByName(i.Name)
		if err != nil {
			return GetLocalAddressOld()
		}

		addresses, err := byNameInterface.Addrs()
		if len(addresses) > 0 {
			ips[i.Name] = strings.Split(addresses[0].String(), "/")[0]
		}
	}
	for _, ifname := range []string{"br1", "bond1", "eth1", "eth0"} {
		if value, ok := ips[ifname]; ok {
			return value
		}
	}

	return GetLocalAddressOld()
}

func FileExists(dir string) bool {
	_, err := os.Stat(dir)
	return err == nil
}

func Mkdir(dir string) error {
	return os.MkdirAll(dir, os.ModePerm)
}

func EnsureDirExists(dir string) (err error) {
	f, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			//如果不存在，创建
			return os.MkdirAll(dir, os.FileMode(0755))
		} else {
			return err
		}
	}

	if !f.IsDir() {
		//已存在，但是不是文件夹
		return fmt.Errorf("path %s is exist,but not dir", dir)
	}

	return nil
}

func DecompressFile(compressedFile string, decompressDir string) error {
	// check the file
	if !strings.HasSuffix(compressedFile, "tar.gz") {
		return fmt.Errorf("%s is not a tar.gz file", compressedFile)
	}

	f, err := os.Stat(compressedFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s is not exist", compressedFile)
		} else {
			return fmt.Errorf("unknow error when get info of %s", compressedFile)
		}
	}

	if f.IsDir() {
		return fmt.Errorf("%s is a directory", compressedFile)
	}

	// ensure dest dir exist
	if err := EnsureDirExists(decompressDir); err != nil {
		return fmt.Errorf("ensure target dir:%s exist failed:%v", decompressDir, err)
	}

	var args []string
	args = append(args, "zxf")
	args = append(args, compressedFile)
	args = append(args, "-C")
	args = append(args, decompressDir)
	cmd := exec.Command("tar", args...)
	/*
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("run tar command, failed:%v, arguments:%v", err, args)
		}
	*/
	_, err = cmd.CombinedOutput()

	return nil
}

func RemoveFiles(paths []string) error {
	for _, p := range paths {
		if p == "" {
			continue
		}

		cmd := exec.Command("rm", "-rf", p)
		if _, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("exec rm -rf %s failed:%v", p, err)
		}
	}

	return nil
}

func IsYamlString(yamlString string) bool {
	_, err := yaml.YAMLToJSON([]byte(yamlString))
	if err != nil {
		return false
	}

	return true
}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 50)
	suffix = strings.ToUpper(suffix)                                                     //忽略后缀匹配的大小写
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
