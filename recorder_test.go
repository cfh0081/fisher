package fisher

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/cfh0081/baseutils"
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
	SetStorageTarget(file)

	go StartStorageServer(ctx)

	toWriteFirst := []byte("The data to write!")
	toWriteSecond := make(map[string]interface{})
	toWriteSecond["a"] = "b"
	toWriteSecond["what"] = 32
	toWriteSecondData, err := json.Marshal(toWriteSecond)
	assert.Nil(t, err)
	toWriteThird := "The data to write!"
	toWriteThirdData := []byte(toWriteThird)

	expected := baseutils.JoinBytes([]byte("\n"), toWriteFirst, toWriteSecondData, toWriteThirdData)
	expected = baseutils.JoinBytes([]byte(""), expected, []byte("\n"))

	err = RecordData(toWriteFirst)
	assert.Nil(t, err)
	err = RecordData(toWriteSecond)
	assert.Nil(t, err)
	err = RecordData(toWriteThird)
	assert.Nil(t, err)
	time.Sleep(3 * time.Second)
	err = file.Close()
	assert.Nil(t, err)

	data, err := ioutil.ReadFile(file.Name())
	assert.Nil(t, err)
	assert.Equal(t, expected, data)
}
