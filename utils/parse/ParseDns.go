package parse

import (
	"kdomain/utils"
	"os"
	"sync/atomic"
)

const (
	File uint8 = iota
	Once
)

type Dns interface {
	MustDnsItem(string) func() (string, bool)
	Close()
}

type fdns struct {
	f *os.File
}
type ddns struct{}

var dnsMothed = map[uint8]Dns{
	File: &fdns{},
	Once: &ddns{},
}

func NewDnsServer(mothed uint8) Dns {
	return dnsMothed[mothed]
}

func (fd *fdns) MustDnsItem(fileName string) func() (string, bool) {
	var err error
	fd.f, err = os.Open(fileName)
	if err != nil {
		utils.Die(utils.FormatString(fileName, " ", err.Error()))
	}
	return utils.Readfile(fd.f)
}

func (dd *ddns) MustDnsItem(dnsServer string) func() (string, bool) {
	s := []string{dnsServer}
	var less atomic.Int64

	return func() (string, bool) {
		index := less.Add(1) - 1
		if index != 0 {
			return "", false
		}

		ss := s[index]
		index++
		return ss, true
	}
}

func (fd *fdns) Close() { fd.f.Close() }
func (dd *ddns) Close() {}
