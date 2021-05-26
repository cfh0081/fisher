package fisher

import (
	"context"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRecordData(t *testing.T) {
	// 临时创建一个文件用于验证文件写入
	file, err := ioutil.TempFile("", "recorder-*.jl")
	assert.Nil(t, err)
	// 确保程序结束时删除临时文件
	defer os.Remove(file.Name())
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	toWrite := "The data to write!"
	expected := []byte(toWrite + "\n")
	SetStorageTarget(file)
	go StartStorageServer(ctx)
	RecordData([]byte(toWrite))
	time.Sleep(3 * time.Second)
	err = file.Close()
	assert.Nil(t, err)

	data, err := ioutil.ReadFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
