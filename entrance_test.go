package fisher

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckArgments(t *testing.T) {

	arg0 := []string{"a", "b", "c"}
	rtn0 := checkArgments(arg0)
	assert.NotNil(t, rtn0)

	arg1 := []string{"=a"}
	rtn1 := checkArgments(arg1)
	assert.NotNil(t, rtn1)

	arg2 := []string{"cx=bc", "a-)x="}
	rtn2 := checkArgments(arg2)
	assert.NotNil(t, rtn2)

	arg3 := []string{"a=b", "b==", "c=.", "a-)x=="}
	rtn3 := checkArgments(arg3)
	assert.Nil(t, rtn3)
}

func TestGetCustomArgs(t *testing.T) {
	expected := map[string]string{"a": "b", "b": "=", "c": ".", "a-)x": "="}
	arg1 := []string{"a=b", "b==", "c=.", "a-)x=="}
	rtn1 := getCustomArgs(arg1)
	assert.Equal(t, expected, rtn1)
}

func TestToCrawlWithArgs(t *testing.T) {
	// 不带任何参数的使用
	args := []string{"crawler"}
	err := toCrawlWithArgs(args, func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {
		assert.Equal(t, "", proxyReqUrl)
		assert.False(t, isHeadless)

		expectedMap := make(map[string]string)
		assert.Equal(t, expectedMap, customMap)
		assert.Nil(t, outputFile)
	})
	assert.Nil(t, err)

	// 验证代理地址文件配置、无界面模式、新建截断文件式写入和自定义参数
	// 临时创建一个文件并写入相应内容用于验证代理文件获取
	file, err := ioutil.TempFile("", "proxy-*.txt")
	assert.Nil(t, err)
	// 确保程序结束时删除临时文件
	defer os.Remove(file.Name())
	recordFile, err := ioutil.TempFile("", "output-*.jl")
	assert.Nil(t, err)
	defer os.Remove(recordFile.Name())

	// 写具体内容到文件中
	urlInfo := `https://www.sohu.com`
	_, err = file.Write([]byte(urlInfo))
	assert.Nil(t, err)
	args = []string{"crawler", "--proxy", file.Name(), "--headless", "-O", recordFile.Name(), "-a", "begin=a", "-a", "a#b=="}
	err = toCrawlWithArgs(args, func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {
		assert.Equal(t, urlInfo, proxyReqUrl)
		assert.True(t, isHeadless)

		expectedMap := map[string]string{"begin": "a", "a#b": "="}
		assert.Equal(t, expectedMap, customMap)
		assert.NotNil(t, outputFile)
	})
	assert.Nil(t, err)

	// 验证增量写文件
	recordAddFile, err := ioutil.TempFile("", "output-*.jl")
	assert.Nil(t, err)
	defer os.Remove(recordAddFile.Name())
	args = []string{"crawler", "-o", recordAddFile.Name()}
	err = toCrawlWithArgs(args, func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {
		assert.Equal(t, "", proxyReqUrl)
		assert.False(t, isHeadless)

		expectedMap := make(map[string]string)
		assert.Equal(t, expectedMap, customMap)
		assert.NotNil(t, outputFile)
	})
	assert.Nil(t, err)

	// 验证同时配置两类输出文件的情况
	record01File, err := ioutil.TempFile("", "output-*.jl")
	assert.Nil(t, err)
	defer os.Remove(record01File.Name())
	recordAdd01File, err := ioutil.TempFile("", "output-*.jl")
	assert.Nil(t, err)
	defer os.Remove(recordAdd01File.Name())
	args = []string{"crawler", "-O", record01File.Name(), "-o", recordAdd01File.Name()}
	err = toCrawlWithArgs(args, func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {

	})
	assert.NotNil(t, err)

	// 验证错误输入情况
	args = []string{"crawler", "-proxy", "--headless"}
	err = toCrawlWithArgs(args, func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {

	})
	assert.NotNil(t, err)
}
