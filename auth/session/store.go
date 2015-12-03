package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/raintreeinc/knowledgebase/kb"
)

type Token [64]byte

var ZeroToken = Token{}

func (token *Token) String() string {
	return base64.StdEncoding.EncodeToString((*token)[:])
}

func TokenFromString(s string) (Token, error) {
	var token Token
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return token, err
	}
	if len(data) != len(token) {
		return token, errors.New("Invalid token length.")
	}
	copy(token[:], data)
	return token, nil
}

func GenerateToken() (Token, error) {
	var token Token
	_, err := io.ReadFull(rand.Reader, token[:])
	return token, err
}

type Info struct {
	User kb.User

	Creation   time.Time
	LastAccess time.Time
}

type Store struct {
	mu      sync.Mutex
	entries map[Token]*Info
	maxage  time.Duration
}

func NewStore(maxage time.Duration) *Store {
	store := &Store{
		entries: make(map[Token]*Info),
		maxage:  maxage,
	}

	go store.purge()
	return store
}

func (store *Store) purge() {
	for range time.Tick(time.Minute) {
		store.PurgeOld()
	}
}

func (store *Store) PurgeOld() {
	deadline := time.Now().Add(-store.maxage)

	store.mu.Lock()
	for tok, info := range store.entries {
		if info.LastAccess.Before(deadline) {
			delete(store.entries, tok)
		}
	}
	store.mu.Unlock()
}

func (store *Store) New(user kb.User) (Token, error) {
	token, err := GenerateToken()
	if err != nil {
		return ZeroToken, err
	}

	info := &Info{
		User:       user,
		Creation:   time.Now(),
		LastAccess: time.Now(),
	}

	store.mu.Lock()
	store.entries[token] = info
	store.mu.Unlock()

	return token, nil
}

func (store *Store) Load(token Token) (user kb.User, ok bool) {
	var info *Info

	store.mu.Lock()
	info, ok = store.entries[token]
	if ok {
		info.LastAccess = time.Now()
		user = info.User
	}
	store.mu.Unlock()

	return user, ok
}

func (store *Store) Delete(token Token) {
	store.mu.Lock()
	delete(store.entries, token)
	store.mu.Unlock()
}
