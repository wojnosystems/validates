package issers

import (
	"fmt"
	"github.com/wojnosystems/validates/ifaces"
	"github.com/wojnosystems/validates/tree"
	"net/url"
	"regexp"
	"strings"
)

// Validater describes how structures should behave
// if they want to be validatable
type Validater interface {
	// Validate performs the validation on an object
	Validate(is *Is) (*Is, error)
}

// Is is where the validates.Is validation errors are stored
// Is only writes, you cannot "undo" an error written to it.
// You can only keep adding more errors. A structure is valid
// if no ValidationErrors are reported. This is very much like
// Rail's Validations objects
type Is struct {
	// currentPath is the textual representation of our
	// current location in the structure. Errors are tied
	// to this location when the validation methods are called
	currentPath tree.Path

	// errorsRoot is the ever-present errorsRoot of our
	// error tree. This is where the errors are stored
	errorsRoot *tree.ErrorNode

	// errorsCount is a summation of all of the errors
	errorsCount int
}

// NewRoot creates a new Is with the current path as the Root (/)
func NewRoot() *Is {
	return &Is{
		currentPath: tree.NewPath(),
	}
}

// CurrentPath is the path we're currently set to
// @return the current path (read-only)
func (i Is) CurrentPath() tree.Path {
	return i.currentPath
}

// HasErrors returns true if there is at least 1 error
func (i Is) HasErrors() bool {
	return i.errorsCount != 0
}

// Errors returns the errorsRoot of the errors
func (i Is) Errors() *tree.ErrorNode {
	return i.errors()
}

// errors is the lazy-loading way to get the root of this tree
// @return a new ErrorNode, attached to receiver, or the already created ErrorNode
func (i *Is) errors() *tree.ErrorNode {
	if i.errorsRoot == nil {
		i.errorsRoot = tree.NewErrorNode(nil)
	}
	return i.errorsRoot
}

// Len returns the number of errors
func (i Is) Len() int {
	return i.errorsCount
}

// ValidStructField performs validation on a nested struct field
// @param fieldName of this structure. If you have a nested struct
//   with json struct tag with name "thing" then "thing" should be
//   this value
// @param validator is a structure to recurse into and test for validation errors
// @return err if there was a problem validating input for some
//   reason that cause validation to stop prematurely
func (i *Is) ValidStructField(fieldName string, validator Validater) (err error) {
	i.WithField(fieldName, func(is *Is) {
		_, err = validator.Validate(i)
	})
	return err
}

// ValidStructIndex performs validation on a nested struct field
// @param index of this structure. If you have a nested struct
//   in an array at index 5: "thing[5]" then 5 should be this value
// @param validator is a structure to recurse into and test for validation errors
// @return err if there was a problem validating input for some
//   reason that cause validation to stop prematurely
func (i *Is) ValidStructIndex(index int, validator Validater) (err error) {
	i.WithIndex(index, func(is *Is) {
		_, err = validator.Validate(i)
	})
	return err
}

// ValidEachStruct performs validation on a nested struct at each index
// This is a convenience method for WithField + ValidStructIndex
// @return err an error that caused validation to stop prematurely
//   if any struct returns this value, no further validation will be
//   performed and the error will be returned along with any validations
//   performed at the time the error was returned
func (i *Is) ValidEachStruct(fieldName string, values []Validater) (err error) {
	i.WithField(fieldName, func(is *Is) {
		for idx, value := range values {
			err = is.ValidStructIndex(idx, value)
			if err != nil {
				break
			}
		}
	})

	return err
}

// WithField is a convenience method to group fields together
// it's called with a function context because when that
// function completes, the current path in receiver is
// reverted back to what it was before starting WithField
func (i *Is) WithField(fieldName string, wrap func(is *Is)) {
	i.With(func() {
		i.currentPath = i.currentPath.DownField(fieldName)
		wrap(i)
	})
}

// WithIndex is a convenience method to group fields together
// it's called with a function context because when that
// function completes, the current path in receiver is
// reverted back to what it was before starting WithIndex
func (i *Is) WithIndex(index int, wrap func(is *Is)) {
	i.With(func() {
		i.currentPath = i.currentPath.DownIndex(index)
		wrap(i)
	})
}

// With is a convenience method to group fields together
// it's called with a function context because when that
// function completes, the current path in receiver is
// reverted back to what it was before starting With
func (i *Is) With(wrap func()) {
	originalPath := i.currentPath
	defer func() {
		i.currentPath = originalPath
	}()
	wrap()
}

// currentErrorNode Builds and/or navigates to the current error node
// Nodes are only created as currentErrorNode is used
func (i *Is) currentErrorNode() (n *tree.ErrorNode) {
	n = i.errors()
	i.currentPath.EachComponent(func(fieldName string) bool {
		n = n.DownField(fieldName)
		return true
	}, func(index int) bool {
		n = n.DownIndex(index)
		return true
	})
	return
}

// Invalid is the input error assertion that states that some
// input validation has failed and the structure that this
// validator describes is not valid. All other assertions can
// be built upon this.
// @param msg is the message to use. There is no default message
//   for Invalid
func (i *Is) Invalid(msg ifaces.ValidateError) {
	i.errorsCount++
	// Only creates the chain if we have an error
	// We do not want to pre-allocate memory unless we know we're going to use it
	n := i.currentErrorNode()
	n.Add(msg)
}

// True ensures that value is true, otherwise it appends an error.
// All other assertions can be built upon this.
// Given the value, if the value is not true, the error message: msg
// will be generated using the function call.
// If nil is returned from the method, or the method itself is nil,
// the default ShouldBeTrueErr object is returned
// @param value does nothing if this value is true, if false, will
//   append the error with the message
// @param msg is a callback used to generate the ifaces.ValidateError.
//   If nil, the default will be called. This allows messages to be overwritten
// @return true if valid (no errors added) false if not
func (i *Is) True(value bool, msg func() ifaces.ValidateError) bool {
	if !value {
		i.Invalid(msgOrDefault(msg, ShouldBeTrueErr))
		return false
	}
	return true
}

// False ensures that the value is true and appends an error if that is not the case
// @param value is the value to evaluate
// @msg is the callback used to generate the message. Leave nil or return nil
//   to use the default: ShouldBeFalseErr
// @return true if valid (no errors added) false if not
func (i *Is) False(value bool, msg func() ifaces.ValidateError) bool {
	return i.True(!value, func() ifaces.ValidateError {
		return msgOrDefault(msg, ShouldBeFalseErr)
	})
}

// msgOrDefault helps resolve the error message to use.
// if the in function is nil or returns nil, def (default) will be used instead
func msgOrDefault(in func() ifaces.ValidateError, def ifaces.ValidateError) ifaces.ValidateError {
	if in != nil {
		ret := in()
		if ret == nil {
			return def
		}
		return ret
	}
	return def
}

// Required is the gateway for performing the requirement test. If isPresent returns true, no error will be added and this method will return true
// If isPresent is false, the error will be recorded on the current field and the function will return false. Here's how it's intended to be used:
//
// is := NewIs()
// is.WithField("myName", func(is *Is) {
//   if is.Required( validates.NotEmptyString(myName) ) {
//     is.StringLengthBetween( myName, 1, 15 )
//   }
// } )
//
// If the value could be empty, Required is assumed to have been called by the user of the library (that's probably you reading this now) as per the above.
// Calling methods that require the object to exist will attempt to use the value as normal.
// When using Required, it returns false if the value is missing. This allows you to skip validates if they don't make sense
// @return true if valid (no errors added) false if not
func (i *Is) Required(isPresent bool) bool {
	if !isPresent {
		i.True(false, func() ifaces.ValidateError {
			return ShouldBePresentErr
		})
		return false
	}
	return true
}

// IntBetween creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) IntBetween(value, low, high int, msg func() ifaces.ValidateError) bool {
	if low > high {
		panic("low cannot be greater than high")
	}
	return i.True(low <= value && value <= high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeIntBetween(low, high))
	})
}

// IntGreaterThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) IntGreaterThan(value, low int, msg func() ifaces.ValidateError) bool {
	return i.True(low < value, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeIntGreaterThan(low))
	})
}

// IntLessThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) IntLessThan(value, high int, msg func() ifaces.ValidateError) bool {
	return i.True(value < high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeIntLessThan(high))
	})
}

// IntGreaterThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) IntGreaterThanOrEqual(value, low int, msg func() ifaces.ValidateError) bool {
	return i.True(low <= value, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeIntGreaterThanOrEqual(low))
	})
}

// IntLessThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) IntLessThanOrEqual(value, high int, msg func() ifaces.ValidateError) bool {
	return i.True(value <= high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeIntLessThanOrEqual(high))
	})
}

// Float64Between creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) Float64Between(value, low, high float64, msg func() ifaces.ValidateError) bool {
	if low > high {
		panic("low cannot be greater than high")
	}
	return i.True(low <= value && value <= high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeFloat64Between(low, high))
	})
}

// Float64GreaterThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) Float64GreaterThan(value, low float64, msg func() ifaces.ValidateError) bool {
	return i.True(low < value, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeFloat64GreaterThan(low))
	})
}

// Float64LessThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) Float64LessThan(value, high float64, msg func() ifaces.ValidateError) bool {
	return i.True(value < high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeFloat64LessThan(high))
	})
}

// Float64GreaterThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) Float64GreaterThanOrEqual(value, low float64, msg func() ifaces.ValidateError) bool {
	return i.True(low <= value, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeFloat64GreaterThanOrEqual(low))
	})
}

// Float64LessThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) Float64LessThanOrEqual(value, high float64, msg func() ifaces.ValidateError) bool {
	return i.True(value <= high, func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldBeFloat64LessThanOrEqual(high))
	})
}

// StringLengthBetween creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) StringLengthBetween(value string, low, high int, msg func() ifaces.ValidateError) bool {
	if low > high {
		panic("low cannot be greater than high")
	}
	return i.IntBetween(len(value), low, high, func() ifaces.ValidateError {
		defMsg := NewShouldBeIntBetween(low, high)
		defMsg.MsgFmt = fmt.Sprintf("length %s", defMsg.MsgFmt)
		return msgOrDefault(msg, defMsg)
	})
}

// StringLengthGreaterThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) StringLengthGreaterThan(value string, low int, msg func() ifaces.ValidateError) bool {
	return i.IntGreaterThan(len(value), low, func() ifaces.ValidateError {
		defMsg := NewShouldBeIntGreaterThan(low)
		defMsg.MsgFmt = fmt.Sprintf("length %s", defMsg.MsgFmt)
		return msgOrDefault(msg, defMsg)
	})
}

// StringLengthLessThan creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) StringLengthLessThan(value string, high int, msg func() ifaces.ValidateError) bool {
	return i.IntLessThan(len(value), high, func() ifaces.ValidateError {
		defMsg := NewShouldBeIntGreaterThan(high)
		defMsg.MsgFmt = fmt.Sprintf("length %s", defMsg.MsgFmt)
		return msgOrDefault(msg, defMsg)
	})
}

// StringLengthGreaterThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) StringLengthGreaterThanOrEqual(value string, low int, msg func() ifaces.ValidateError) bool {
	return i.IntGreaterThanOrEqual(len(value), low, func() ifaces.ValidateError {
		defMsg := NewShouldBeIntGreaterThanOrEqual(low)
		defMsg.MsgFmt = fmt.Sprintf("length %s", defMsg.MsgFmt)
		return msgOrDefault(msg, defMsg)
	})
}

// StringLengthLessThanOrEqual creates an error unless string's length is between the provided values (inclusive)
// @return true if valid (no errors added) false if not
func (i *Is) StringLengthLessThanOrEqual(value string, high int, msg func() ifaces.ValidateError) bool {
	return i.IntLessThanOrEqual(len(value), high, func() ifaces.ValidateError {
		defMsg := NewShouldBeIntGreaterThanOrEqual(high)
		defMsg.MsgFmt = fmt.Sprintf("length %s", defMsg.MsgFmt)
		return msgOrDefault(msg, defMsg)
	})
}

// StringNotEmpty creates an error unless string's length non-zero
// @return true if valid (no errors added) false if not
func (i *Is) StringNotEmpty(value string, msg func() ifaces.ValidateError) bool {
	return i.True(len(value) != 0, func() ifaces.ValidateError {
		return NewShouldBeNotEmpty()
	})
}

// StringInStringSlice creates an error unless the value exists in the values array
// @return true if valid (no errors added) false if not
func (i *Is) StringInStringSlice(value string, values []string, msg func() ifaces.ValidateError) bool {
	for _, v := range values {
		if 0 == strings.Compare(value, v) {
			return true
		}
	}
	i.Invalid(msgOrDefault(msg, NewShouldBeInStringSlice()))
	return false
}

// MatchingRegexp creates a ValidationError unless the value matches the regular expression provided
// @return true if valid (no errors added) false if not
func (i *Is) MatchingRegexp(value string, reg *regexp.Regexp, msg func() ifaces.ValidateError) bool {
	return i.True(reg.MatchString(value), func() ifaces.ValidateError {
		return msgOrDefault(msg, NewShouldMatchingRegexp())
	})
}

// EmailAddress creates a ValidationError unless the value matches the email validation
// regular expression: `^[^@]+@.+\.[^.]{2,}$`
// @return true if valid (no errors added) false if not
func (i *Is) EmailAddress(value string, msg func() ifaces.ValidateError) bool {
	return i.MatchingRegexp(value, emailRegexpCompiled, func() ifaces.ValidateError {
		defMsg := NewShouldMatchingRegexp()
		defMsg.MsgFmt = shouldBeEmailMsg
		return msgOrDefault(msg, defMsg)
	})
}

// URI creates a ValidationError unless the value is a url as per Go's url.ParseRequestURI method
// @return true if no error (it was a URL), false if error
// @return true if valid (no errors added) false if not
func (i *Is) URI(value string, msg func() ifaces.ValidateError) bool {
	if i.StringNotEmpty(value, msg) {
		_, err := url.ParseRequestURI(value)
		if err != nil {
			i.Invalid(msgOrDefault(msg, NewShouldBeURL(err.(*url.Error))))
			return false
		}
		return true
	}
	return false
}

var (
	isEmailRegexp = `^[^@]+@.+\.[^.]{2,}$`
)
var (
	emailRegexpCompiled *regexp.Regexp
)

func init() {
	emailRegexpCompiled = regexp.MustCompile(isEmailRegexp)
}
