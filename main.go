package main

import (
	"bufio"
	"errors"
	"flag"
	"kdomain/utils"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/miekg/dns"
	"github.com/molikatty/mlog"
	"github.com/molikatty/molix"
)

var (
	FileName  string
	Domain    string
	DnsServer string
	OutFile   string
)

var (
	async  sync.WaitGroup
	log    = mlog.Logger()
	format = strings.Repeat(" ", 50)
)

var (
	ErrInput = errors.New("cant parse parameter please check your input")
)

// 初始化
func init() {
	banner()
	flag.StringVar(&FileName, "f", "subdomains.txt", "set subdomain dict file")
	flag.StringVar(&Domain, "d", "", "target domain")
	flag.StringVar(&DnsServer, "dns", "8.8.8.8", "target domain")
	flag.StringVar(&OutFile, "o", "result.txt", "set out put file")
	flag.Parse()
	if Domain == "" {
		flag.PrintDefaults()
		utils.Die(ErrInput)
	}
}

func main() {
	// 给dns服务器IP添加端口
	DnsServer = utils.JoinHostPort(DnsServer, "53")
	aw := make(chan string, 1e3)
	f, err := os.Open(FileName)
	if err != nil {
		utils.Die(err.Error())
	}
	go result(aw)
	start := time.Now()
	fqdnItem := getFqdn(f)
	for {
		fqdn, ok := fqdnItem()
		if !ok {
			f.Close()
			break
		}

		submit(fqdn, DnsServer, dns.TypeA, aw)
	}
	async.Wait()
	molix.Stop()
	close(aw)
	log.Info("[DONE]", time.Since(start).String()+format+"\n")
}

// 字典文件迭代器
func getFqdn(f *os.File) func() (string, bool) {
	buf := bufio.NewScanner(f)
	return func() (string, bool) {
		if !buf.Scan() {
			return "", false
		}
		return dns.Fqdn(utils.FormatString(
			[]string{
				utils.String(buf.Bytes()),
				".",
				Domain,
			},
		)), true
	}
}

// 提交任务
func submit(fqdn, dnsServer string, dnsType uint16, aw chan string) {
	async.Add(1)
	molix.Submit(func() {
		defer async.Done()
		item := lookup(fqdn, dnsServer, dnsType)
		if item == nil {
			log.Info("[SKIP]", utils.FormatString([]string{fqdn, format, "\r"}))
			return
		}
		for {
			data, ok := item()
			if !ok {
				break
			}

			aw <- data
		}
	})
}

// 初始化DNS结果的迭代器
func lookup(fqdn, dnsServer string, dnsType uint16) func() (string, bool) {
	r, err := dns.Exchange(
		(&dns.Msg{}).SetQuestion(fqdn, dnsType),
		dnsServer,
	)
	if err != nil {
		log.Err("[ERROR]", utils.FormatString([]string{
			err.Error(),
			"\r",
		}))
		return nil
	}

	n := len(r.Answer) - 1
	if n < 1 {
		log.Err("[WARNING]", utils.FormatString([]string{
			fqdn,
			"No answer\r",
		}))
		return nil
	}

	var (
		less atomic.Int64
		max  = int64(n)
	)
	return func() (string, bool) {
		index := less.Add(1) - 1
		if index > max {
			less.Store(max)
			less.Add(1)
			return "", false
		}

		return r.Answer[index].String(), true
	}
}

// 导出结果
func result(aw chan string) {
	file, err := os.OpenFile(OutFile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		utils.Die(err)
	}

	defer file.Close()
	for data := range aw {
		data = utils.FormatString([]string{data, "\n"})
		log.Info("[FOUND]", data)
		file.WriteString(data)
	}
}

func banner() {
	banner := `
 ____  __.__  __    __           ________                        .__        
|   |/ _|__|/  |__/  |_ ___.__. \______ \   ____   _____ _____  |__| ____  
|     < |  \   __\   __<   |  |  |    |  \ /  _ \ /     \\__  \ |  |/    \ 
|   |  \|  ||  |  |  |  \___  |  |    /   (  <_> )  Y Y  \/ __ \|  |   |  \
|___|__ \__||__|  |__|  / ____| /_______  /\____/|__|_|  (____  /__|___|  /
        \/               \/              \/             \/     \/        \/ 
		kitty domain version: 0.1
`
	log.OutMsg(log.Stdout, "", banner)
}
