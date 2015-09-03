package trust

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	// Maximum skew in times between servers
	MaxRequestSkew = time.Minute
	// Maximum authorization field size
	MaxAuthorizationSize = 10 << 10
)

var (
	PeerNotTrusted    = errors.New("Peer is not trusted.")
	RequestNotTrusted = errors.New("Request is not trusted.")
	TimeSkewed        = errors.New("Time is skewed.")
)

const timeLayout = time.RFC3339

// Peer is trusted server information
type Peer struct {
	// key that will be used to authenticate
	Key []byte
}

// Sign signs the request from id.
// After calling this function request mustn't be further modified.
func (peer Peer) Sign(id string, req *http.Request) error {
	ts := time.Now().Format(timeLayout)
	nonce, err := generateNonce()
	if err != nil {
		return err
	}

	mac := Sign(SerializeParams(id, ts, nonce, req), peer.Key)

	v := url.Values{}
	v.Set("id", id)
	v.Set("ts", ts)
	v.Set("nonce", nonce)
	v.Set("mac", string(mac))

	req.Header.Set("Authorization", "KB "+v.Encode())
	return nil
}

// Verify authenticates whether the request is trusted.
func (peer Peer) Verify(req *http.Request) (id string, err error) {
	now := time.Now()

	auth := req.Header.Get("Authorization")
	if auth == "" || !strings.HasPrefix(auth, "KB ") {
		return "", RequestNotTrusted
	}

	if len(auth) > MaxAuthorizationSize {
		return "", RequestNotTrusted
	}

	v, err := url.ParseQuery(auth[3:])
	if err != nil {
		return "", RequestNotTrusted
	}

	id, ts, nonce, mac := v.Get("id"), v.Get("ts"), v.Get("nonce"), v.Get("mac")

	if id == "" || ts == "" || nonce == "" {
		return "", RequestNotTrusted
	}

	sentAt, err := time.Parse(timeLayout, ts)
	if err != nil {
		return "", RequestNotTrusted
	}
	skew := now.Sub(sentAt)
	if skew > MaxRequestSkew || -skew > MaxRequestSkew {
		return "", TimeSkewed
	}

	if err := nonceHistory.Record(nonce, now); err != nil {
		return "", RequestNotTrusted
	}

	expected := Sign(SerializeParams(id, ts, nonce, req), peer.Key)
	if !hmac.Equal(expected, []byte(mac)) {
		return "", RequestNotTrusted
	}

	return id, nil
}

func Sign(serialized, key []byte) []byte {
	m := hmac.New(sha1.New, key)
	m.Write(serialized)
	return m.Sum(nil)
}

func SerializeParams(id, ts, nonce string, req *http.Request) []byte {
	return SerializeValues(
		id, ts, nonce,
		req.Host, req.URL.Path,
	)
}

func SerializeValues(values ...string) []byte {
	t := 0
	for _, v := range values {
		t += len(v) + 1
	}
	r := make([]byte, 0, t)
	for _, v := range values {
		r = append(r, []byte(v)...)
		r = append(r, '\n')
	}
	return r
}
