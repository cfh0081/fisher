package fisher

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

// 先对 json 格式进行 struct 结构定义
type ProxyInfo struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
}

type ProxyResponse struct {
	Code    int           `json:"code"`
	Data    [16]ProxyInfo `json:"data"`
	Msg     string        `json:"msg"`
	Success bool          `json:"success"`
}

// 获取代理地址信息
func GetProxyInfo(ctx context.Context, url string) (ret ProxyInfo, err error) {
	// req, _ := http.NewRequest("GET", url, nil)
	// req = req.WithContext(ctx)

	req, newErr := http.NewRequestWithContext(ctx, "GET", url, nil)
	if newErr != nil {
		return ret, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ret, err
	}
	defer resp.Body.Close()

	var rtnData ProxyResponse

	//返回的状态码
	bodyContent, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(bodyContent, &rtnData); err == nil {
		if rtnData.Code == 0 {
			ret = rtnData.Data[0]
			return ret, err
		} else {
			newErr = errors.New(rtnData.Msg)
			return ret, newErr
		}

		// fmt.Printf("target: %v\n", ret)
	}

	return ret, err
}

// 获取SOCKS5代理地址
func GetSocks5Proxy(ctx context.Context, url string) (ret string, err error) {
	cntButt := 3
	cnt := 0
	var proxy string
	var proxyInfo ProxyInfo
	for {
		if proxyInfo, err = GetProxyInfo(ctx, url); err == nil {
			proxy = fmt.Sprintf("SOCKS5://%v:%v", proxyInfo.IP, proxyInfo.Port)
			// fmt.Printf("proxy: %v.\n", proxy)
			break
		} else {
			// 至少要等待1秒以上才能重新请求，避免被服务器拒绝
			// fmt.Printf("call GetProxyInfo failed with %v!\n", err)
			time.Sleep(2 * time.Second)
		}

		cnt += 1
		if cnt >= cntButt {
			break
		}
	}

	return proxy, err
}

// 验证代码逻辑，先注释，后续有需要可参考
// func proxyCheck() {
// 	url := "http://http.tiqu.letecs.com/getip3?num=1&type=2&pro=440000&city=441900&yys=0&port=2&time=1&ts=1&ys=1&cs=1&lb=1&sb=0&pb=4&mr=1&regions="
// 	ret, err := GetSocks5Proxy(context.Background(), url)
// 	if err != nil {
// 		fmt.Printf("err: %v.\n", err)
// 	}

// 	fmt.Printf("ret: %v.\n", ret)
// }
