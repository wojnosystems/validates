# Overview

The Validates package provides a way to uniformly validate structs and their components/fields in a Rails-like way. I've always loved the Validates class in the Rails package for model and input validation.

I've opted against using tags because they are not extensible or programmatic. Far better to have compose-able validations than a set of meta code that cannot be altered or extended except by the author. In Go, meta programming using reflect is a hack and a dangerous, time-consuming one at that. Having written several reflect-based libraries, if you ever feel as though you must do it that way, stop and reconsider. reflect should always be a last resort.

# Requirements

Goal: make validating inputs easy.

A uniform method of collecting validation errors over structs. 

 * Uniform
 * Support Nested structs
 * Support Embedded structs
 * Single point to know if IsInvalid
 * Re-enterable: running the validations multiple times should return the same results if the data has not changed between runs
 
# Components

 * ValidationError: these should not be go-errors as go-errors are usually reserved for fatal errors. There are some exceptions as with the strconv.ParseInt returns an Error implementer.
 * Field: each field is a component of a struct. Since struct can be embedded or nested, each field has an address, local to the root of the validation
 * Embedded structs are fields. You cannot have a primitive value and an embedded struct with the same name such that a path would conflict.
 
 # Examples
 
 ```go
package main

import "github.com/wojnosystems/validates/issers"
import "github.com/wojnosystems/validates/ifaces"
import "github.com/wojnosystems/validates/tree"
import "json"

// define a nested struct
type testName struct {
	First  string `json:"first"`
	Last   string `json:"last"`
	Middle string `json:"middle"`
}

// validate it
func (r testName) Validate(is *issers.Is) (*issers.Is, error) {
	is.WithField("first", func(is *issers.Is) {
		value := r.First
		if is.Required(issers.NotEmptyString(value)) {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	is.WithField("last", func(is *issers.Is) {
		value := r.Last
		if is.Required(issers.NotEmptyString(value)) {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	is.WithField("middle", func(is *issers.Is) {
		value := r.Middle
		if issers.NotEmptyString(value) {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	return is, nil
}

// define a root struct
type testRoot struct {
	Name testName `json:"name"`
	Age  int      `json:"age"`
	Emails []string `json:"emails"`
}

// validate it
func (r testRoot) Validate(is *issers.Is) (*issers.Is, error) {
	err := is.ValidStructField( "name", &r.Name)
	if err != nil {
		return is, err
	}
	is.WithField("age", func(is *issers.Is) {
		value := r.Age
		is.IntGreaterThanOrEqual(value, 18, nil)
	})
	is.WithField("emails", func(is *issers.Is) {
		value := r.Emails
		is.IntGreaterThanOrEqual(len(value), 1, func() ifaces.ValidateError {
			e := tree.NewShouldBeIntGreaterThanOrEqual(1)
			// custom message
			e.MsgFmt = "requires more than %d"
			return e
		})
		for i, email := range r.Emails {
			is.WithIndex(i, func(is *issers.Is) {
				is.EmailAddress(email, nil)
			})
		}
	})
	return is, err
}

func main() {
	var data testRoot
	err := json.Unmarshal(byte[](`{"name":{"first":"chris","last":"wojno"}, "age": 18, "emails": ["faketest@wojno.com"]}`), &data )
	
	vErr, err := validates.On(&data)
    if err != nil {
	    log.Fatal(err)
	}
	// vErr are your errors from the validation
	if vErr.HasErrors() {
		// return errors to user to fix
	}
}
```

# Copyright

Copyright © 2019 Chris Wojno. All rights reserved.

No Warranties. Use this software at your own risk.

# License

[Creative Commons: Attribution-NonCommercial-ShareAlike 4.0 International](http://creativecommons.org/licenses/by-nc-sa/4.0/)