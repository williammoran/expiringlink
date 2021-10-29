package expiringlink

import (
	"testing"
	"time"
)

func TestMaxRounds(t *testing.T) {
	el := ExpiringLink{
		Epoch:     epoch,
		Expire:    10 * time.Second,
		Rounds:    5,
		MaxRounds: 4,
	}
	hash := el.Generate("meow")
	if err := el.Check(hash, "meow"); err != CorruptHashError {
		t.Fatalf("Should have been corrupt but %s", err.Error())
	}
}

func TestInvalidHash(t *testing.T) {
	el := ExpiringLink{
		Epoch:  epoch,
		Expire: 10 * time.Second,
		Rounds: 5,
	}
	if err := el.Check("ff", "25"); err != CorruptHashError {
		t.Log(err.Error())
		t.Fatal("Should have returned CorruptHashError for invalid hash")
	}
}

func TestHashCheck(t *testing.T) {
	el := ExpiringLink{
		Epoch:  epoch,
		Expire: 24 * time.Hour,
		Rounds: 5,
	}
	for _, val := range hashTestStrings {
		hash := el.Generate(val)
		if err := el.Check(hash, val); err != nil {
			t.Logf("'%s' hashed to '%s' but didn't check", val, hash)
			t.Logf("%s", err)
			t.Fail()
		}
		if el.Check(hash, "UnusedSecretValue") == nil {
			t.Logf("'%s' incorrectly checked", val)
			t.Fail()
		}
	}
}

func TestExpire(t *testing.T) {
	el := ExpiringLink{
		Epoch:  epoch,
		Expire: 2 * time.Second,
		Rounds: 5,
	}
	for _, val := range hashTestStrings {
		hash := el.Generate(val)
		if err := el.Check(hash, val); err != nil {
			t.Logf("Should be valid but %s", err.Error())
			t.Fail()
		}
		time.Sleep(3 * time.Second)
		if el.Check(hash, val) != HashExpiredError {
			t.Log("Didn't expire")
			t.Fail()
		}
	}
}

var hashTestStrings = []string{
	"meow",
	"a reasonlbly but not insanely long string",
	"UTF-8 stuff: 😁🚧🚀",
}

var epoch time.Time

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2021-09-29T12:45:26.371Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		panic(err)
	}
	epoch = t
}
