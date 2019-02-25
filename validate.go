package validates

import (
	"validates/issers"
)

// The Validates package provides a way to uniformly validate structs and
// their components/fields.
//
// I've opted against using struct tags because they are not extensible or
// programmatic. They also make internationalization impossible. Far better
// to have compose-able validations than a set of meta code that cannot be
// altered or extended except by the author. In Go, meta-programming using
// reflect is a hack and a dangerous, time-consuming one at that. Having
// written several reflect-based libraries, if you ever feel as though you
// must do it that way, stop and reconsider. reflect should always be a
// last resort.

// On performs the validation on a root struct. It's a convenience method
// to kick off validation
//
// @param on the struct to start performing validations on
// @return i the validation result containing all of the validation errors
// @return err any errors encountered while processing the validations.
//   Only returned if some abnormal condition prevented validation
//   from occurring (like a missing configuration or service), validations
//   are NOT stored here, this is only to document abnormal validation
//   conditions, not for bad inputs.
//
// @example
// ```go
// if i, err := validates.On(myStruct); err != nil {
//   log.Fatal(err)
// } else {
//	 if i.HasErrors() {
//     // do format your response
//   }
// }
// ```
func On(on issers.Validater) (i *issers.Is, err error) {
	return on.Validate(issers.NewRoot())
}

// NotEmptyString returns true if the string is not empty
func NotEmptyString(s string) bool {
	return len(s) != 0
}
