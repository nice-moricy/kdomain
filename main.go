package main

import (
	"errors"
	"flag"
	"fmt"
	"kdomain/handle"
	"kdomain/utils"
	"kdomain/utils/parse"
	"time"
)

var (
	FileName  string
	Domain    string
	Fdomain   string
	Rdomain   string
	Reg       string
	DnsServer string
	DnsFile   string
	OutFile   string
)

var (
	dmMothed uint8
	dnMothed = parse.Once
	server   string
	domain   string
	version  = "0.3"
)

var (
	ErrInput = errors.New("cant parse parameter please check your input")
)

func banner() {
	banner := `
   __ __     __  __  _        ___                  _    
  / //_/__ _/ /_/ /_(_)_ __  / _ \___  __ _  ___ _(_)__ 
 / ,< / _ ./ __/ __/ / // / / // / _ \/  ' \/ _ ./ / _ \
/_/|_|\_,_/\__/\__/_/\_, / /____/\___/_/_/_/\_,_/_/_//_/
                    /___/                               
		kitty domain version: ` + version + `
`
	fmt.Println(banner)
}

// 初始化
func init() {
	banner()
	flag.StringVar(&FileName, "f",
		"subdomains.txt", "set subdomain dict file")
	flag.StringVar(&Domain, "d",
		"", "target domain")
	flag.StringVar(&Fdomain, "fd",
		"", "fuzz domain the placeholder is ? example: w?.baidu.com -> www.baidu.com")
	flag.StringVar(&Rdomain, "rd",
		"", `use regexp fuzz domain must be used -reg example: -rd w?.baidu.com -reg "\?"`)
	flag.StringVar(&Reg, "reg", "", "set regexp for -rd")
	flag.StringVar(&DnsServer, "dns",
		"8.8.8.8", "dns server")
	// flag.StringVar(&DnsFile, "fdns",
	// "", "dns server file")
	flag.StringVar(&OutFile, "o", "", "set out put file")
	flag.Parse()

	server = DnsServer
	switch {
	case Domain != "" && Fdomain == "" && Rdomain == "" && Reg == "":
		dmMothed = parse.Ddomian
		domain = Domain
	case Fdomain != "" && Domain == "" && Rdomain == "" && Reg == "":
		dmMothed = parse.Fdomain
		domain = Fdomain
	case Rdomain != "" && Reg != "" && Domain == "" && Fdomain == "":
		dmMothed = parse.Rdomain
		domain = Rdomain
		parse.Reg = Reg
	}

	if dmMothed == 0 || domain == "" {
		flag.PrintDefaults()
		utils.Die(utils.ErrPlaceholder.Error())
	}

	// if DnsFile != "" {
	// dnMothed = parse.File
	// server = DnsFile
	// }
	if OutFile == "" {
		OutFile = domain + ".txt"
	}
}

func main() {
	start := time.Now()
	run := handle.NewRun(dmMothed, dnMothed)
	run.SetDns(server)
	run.SetSubdomain(FileName)
	run.SetParseDomain(domain)
	go run.GetAnwser(OutFile)
	run.Run()
	run.Info("[DONE]",
		utils.FormatString(time.Since(start).String(), utils.Format, "\n"),
	)
}
