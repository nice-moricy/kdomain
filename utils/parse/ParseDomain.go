package parse

import (
	"kdomain/utils"
	"regexp"
	"strings"
)

const (
	Fdomain uint8 = iota + 1
	Rdomain
	Ddomian
)

var Reg string

type Domain interface {
	Parse(string) []string
	Format(string, []string) string
}

type formatDomain struct{}
type regexpDomain struct{}
type defaultDomain struct{}

var domainMothed = map[uint8]Domain{
	Fdomain: formatDomain{},
	Rdomain: regexpDomain{},
	Ddomian: defaultDomain{},
}

func NewParseDomain(mode uint8) Domain {
	return domainMothed[mode]
}

func (f formatDomain) Parse(domain string) (d []string) {
	d = strings.Split(domain, "?")
	if len(d) < 2 {
		utils.Die(utils.FormatString(domain, " ", utils.ErrPlaceholder.Error(), "\n"))
	}
	return d
}

func (f formatDomain) Format(subdomain string, domain []string) string {
	return strings.Join(domain, subdomain)
}

func (r regexpDomain) Parse(domain string) []string {
	reg, err := regexp.Compile(Reg)
	if err != nil {
		utils.Die(utils.FormatString(Reg, " ", err.Error(), "\n"))
	}

	return reg.Split(domain, 20)
}

func (f regexpDomain) Format(subdomain string, domain []string) string {
	return strings.Join(domain, subdomain)
}

func (d defaultDomain) Parse(domain string) []string {
	return []string{domain}
}

func (f defaultDomain) Format(subdomain string, domain []string) string {
	return utils.FormatString(
		append(domain[:0], append([]string{subdomain, "."}, domain...)...)...,
	)
}
