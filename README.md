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

** Note on JWT **

It is entirely possible to do the same thing as this
library using JWT. The only minor advantages that this
library has is that it generates links that are a bit
shorter than a JWT token (although that's not likely to
be an issue for most uses) and that it requires a little
less code to be imported (which probably isn't an issue
in 99% of cases either).

As a result, this library serves more as an example on
how to implement than a direct "import and use me". You
could literally replace the HMAC code with JWT code and
end up with something that worked the same way.

## Detail of the problem statement

There are many methods of creating these links. Some
common ones:

(Note: I'm going to use "DB" here to refer to any persistent
storage. I understand fully that some people will assume
that this means a relational database, but it could be
anything that can do persistent storage)


### 1. Store a URL and expiry with each account and provide those when required.
   
Pros:
* Simple and effective

Cons:
* May result in a link that expires before it can be used
  (rare)
* May result in a DB write if a new link needs to be
  generated with a new expirey

### 2. Have a special store for link URLs

Pros:
* No danger of link expiring before it can be used
  
Cons:
* Complexity of another DB storage area
* Extra DB write on each request
* Possible DoS attacks
* Extra work to prune expired URLs from the DB store

### 3. Use encryption to include all needed data in the URL

Pros:
* Really secure, no extra DB access

Cons:
* Requires proper key rotation to be secure

### 4. ExpiringLink

Pros:
* As simple as #1
* Can be implemented such that no additional DB lookups are
  required

Cons:
* Isn't a standard, so won't impress your boss

## Usage

You need to store 1 piece of infomation globally, the
Epoch value. The Epoch can be any date in the past, but
must never change or it will invalidate the system. Values
close to the present will result in smaller hashes.

The expiration and can be changed at any time and
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

## Usage example

ExpiringLink doesn't do everything, it only generates the
link with the embedded expiry. Here are some additional
things that are needed to get the desired result.

1. Add a URL secret to the user data and update it with
a random value any time the password is changed.

2. Use the hash provided by ExpiringLink as _part_ of the
URL. The URL will also have to include something to
identify the user to which the hash applies.

   Examples:

   /resetpassword?userid=5&hash=296c37g5gef7d38828b4aa5df43ef156e86a90b6c2823be37
   /resetpassword/5/296c37g5gef7d38828b4aa5df43ef156e86a90b6c2823be37

   Of course, exactly how you do this depends on you.

3. Store the Epoch for ExpiringLink somewhere that will
never change. It could be a const in your code.

See the example.go file ...