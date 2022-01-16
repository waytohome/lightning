package confx

import (
	"strings"
	"testing"
)

func TestReadFileNotExist(t *testing.T) {
	_, err := NewFileConfigure("abc.json", nil)
	if err != nil && strings.Contains(err.Error(), "no such file or directory") {
		t.Logf("文件不存在")
	} else {
		t.Fatal(err)
	}
}

func TestReadYaml(t *testing.T) {
	configure, err := NewFileConfigure("test.yaml", nil)
	if err != nil {
		t.Fatal(err)
	}
	port, _ := configure.GetString("server.port", ":8080")
	timeout, _ := configure.GetInt("server.timeout", 30)
	exist := configure.Exist("server.upload-size")
	t.Logf("port = %v", port)
	t.Logf("timeout = %v", timeout)
	t.Logf("exist = %v", exist)
}
