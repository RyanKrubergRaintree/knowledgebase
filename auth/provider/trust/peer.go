package trust

import (
	"crypto/hmac"
	"crypto/sha1"
	"errors"
	"net/url"
	"time"
)

var (
	// Maximum skew in times between servers
	MaxRequestSkew = time.Minute
	// Maximum authorization field size
	MaxAuthorizationSize = 10 << 10
)

var (
	ErrUnauthorized = errors.New("Invalid key.")
	ErrTimeSkewed   = errors.New("Time is skewed.")
)

const timeLayout = time.RFC3339

// Peer is trusted server information
type Peer struct {
	// key that will be used to authenticate
	Key []byte
}

// Sign signs the request from id.
func (peer Peer) Sign(id string) (string, error) {
	ts := time.Now().Format(timeLayout)
	nonce, err := generateNonce()
	if err != nil {
		return "", err
	}

	mac := Sign(SerializeValues(id, ts, nonce), peer.Key)

	v := url.Values{}
	v.Set("id", id)
	v.Set("ts", ts)
	v.Set("nonce", nonce)
	v.Set("mac", string(mac))

	return "KB " + v.Encode(), nil
}

// Verify authenticates whether the request is trusted.
func (peer Peer) Verify(auth string) (id string, err error) {
	now := time.Now()

	if len(auth) > MaxAuthorizationSize {
		return "", ErrUnauthorized
	}

	v, err := url.ParseQuery(auth)
	if err != nil {
		return "", ErrUnauthorized
	}

	id, ts, nonce, mac := v.Get("id"), v.Get("ts"), v.Get("nonce"), v.Get("mac")

	if id == "" || ts == "" || nonce == "" {
		return "", ErrUnauthorized
	}

	sentAt, err := time.Parse(timeLayout, ts)
	if err != nil {
		return "", ErrUnauthorized
	}
	skew := now.Sub(sentAt)
	if skew > MaxRequestSkew || -skew > MaxRequestSkew {
		return "", ErrTimeSkewed
	}

	if err := nonceHistory.Record(nonce, now); err != nil {
		return "", ErrUnauthorized
	}

	expected := Sign(SerializeValues(id, ts, nonce), peer.Key)
	if !hmac.Equal(expected, []byte(mac)) {
		return "", ErrUnauthorized
	}

	return id, nil
}

func Sign(serialized, key []byte) []byte {
	m := hmac.New(sha1.New, key)
	m.Write(serialized)
	return m.Sum(nil)
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
