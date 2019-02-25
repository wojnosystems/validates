package issers

import (
	"github.com/wojnosystems/validates/ifaces"
	"github.com/wojnosystems/validates/tree"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"testing"
)

var defTestMessagePrinter = message.NewPrinter(language.AmericanEnglish)

func TestIs_True(t *testing.T) {
	i := &Is{}
	i.WithField("zoey", func(is *Is) {
		is.True(true, nil)
	})
	if i.Len() != 0 || i.HasErrors() {
		t.Error("expected no errors")
	}

	i.WithField("zoey", func(is *Is) {
		is.True(false, nil)
	})
	if i.Len() != 1 || !i.HasErrors() {
		t.Error("expected errors")
	}

	zoeyErrors := i.Errors().NamedChildren["zoey"]
	if !zoeyErrors.HasErrors() {
		t.Error("expected errors")
	}
	if zoeyErrors.Errors()[0] != ShouldBeTrueErr {
		t.Errorf(`expected: "%s" but got: "%s"`, zoeyErrors.Errors()[0].ErrorI18n(defTestMessagePrinter), ShouldBeTrueErr.ErrorI18n(defTestMessagePrinter))
	}

	i.WithField("slater", func(is *Is) {
		is.False(true, nil)
	})
	slaterErrors := i.Errors().NamedChildren["slater"]
	if slaterErrors.Errors()[0] != ShouldBeFalseErr {
		t.Errorf(`expected: "%s" but got: "%s"`, slaterErrors.Errors()[0].ErrorI18n(defTestMessagePrinter), ShouldBeFalseErr.ErrorI18n(defTestMessagePrinter))
	}
}

func TestIs_StringLengthBetween(t *testing.T) {
	cases := map[string]struct {
		value     string
		low, high int
		expect    bool
	}{
		"zoey good": {
			value:  "zoey",
			low:    1,
			high:   4,
			expect: true,
		},
		"zoey short": {
			value: "zoey",
			low:   5,
			high:  12,
		},
		"zoey long": {
			value: "zoey",
			low:   1,
			high:  2,
		},
	}

	for caseName, c := range cases {
		is := &Is{}
		ret := is.StringLengthBetween(c.value, c.low, c.high, nil)
		if c.expect != ret {
			t.Errorf("%s: expected %s but got %s", caseName, boolToString(c.expect), boolToString(ret))
		}
	}
}

func boolToString(b bool) string {
	if b {
		return "t"
	} else {
		return "f"
	}
}

func TestIs_StringLengthGreaterThan(t *testing.T) {
	cases := map[string]struct {
		value  string
		bound  int
		expect bool
	}{
		"zoey good": {
			value:  "zoey",
			bound:  1,
			expect: true,
		},
		"zoey short": {
			value: "zoey",
			bound: 4,
		},
	}

	for caseName, c := range cases {
		is := &Is{}
		ret := is.StringLengthGreaterThan(c.value, c.bound, nil)
		if c.expect != ret {
			t.Errorf("%s: expected %s but got %s", caseName, boolToString(c.expect), boolToString(ret))
		}
	}
}

func TestIs_StringLengthLessThan(t *testing.T) {
	cases := map[string]struct {
		value  string
		bound  int
		expect bool
	}{
		"zoey good": {
			value:  "zoey",
			bound:  5,
			expect: true,
		},
		"zoey short": {
			value: "zoey",
			bound: 2,
		},
	}

	for caseName, c := range cases {
		is := &Is{}
		ret := is.StringLengthLessThan(c.value, c.bound, nil)
		if c.expect != ret {
			t.Errorf("%s: expected %s but got %s", caseName, boolToString(c.expect), boolToString(ret))
		}
	}
}

func TestIs_StringLengthGreaterThanOrEqual(t *testing.T) {
	cases := map[string]struct {
		value  string
		bound  int
		expect bool
	}{
		"zoey good": {
			value:  "zoey",
			bound:  4,
			expect: true,
		},
		"zoey bad": {
			value: "zoey",
			bound: 5,
		},
	}

	for caseName, c := range cases {
		is := &Is{}
		ret := is.StringLengthGreaterThanOrEqual(c.value, c.bound, nil)
		if c.expect != ret {
			t.Errorf("%s: expected %s but got %s", caseName, boolToString(c.expect), boolToString(ret))
		}
	}
}

func TestIs_StringLengthLessThanOrEqual(t *testing.T) {
	cases := map[string]struct {
		value  string
		bound  int
		expect bool
	}{
		"zoey good": {
			value:  "zoey",
			bound:  4,
			expect: true,
		},
		"zoey bad": {
			value: "zoey",
			bound: 2,
		},
	}

	for caseName, c := range cases {
		is := &Is{}
		ret := is.StringLengthLessThanOrEqual(c.value, c.bound, nil)
		if c.expect != ret {
			t.Errorf("%s: expected %s but got %s", caseName, boolToString(c.expect), boolToString(ret))
		}
	}
}

type testName struct {
	First  string `json:"first"`
	Last   string `json:"last"`
	Middle string `json:"middle"`
}

func (r testName) Validate(is *Is) (*Is, error) {
	is.WithField("first", func(is *Is) {
		value := r.First
		if is.Required(len(value) != 0) {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	is.WithField("last", func(is *Is) {
		value := r.Last
		if is.Required(len(value) != 0) {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	is.WithField("middle", func(is *Is) {
		value := r.Middle
		if len(value) != 0 {
			is.StringLengthBetween(value, 1, 32, nil)
		}
	})
	return is, nil
}

type testRoot struct {
	Name   testName `json:"name"`
	Age    int      `json:"age"`
	Emails []string `json:"emails"`
}

func (r testRoot) Validate(is *Is) (*Is, error) {
	err := is.ValidStructField("name", &r.Name)
	if err != nil {
		return is, err
	}
	is.WithField("age", func(is *Is) {
		value := r.Age
		is.IntGreaterThanOrEqual(value, 18, nil)
	})
	is.WithField("emails", func(is *Is) {
		value := r.Emails
		is.IntGreaterThanOrEqual(len(value), 1, func() ifaces.ValidateError {
			e := NewShouldBeIntGreaterThanOrEqual(1)
			// custom message
			e.MsgFmt = "requires more than %d"
			return e
		})
		for i, email := range r.Emails {
			is.WithIndex(i, func(is *Is) {
				is.EmailAddress(email, nil)
			})
		}
	})
	return is, err
}

func goldenTestRoot() *testRoot {
	return &testRoot{
		Name: testName{
			First:  "chris",
			Last:   "wojno",
			Middle: "r",
		},
		Age: 30,
		Emails: []string{
			"clearlyFake@wojno.com",
		},
	}
}

func TestIs_ValidStructField(t *testing.T) {
	cases := map[string]struct {
		input          testRoot
		expectedErrors *tree.ErrorNode
	}{
		"working": {
			input:          *goldenTestRoot(),
			expectedErrors: tree.NewErrorNode(nil),
		},
		"missing required value": {
			input: func() testRoot {
				g := *goldenTestRoot()
				g.Name.First = ""
				return g
			}(),
			expectedErrors: func() *tree.ErrorNode {
				e := tree.NewErrorNode(nil)
				e.DownField("name").DownField("first").Add(ShouldBePresentErr)
				return e
			}(),
		},
		"bad age": {
			input: func() testRoot {
				g := *goldenTestRoot()
				g.Age = 15
				return g
			}(),
			expectedErrors: func() *tree.ErrorNode {
				e := tree.NewErrorNode(nil)
				e.DownField("age").Add(NewShouldBeIntGreaterThanOrEqual(18))
				return e
			}(),
		},
	}

	for caseName, c := range cases {
		is, err := c.input.Validate(&Is{})
		if err != nil {
			t.Error("not expecting an error")
		}
		if !c.expectedErrors.IsEqual(is.Errors()) {
			t.Errorf("%s: errors were not the same, expected: %v, got %v", caseName, *c.expectedErrors, *is.Errors())
		}
	}
}

func TestIs_URL(t *testing.T) {
	cases := map[string]struct {
		input    string
		expected bool
	}{
		"empty": {
			input:    "",
			expected: false,
		},
		"ok url": {
			input:    "https://www.wojno.com",
			expected: true,
		},
		"bad": {
			input:    "puppy",
			expected: false,
		},
	}

	for caseName, c := range cases {
		is := NewRoot()
		actual := is.URI(c.input, nil)

		if c.expected {
			if !actual {
				t.Errorf("%s: expected url", caseName)
			}
		} else {
			if actual {
				t.Errorf("%s: expected not url", caseName)
			}
		}
	}
}
