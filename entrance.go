package fisher

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/akamensky/argparse"
	"github.com/cfh0081/runutils"
)

type CrawlAction = func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File)

func checkArgments(args []string) error {
	argReg := regexp.MustCompile(`^[^=]+=.+$`) // 正则匹配***=***样式
	if cnt := len(args); cnt > 0 {
		for _, v := range args {
			isMatch := argReg.MatchString(v)
			if !isMatch {
				return fmt.Errorf("argment %v is invalid", v)
			}
		}
	}

	return nil
}

func getCustomArgs(args []string) map[string]string {
	result := make(map[string]string)
	argReg := regexp.MustCompile(`^([^=]+)=(.+)$`) // 正则匹配***=***样式
	if cnt := len(args); cnt > 0 {
		for _, v := range args {
			rtn := argReg.FindStringSubmatch(v)
			result[rtn[1]] = rtn[2]
		}
	}

	return result
}

func toCrawlWithArgs(args []string, action CrawlAction) (err error) {
	var emptyFile os.File
	var outputFile *os.File = nil
	proxyReqUrl := ""
	isHeadless := false
	err = nil

	// Create new parser object
	parser := argparse.NewParser("crawl", "crawl the information needed.")

	// 用于添加自定义参数
	var customArgs *[]string = parser.StringList("a", "arg", &argparse.Options{Required: false, Validate: checkArgments, Help: "Custom arguments."})

	// 获取os.File用于读取请求获取代理地址的链接
	proxyFile := parser.File("", "proxy", os.O_RDONLY, 0644, &argparse.Options{Required: false, Help: "The file stored the url to get the proxy."})

	headlessFlag := parser.Flag("", "headless", &argparse.Options{Required: false, Help: "Specify the chrome browser to enable headless."})

	// 获取os.File用于读取请求获取代理地址的链接
	outputFileA := parser.File("o", "output_append", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666, &argparse.Options{Required: false, Help: "The file to append the crawled data."})

	// 获取os.File用于读取请求获取代理地址的链接
	outputFileB := parser.File("O", "output_new", os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666, &argparse.Options{Required: false, Help: "The file as a new one to store the crawled data."})

	// Parse input
	err = parser.Parse(args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return
	}

	if *proxyFile != emptyFile {
		defer proxyFile.Close()

		if data, errRead := ioutil.ReadAll(proxyFile); errRead == nil {
			err = errRead
			proxyReqUrl = string(data)
			// fmt.Printf("proxyReqUrl: %v\n", proxyReqUrl)
		} else {
			err = errRead
			fmt.Printf("%v call ioutil.ReadAll with error %v.\n", runutils.RunFuncName(), errRead)
			return
		}
	}

	isHeadless = *headlessFlag
	customMap := getCustomArgs(*customArgs)

	if *outputFileA != emptyFile && *outputFileB != emptyFile {
		defer outputFileA.Close()
		defer outputFileB.Close()

		errInfo := "can only choose between -o and -O item"
		err = errors.New(errInfo)
		fmt.Print(errInfo)
		return
	} else {
		if *outputFileA != emptyFile {
			outputFile = outputFileA
		} else if *outputFileB != emptyFile {
			outputFile = outputFileB
		}

		if outputFile != nil {
			defer outputFile.Close()
		}
	}

	action(proxyReqUrl, isHeadless, customMap, outputFile)

	return
}

func ToCrawl(action CrawlAction) {
	toCrawlWithArgs(os.Args, action)
}
