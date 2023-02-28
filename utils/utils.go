package utils

import (
	"fmt"
	"net"
	"os"
	"strings"
	"unsafe"
)

func FormatString(s []string) string {
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

func String(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func Die(s interface{}) {
	fmt.Fprintln(os.Stderr, s)
	os.Exit(1)
}
