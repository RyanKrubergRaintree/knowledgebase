// +build ignore

package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"

	"github.com/raintreeinc/knowledgebase/auth/trust"
)

var (
	key = flag.String("key", "123456789", "key used for signing")

	id    = flag.String("id", "User:Company", "id to be signed")
	ts    = flag.String("ts", "2006-01-02T15:04:05Z07:00", "timestamp to be signed")
	nonce = flag.String("nonce", "1234567890", "")
	host  = flag.String("host", "example.com", "")
	path  = flag.String("path", "/system/auth/callback/community", "")
)

func main() {
	flag.Parse()
	signature := trust.Sign(
		trust.SerializeValues(
			*id, *ts, *nonce,
			*host, *path,
		),
		[]byte(*key),
	)

	fmt.Println(signature)
	fmt.Println(hex.Dump(signature))
	fmt.Println(base64.StdEncoding.EncodeToString(signature))
}
