package validates

import "testing"

func TestNotEmptyString(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{
			input: "",
		},
		{
			input:    "test",
			expected: true,
		},
	}

	for _, c := range cases {
		actual := NotEmptyString(c.input)
		if c.expected {
			if !actual {
				t.Error("expected a not empty string")
			}
		} else {
			if actual {
				t.Error("expected an empty string")
			}
		}
	}
}
