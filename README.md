ExpiringLink
============
This is a simple library for creating unique strings that
have a built in expiration.

The target use case is web links for password resets that
expire after some amount of time.

The typical approach to this problem is to generate a
random code for each request and keep track of when each
code expires. Theoretically, such a system could be
DoSed by generating password resets as fast as possible
and overwhelming whatever system is keeping track of
the expirations.

That being said, I'm not aware of any attack like that
ever happening, and there are multiple ways to defend
against it.

However, I believe I've got a simpler solution. Using this
library, you only need to ever generate a single random
value per user, no matter how many password resets they
request. The unique codes are self-expiring, so the
library can easily detect when a code has expired.

This system is **theoretical** at this time. I'm waiting to
hear if people much smarter at math than me can prove that
it's not secure before calling it anything other than that.

## Usage

You need to store 1 piece of infomation globally, the
Epoch value. The Epoch can be any date in the past, but
must never change or it will invalidate the system. Values
close to the present will result in smaller hashes.

The expiration and rounds can be changed at any time and
will only affect tokens generated after the change
(previously generated tokens will still validate correctly)

In the target use case, each user would have a random value
generated each time their password is changed (note, for
security, this should be **random** and not tied to any
other value) That random value is fed into the Generate()
function creating a unique string with a builtin
expiration. The user can generate as many pasword reset
links as they like. Each one will be unique with a built
in expiration. Once the password is successfully reset the
new random value will immediately invalidate all previous
password reset links.

See the test file for examples.
