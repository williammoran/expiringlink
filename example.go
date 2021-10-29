package expiringlink

import (
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// This is an example file to illustrate how to use the
// library. It's illustrative, but not executable.

var pwResetEpoch time.Time

var expiringLink *ExpiringLink

func init() {
	layout := "2006-01-02T15:04:05.000Z"
	str := "2021-09-29T12:45:26.371Z"
	t, err := time.Parse(layout, str)
	if err != nil {
		panic(err)
	}
	pwResetEpoch = t
	expiringLink = &ExpiringLink{
		Epoch:     pwResetEpoch,
		Expire:    24 * time.Hour,
		Rounds:    8,
		MaxRounds: 16,
	}
}

// account stores user information
type account struct {
	ID        int
	Password  string
	URLSecret string
	// Would obviusly contain more fields ...
}

func (a *account) setPassword(n string) {
	// Password hashing omitted. It's necessary for a
	// real implementation, but just imagine that it's
	// here for this example. Note, probably not the
	// best way to generate random data, but hopefully
	// gets the point across.
	a.URLSecret = strconv.Itoa(rand.Int())
}

// accountStorage would have an implementation somewhere
// that talked to a database or an auth provider or whatever
type accountStorage interface {
	Account(int) *account
}

// generatePWresetLink returns a string to use as a
// clickable link
func generatePWresetLink(a *account) string {
	hash := expiringLink.Generate(a.URLSecret)
	return fmt.Sprintf(
		"http//example.com/resetpassword?id=%d&confirm=%s",
		a.ID, hash,
	)
}

// verifyPWResetLink returns true if the clicked link is
// valid along with the account of the user it's valid for.
// Obviously, the http handler setup would have to handle
// the link in generatePWresetLink()
func verifyPWResetLink(as accountStorage, r *http.Request) (*account, bool) {
	// Much error checking omitted for clarity
	idString := r.URL.Query()["id"][0]
	id, _ := strconv.Atoi(idString)
	acct := as.Account(id)
	confirmation := r.URL.Query()["confirm"][0]
	return acct, expiringLink.Check(confirmation, acct.URLSecret) == nil
}
