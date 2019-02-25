package issers

import (
	"golang.org/x/text/message"
	"reflect"
	"validates/ifaces"
)

type SimpleValidateError string

func (v SimpleValidateError) ErrorI18n(p *message.Printer) string {
	return p.Sprint(string(v))
}

func (v SimpleValidateError) IsEqual(e ifaces.ValidateError) bool {
	if t, ok := e.(*SimpleValidateError); !ok {
		return false
	} else {
		return string(v) == string(*t)
	}
}

func NewSimpleValidateError(v string) *SimpleValidateError {
	sve := SimpleValidateError(v)
	return &sve
}

var (
	shouldBeTrueMsg = "should be true"
	ShouldBeTrueErr = NewSimpleValidateError(shouldBeTrueMsg)

	shouldBeFalseMsg = "should be false"
	ShouldBeFalseErr = NewSimpleValidateError(shouldBeFalseMsg)

	shouldBePresentMsg = "should be present"
	ShouldBePresentErr = NewSimpleValidateError(shouldBePresentMsg)

	shouldBeIntBetweenMsg            = "should be between %d and %d"
	shouldBeIntGreaterThanMsg        = "should be greater than %d"
	shouldBeIntLessThanMsg           = "should be less than %d"
	shouldBeIntGreaterThanOrEqualMsg = "should be greater than or equal to %d"
	shouldBeIntLessThanOrEqualMsg    = "should be less than or equal to %d"

	shouldBeFloat64BetweenMsg            = "should be between %f and %f"
	shouldBeFloat64GreaterThanMsg        = "should be greater than %f"
	shouldBeFloat64LessThanMsg           = "should be less than %f"
	shouldBeFloat64GreaterThanOrEqualMsg = "should be greater than or equal to %f"
	shouldBeFloat64LessThanOrEqualMsg    = "should be less than or equal to %f"

	shouldBeMatchingRegexpMsg = "should be formatted properly"
	shouldBeEmailMsg          = "should be a valid email address"

	shouldBeInStringSlice = "not an acceptable value"
)

type ShouldBeMsg struct {
	ifaces.ValidateError
	MsgFmt string
	Args   []interface{}
}

// ErrorI18n is the error, but internationalized
// I know English so my errors are all in English
func (v ShouldBeMsg) ErrorI18n(p *message.Printer) string {
	return p.Sprintf(v.MsgFmt, v.Args...)
}

func (v ShouldBeMsg) IsEqual(e ifaces.ValidateError) bool {
	if t, ok := e.(*ShouldBeMsg); !ok {
		return false
	} else {
		return v.MsgFmt == t.MsgFmt &&
			reflect.DeepEqual(v.Args, t.Args)
	}
}

func NewShouldBeIntBetween(low, high int) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeIntBetweenMsg,
		Args:   []interface{}{low, high},
	}
}

func NewShouldBeIntGreaterThan(low int) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeIntGreaterThanMsg,
		Args:   []interface{}{low},
	}
}
func NewShouldBeIntLessThan(high int) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeIntLessThanMsg,
		Args:   []interface{}{high},
	}
}

func NewShouldBeIntGreaterThanOrEqual(low int) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeIntGreaterThanOrEqualMsg,
		Args:   []interface{}{low},
	}
}

func NewShouldBeIntLessThanOrEqual(high int) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeIntLessThanOrEqualMsg,
		Args:   []interface{}{high},
	}
}

func NewShouldBeFloat64Between(low, high float64) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeFloat64BetweenMsg,
		Args:   []interface{}{low, high},
	}
}

func NewShouldBeFloat64GreaterThan(low float64) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeFloat64GreaterThanMsg,
		Args:   []interface{}{low},
	}
}
func NewShouldBeFloat64LessThan(high float64) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeFloat64LessThanMsg,
		Args:   []interface{}{high},
	}
}

func NewShouldBeFloat64GreaterThanOrEqual(low float64) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeFloat64GreaterThanOrEqualMsg,
		Args:   []interface{}{low},
	}
}

func NewShouldBeFloat64LessThanOrEqual(high float64) *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeFloat64LessThanOrEqualMsg,
		Args:   []interface{}{high},
	}
}

func NewShouldMatchingRegexp() *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeMatchingRegexpMsg,
		Args:   []interface{}{},
	}
}
func NewShouldBeInStringSlice() *ShouldBeMsg {
	return &ShouldBeMsg{
		MsgFmt: shouldBeInStringSlice,
		Args:   []interface{}{},
	}
}
