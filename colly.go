package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mousemin/colly-examples/pipeline"

	"github.com/gocolly/colly/v2"
	"github.com/mousemin/colly-examples/config"
	"github.com/mousemin/colly-extra/configurable"
	"github.com/mousemin/colly-extra/queue"
)

var (
	collector *configurable.Collector
	storage   configurable.Storage
	q         queue.Interface
	factory   *pipeline.Factory
	urlStr    string
	collyName string
	savePath  string
)

func init() {
	flag.StringVar(&urlStr, "url", "", "抓取地址")
	flag.StringVar(&collyName, "name", "", "抓取配置: bswtan")
	flag.StringVar(&savePath, "path", "/Users/mouse/workspace/punch/colly-examples/colly", "存储地址")
}

func Init() (err error) {
	storage, err = configurable.NewDirStorage("", &config.FS)
	if err != nil {
		return err
	}
	q, _ = queue.New(12, nil)
	factory, err = pipeline.New(savePath)
	if err != nil {
		return err
	}
	pipeline.RegisterPipeline(factory)
	collector, err = configurable.New(
		"collector",
		configurable.WithQueue(q),
		configurable.WithConfigStorage(storage),
		configurable.WithPipeline(factory.Pipeline),
	)
	return
}

func Error(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
	os.Exit(1)
}

func main() {
	flag.Parse()
	if len(collyName) == 0 {
		Error("配置名称能为空", collyName)
	}
	if err := Init(); err != nil {
		Error("爬虫初始化失败, err: %s", err.Error())
	}
	conf, err := storage.GetConfig(collyName)
	if err != nil {
		Error("配置名称:%s 的配置不存在 err:%s", collyName, err)
	}
	request := new(colly.Request)
	if len(urlStr) == 0 {
		request = conf.GetBaseRequest()
	} else {
		request = conf.GetBaseRequest(urlStr)
	}
	if err := q.AddRequest(request); err != nil {
		Error("[%s] 爬虫添加配置失败", collyName)
	}
	if err := collector.Init(); err != nil {
		Error("爬虫初始化失败, err: %s", err)
	}
	if err := collector.Start(); err != nil {
		Error("爬虫启动失败, err: %s", err)
	}
	if err := factory.Merge(); err != nil {
		Error("数据merger失败, err: %s", err)
	}
}
