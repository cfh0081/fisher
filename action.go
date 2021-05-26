package fisher

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/cfh0081/runutils"
	"github.com/chromedp/chromedp"
)

type Crawler interface {
	Crawl(ctx context.Context) error
	ParseArgs(argMap map[string]string) error
}

func GetAction(crawler Crawler) CrawlAction {
	return func(proxyReqUrl string, isHeadless bool, customMap map[string]string, outputFile *os.File) {
		opts := make([]chromedp.ExecAllocatorOption, 0, 16)
		if proxyReqUrl != "" {
			proxy, err := GetSocks5Proxy(context.Background(), proxyReqUrl)
			if err != nil {
				fmt.Printf("%v call GetSocks5Proxy with error %v.\n", runutils.RunFuncName(), err)
				return
			} else {
				opts = append(opts, chromedp.ProxyServer(proxy))
				fmt.Printf("to process with proxy %v.\n", proxy)
			}
		}

		if isHeadless {
			opts = append(opts,
				chromedp.Flag("headless", true),
				// Like in Puppeteer.
				chromedp.Flag("hide-scrollbars", true),
				chromedp.Flag("mute-audio", true),
			)
		} else {
			opts = append(opts,
				chromedp.Flag("headless", false),
			)
		}

		err := crawler.ParseArgs(customMap)
		if err != nil {
			fmt.Printf("%v call crawler.ParseArgs with error %v.\n", runutils.RunFuncName(), err)
			return
		}

		SetStorageTarget(outputFile)

		RunWithCrawler(crawler, opts...)
	}
}

func RunWithCrawler(crawler Crawler, opts ...chromedp.ExecAllocatorOption) {
	targetOpts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("start-maximized", true), // 最大化运行（全屏窗口）
		chromedp.Flag("incognito", true),       // 隐身模式（无痕模式）
		chromedp.Flag("log-level", "2"),        // 日志级别 ( info(default) = 0 warning = 1 LOG_ERROR = 2 LOG_FATAL = 3 )
		chromedp.Flag("lang", `zh-CN,zh,zh-TW,en-US,en`),
		chromedp.Flag("enable-automation", false), // 禁用浏览器正在被自动化程序控制的提示
	)

	targetOpts = append(targetOpts, opts...)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), targetOpts...)
	defer cancel()

	// create chrome instance
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()

	// 开启保存数据服务
	go StartStorageServer(ctx)

	err := crawler.Crawl(ctx)
	if err != nil {
		fmt.Printf("%v call crawler.Crawl with error %v.\n", runutils.RunFuncName(), err)
	}
}
