package handle

import (
	"kdomain/utils"
	"kdomain/utils/parse"
	"os"
	"sync/atomic"

	"github.com/miekg/dns"
	"github.com/molikatty/mlog"
)

type digSubdomain struct {
	aw        chan string
	domain    []string
	dnsPool   func() (string, bool)
	subFile   *os.File
	parse.Domain
	parse.Dns
	*mlog.Log
}

func newDig(dmMothed, dnMothed uint8) *digSubdomain {
	return &digSubdomain{
		aw:     make(chan string, 1e3),
		Domain: parse.NewParseDomain(dmMothed),
		Dns:    parse.NewDnsServer(dnMothed),
	}
}

func (dig *digSubdomain) setDns(s string) {
	dig.dnsPool = dig.MustDnsItem(s)
}

func (dig *digSubdomain) setParseDomain(s string) {
	dig.domain = dig.Parse(s)
}

func (dig *digSubdomain) setSubdomain(file string) {
	var err error
	dig.subFile, err = os.Open(file)
	if err != nil {
		utils.Die(utils.FormatString(file, " ", err.Error()))
	}
}

func (dig *digSubdomain) close() error {
	close(dig.aw)
	return dig.subFile.Close()
}

func (dig *digSubdomain) getFqdn(s string) string {
	return dns.Fqdn(dig.Format(s, dig.domain))
}

func (dig *digSubdomain) runLookup(fqdn, dnsServer string) bool {
	item := dig.lookup(fqdn, dnsServer, dns.TypeA)
	if item == nil {
		return false
	}
	for {
		data, ok := item()
		if !ok {
			break
		}

		dig.aw <- data
	}

	return true
}

func (dig *digSubdomain) Fqdn(sub string, domain []string) string {
	return dns.Fqdn(dig.Format(sub, domain))
}

// 初始化DNS结果的迭代器
func (dig *digSubdomain) lookup(fqdn, dnsServer string, dnsType uint16) func() (string, bool) {
	r, err := dns.Exchange(
		(&dns.Msg{}).SetQuestion(fqdn, dnsType),
		dnsServer,
	)
	if err != nil {
		return nil
	}

	n := len(r.Answer) - 1
	if n < 1 {
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
