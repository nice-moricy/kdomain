package handle

import (
	"kdomain/utils"
	"os"
	"sync"

	"github.com/molikatty/mlog"
	"github.com/molikatty/molix"
)

type RunDig struct {
	*digSubdomain
	*mlog.Log
	sync.WaitGroup
}

func NewRun(dmMothed, dnMothed uint8) *RunDig {
	return &RunDig{
		Log:          mlog.Logger(),
		digSubdomain: newDig(dmMothed, dnMothed),
	}
}

func (r *RunDig) Run() {
	for {
		server, ok := r.dnsPool()
		if !ok {
			break
		}
		r.Info("[DNS]", utils.FormatString(server, utils.Format, "\n"))
		server = utils.JoinHostPort(server, "53")
		r.Info("[DNS]", utils.FormatString(server, utils.Format, "\n"))

		subItem := utils.Readfile(r.subFile)
		for subItem.Scan() {
			r.Add(1)

			molix.Submit(func() {
				domain := r.getFqdn(subItem.Text())
				if !r.runLookup(domain, server) {
					r.Warning("[SKIP]", utils.FormatString(domain, utils.Format, "\r"))
				}
				r.Done()
			})
		}
	}
	r.Wait()
	r.Close()
	molix.Stop()
}

func (r *RunDig) SetDns(dnsServer string) {
	r.setDns(dnsServer)
}

func (r *RunDig) SetParseDomain(domain string) {
	r.setParseDomain(domain)
}

func (r *RunDig) SetSubdomain(subFile string) {
	r.setSubdomain(subFile)
}

func (r *RunDig) GetAnwser(outfile string) {
	file, err := os.OpenFile(outfile, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		utils.Die(err.Error())
	}

	defer file.Close()
	for data := range r.aw {
		if !utils.RemoveDeduplication(data) {
			continue
		}
		data = utils.FormatString(data, "\n")
		// r.Info("[FOUND]", data)
		file.WriteString(data)
	}
}
