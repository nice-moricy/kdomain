package utils

import (
	"bufio"
	"errors"
	"net"
	"os"
	"strings"
	"unsafe"

	"github.com/molikatty/mlog"
)

var (
	deduplication = make(map[string]struct{})
	Format        = strings.Repeat(" ", 30)
)

var (
	ErrPlaceholder = errors.New("placecholder error please check your input")
)

func FormatString(s ...string) string {
	var (
		b strings.Builder
		n int
	)

	for i := range s {
		n += len(s[i])
	}

	b.Grow(n)

	for i := range s {
		b.WriteString(s[i])
	}

	return b.String()
}

func JoinHostPort(host, port string) string {
	return net.JoinHostPort(host, port)
}

func ToString(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func Readfile(f *os.File) func() (string, bool) {
	buf := bufio.NewScanner(f)
	return func() (string, bool) {
		return buf.Text(), buf.Scan()
	}
}

func RemoveDeduplication(s string) bool {
	if _, ok := deduplication[s]; ok {
		return false
	}

	deduplication[s] = struct{}{}
	return true
}

func Die(s string) {
	mlog.Logger().Err("[ERROR]", s)
	os.Exit(1)
}
