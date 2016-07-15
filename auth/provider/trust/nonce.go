package trust

import (
	"crypto/rand"
	"errors"
	"time"
)

var nonceHistory = make(nonces, 0, maxNonces)

var (
	nonceUsed      = errors.New("nonce has been used")
	nonceTableFull = errors.New("nonce table is full")
)

type nonce [6]byte

const maxNonces = 60 << 10

type expiringNonce struct {
	Nonce nonce
	Time  time.Time
}

type nonces []expiringNonce

func (ns *nonces) purge(now time.Time) {
	xs := *ns
	for i := len(xs) - 1; i >= 0; i-- {
		delta := xs[i].Time.Sub(now)
		if delta > MaxRequestSkew || -delta > MaxRequestSkew {
			xs[i] = xs[len(xs)-1]
			xs = xs[:len(xs)-1]
		}
	}
	*ns = xs
}

func (ns *nonces) Record(xs string, now time.Time) error {
	ns.purge(now)

	var x nonce
	copy(x[:], []byte(xs))

	for _, n := range *ns {
		if n.Nonce == x {
			return nonceUsed
		}
	}

	if len(*ns) == cap(*ns) {
		return nonceTableFull
	}

	*ns = append(*ns, expiringNonce{
		Nonce: x,
		Time:  now,
	})
	return nil
}

func generateNonce() (string, error) {
	var nonce nonce
	n, err := rand.Read(nonce[:])
	if n != len(nonce) {
		return "", errors.New("failed to read enough random")
	}
	if err != nil {
		return "", err
	}
	return string(nonce[:]), nil
}
