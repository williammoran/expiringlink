package expiringlink

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ExpiringLink struct {
	// Epoch defines the oldest possible date the system
	// could use. It can be any date in the past, but once
	// selected must NEVER be changed, as doing so could
	// cause previously expired tokens to be valid again.
	// Dates closer to the present will result in shorter
	// hashes
	Epoch time.Time
	// Expire is how long before a generated token expires
	// (in seconds) This value can be changed at any time
	// and will affect any tokens generate after the change.
	Expire time.Duration
}

func (e *ExpiringLink) Generate(secret string) string {
	expire := time.Since(e.Epoch) + e.Expire
	expTime := uint64(expire / time.Second)
	hash := sign(expTime, secret)
	return formatHash(expTime, hash)
}

type constError string

func (c constError) Error() string { return string(c) }

const (
	CorruptHashError = constError("Corrupt Hash")
	HashExpiredError = constError("Hash expired")
	InvalidHashError = constError("Hash did not validate")
)

func (e *ExpiringLink) Check(hash, secret string) error {
	part := strings.Split(hash, "g")
	if len(part) != 2 {
		return CorruptHashError
	}
	ts, err := strconv.ParseInt(part[0], 16, 64)
	if err != nil {
		return CorruptHashError
	}
	if e.Epoch.Add(time.Second * time.Duration(ts)).Before(time.Now()) {
		return HashExpiredError
	}
	genHash := sign(uint64(ts), secret)
	genFormatted := formatHash(uint64(ts), genHash)
	if genFormatted == hash {
		return nil
	}
	return InvalidHashError
}

func formatHash(age uint64, hash string) string {
	return fmt.Sprintf("%xg%s", age, hash)
}

func sign(expire uint64, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha1.New, key)
	message := fmt.Sprintf("%x", expire)
	h.Write([]byte(message))
	return hex.EncodeToString(h.Sum(nil))
}
