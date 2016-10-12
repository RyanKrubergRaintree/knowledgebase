package trust

import (
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
	"net/url"
	"time"
)

var (
	// Maximum skew in times between servers
	MaxRequestSkew = 4 * time.Minute
	// Maximum authorization field size
	MaxAuthorizationSize = 10 << 10
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

// Verify checks whether the request is trusted.
func (peer Peer) Verify(auth string) (id string, err error) {
	now := time.Now()

	if len(auth) > MaxAuthorizationSize {
		return "", fmt.Errorf("Auth failed: invalid request.")
	}

	v, err := url.ParseQuery(auth)
	if err != nil {
		return "", fmt.Errorf("Auth failed: invalid query.")
	}

	id, ts, nonce, mac := v.Get("id"), v.Get("ts"), v.Get("nonce"), v.Get("mac")

	if id == "" || ts == "" || nonce == "" {
		return "", fmt.Errorf("Auth \"%s\" failed: missing entries.", id)
	}

	sentAt, err := time.Parse(timeLayout, ts)
	if err != nil {
		return "", fmt.Errorf("Auth \"%s\" failed: wrong time layout.", id)
	}
	skew := now.Sub(sentAt)
	if skew > MaxRequestSkew || -skew > MaxRequestSkew {
		return "", fmt.Errorf("Auth \"%s\" failed: time skewed (%v).", id, skew)
	}

	if err := nonceHistory.Record(nonce, now); err != nil {
		return "", fmt.Errorf("Auth \"%s\" failed: nonce check failed (%v).", id, err)
	}

	expected := Sign(SerializeValues(id, ts, nonce), peer.Key)
	if !hmac.Equal(expected, []byte(mac)) {
		return "", fmt.Errorf("Auth \"%s\" failed: invalid signature.", id)
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
