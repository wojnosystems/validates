package tree

import (
	"fmt"
	"testing"
)

func TestPath_IsEqual(t *testing.T) {
	p := NewPath()
	q := NewPath()
	if !p.IsEqual(q) {
		t.Error("expected p to be the same as q")
	}

	l1 := p.DownField("puppy")
	l2 := q.DownField("puppy")

	if l1.IsEqual(q) {
		t.Error("expected l1 to be different from p")
	}

	if !l1.IsEqual(l2) {
		t.Error("expected l1 to be the same as l2")
	}
}

func TestPath_DownField(t *testing.T) {
	p := NewPath()
	phone0 := p.DownField("bob").DownField("phones").DownIndex(0)
	if "/bob/phones[0]" != phone0.String() {
		t.Error(`expected path to be: "/bob/phones[0]", but got `, phone0.String())
	}
}

func TestPath_Up(t *testing.T) {
	cases := []struct {
		path     Path
		expected Path
	}{
		{
			path:     NewPath(),
			expected: NewPath(),
		},
		{
			path:     NewPath().DownField("puppy"),
			expected: NewPath(),
		},
		{
			path:     NewPath().DownField("puppy").DownField("puppy"),
			expected: NewPath().DownField("puppy"),
		},
		{
			path:     NewPath().DownField("puppy").DownField("puppy").DownField("puppy"),
			expected: NewPath().DownField("puppy").DownField("puppy"),
		},
		{
			path:     NewPath().DownField("puppy").DownIndex(0),
			expected: NewPath().DownField("puppy"),
		},
		{
			path:     NewPath().DownField("puppy").DownIndex(1).DownIndex(2),
			expected: NewPath().DownField("puppy").DownIndex(1),
		},
	}

	for _, c := range cases {
		actual := c.path.Up()
		if !c.expected.IsEqual(actual) {
			t.Errorf(`Up for: "%s" expected to be: %s, but got: %s`, c.path.String(), c.expected, actual)
		}
	}
}

func TestPath_Depth(t *testing.T) {
	cases := []struct {
		path          Path
		expectedDepth int
	}{
		{
			path:          NewPath(),
			expectedDepth: 0,
		},
		{
			path:          NewPath().DownField("puppy"),
			expectedDepth: 0,
		},
		{
			path:          NewPath().DownField("puppy").DownField("puppy"),
			expectedDepth: 1,
		},
		{
			path:          NewPath().DownField("puppy").DownField("puppy").DownField("puppy"),
			expectedDepth: 2,
		},
		{
			path:          NewPath().DownField("puppy").DownIndex(0),
			expectedDepth: 1,
		},
		{
			path:          NewPath().DownField("puppy").DownIndex(1).DownIndex(2),
			expectedDepth: 2,
		},
	}

	for _, c := range cases {
		actualDepth := c.path.Depth()
		if actualDepth != c.expectedDepth {
			t.Errorf(`depth for: "%s" expected to be: %d, but got: %d`, c.path.String(), c.expectedDepth, actualDepth)
		}
	}
}

func TestPath_EachComponent(t *testing.T) {
	cases := []struct {
		path     Path
		expected []string
	}{
		{
			path:     NewPath(),
			expected: []string{},
		},
		{
			path:     NewPath().DownField("puppy"),
			expected: []string{"puppy"},
		},
		{
			path:     NewPath().DownField("puppy").DownIndex(1),
			expected: []string{"puppy", "[1]"},
		},
		{
			path:     NewPath().DownField("puppy").DownIndex(1).DownIndex(4),
			expected: []string{"puppy", "[1]", "[4]"},
		},
		{
			path:     NewPath().DownField("puppy").DownIndex(1).DownIndex(33).DownField("zoey"),
			expected: []string{"puppy", "[1]", "[33]", "zoey"},
		},
	}

	for _, c := range cases {
		actual := make([]string, 0)
		c.path.EachComponent(func(fieldName string) bool {
			actual = append(actual, fieldName)
			return true
		}, func(index int) bool {
			actual = append(actual, fmt.Sprintf("[%d]", index))
			return true
		})
		if !isStringArrayEqual(actual, c.expected) {
			t.Errorf(`"%s" expected to be: %v, but got: %v`, c.path.String(), c.expected, actual)
		}
	}
}

func isStringArrayEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for ai := range a {
		if a[ai] != b[ai] {
			return false
		}
	}
	return true
}
