package expiringlink

import (
	"testing"
	"time"
)

func TestHashCheck(t *testing.T) {
	el := ExpiringLink{
		Epoch:  epoch,
		Expire: 24 * time.Hour,
		Rounds: 5,
	}
	for _, val := range hashTestStrings {
		hash := el.Generate(val)
		t.Log(hash)
		if !el.Check(hash, val) {
			t.Logf("'%s' hashed to '%s' but didn't check", val, hash)
			t.Fail()
		}
		if el.Check(hash, "UnusedSecretValue") {
			t.Logf("'%s' incorrectly checked", val)
			t.Fail()
		}
	}
	t.Fail()
}

func TestExpire(t *testing.T) {
	el := ExpiringLink{
		Epoch:  epoch,
		Expire: 2,
		Rounds: 5,
	}
	for _, val := range hashTestStrings {
		hash := el.Generate(val)
		if el.Check(hash, val) {
			t.Log("Expired too soon")
			t.Fail()
		}
		time.Sleep(3 * time.Second)
		if el.Check(hash, val) {
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
