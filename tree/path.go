package tree

import (
	"fmt"
	"strconv"
	"strings"
)

// PathSeparator is the symbol used to identify object paths
var PathSeparator = "/"

// Path is how we address items in the Validate.Is struct
type Path string

// NewPath creates a new path from root
func NewPath() Path {
	return Path(PathSeparator)
}

// String returns the string representation of this path. This is really only useful for debugging
func (p Path) String() string {
	return string(p)
}

// IsEqual returns true if the two paths point to the same location, false otherwise
func (p Path) IsEqual(op Path) bool {
	return string(p) == string(op)
}

// Up moves up the tree until reaching root, in which case, it returns the root
func (p Path) Up() Path {
	if p.IsRoot() {
		return p
	}
	if p.IsArrayElement() {
		// Array elements need to go up to the array itself
		index := strings.LastIndex(string(p), "[")
		// index = -1 should be impossible as we cover this case in IsArrayElement
		return Path(string(p)[0:index])
	} else {
		// Split the path:
		parts := strings.Split(string(p), PathSeparator)
		if p.IsAbsolute() {
			// parts[0] will always be an empty string. We know we're not at the root path, so there should be at least 1 item
			return Path(fmt.Sprintf("/%s", strings.Join(parts[1:len(parts)-1], PathSeparator)))
		} else {
			// no forward slash, if nothing left, should be empty string ""
			return Path(strings.Join(parts[0:len(parts)-1], PathSeparator))
		}
	}
}

// DownField goes down the path and references a specific field or a struct. Fields can be leaves or additional nodes
func (p Path) DownField(fieldName string) Path {
	if !isValidFieldName(fieldName) {
		panic(fmt.Errorf("invalid fieldName provided: %s", fieldName))
	}
	if p.IsAbsolute() && p.IsRoot() {
		return Path(fmt.Sprintf("%s%s", string(p), fieldName))
	}
	return Path(fmt.Sprintf("%s%s%s", string(p), PathSeparator, fieldName))
}

// DownIndex goes down the path assuming that the current element is an array
func (p Path) DownIndex(index int) Path {
	return Path(fmt.Sprintf("%s[%d]", string(p), index))
}

// IsRoot returns true if the path is at the root, false if not. The root is defined as being equal to "NewPath", the path references no fields or child objects
// A root Path has no parent
func (p Path) IsRoot() bool {
	return string(p) == PathSeparator
}

// IsAbsolute is true if the path is absolute (starts with the /)
func (p Path) IsAbsolute() bool {
	return strings.HasPrefix(string(p), PathSeparator)
}

// IsArrayElement returns true if the item currently referenced is in an array
func (p Path) IsArrayElement() bool {
	if p.IsRoot() {
		return false
	}
	// get last element
	parts := strings.Split(string(p), PathSeparator)
	lastPart := parts[len(parts)-1]
	return strings.Index(lastPart, "[") != -1
}

// Index returns the index of the current path, or -1 if invalid
func (p Path) Index() int {
	if !p.IsArrayElement() {
		return -1
	}
	// get last element
	parts := strings.Split(string(p), PathSeparator)
	lastPart := parts[len(parts)-1]
	startOfIndex := strings.LastIndex(lastPart, "[") + 2
	endOfIndex := strings.LastIndex(lastPart, "]") + 1
	indexStr := string(p)[startOfIndex:endOfIndex]
	// Error should be impossible, we're setting indexes
	dex, _ := strconv.Atoi(indexStr)
	return dex
}

// FieldName returns the name of the current field, or empty string if not valid (e.g.: we're at the root or in an index)
func (p Path) FieldName() string {
	if p.IsRoot() || p.IsArrayElement() {
		return ""
	}
	// get last element
	parts := strings.Split(string(p), PathSeparator)
	return parts[len(parts)-1]
}

// Depth returns how nested this element is. if IsRoot is true, Depth returns 0. A field at Depth 0 is also zero.
// Examples:
// / = 0
// /field = 0
// /field/field = 1
// /field/field/field = 2
func (p Path) Depth() int {
	if p.IsRoot() {
		return 0
	}
	countRoot := 0
	if p.IsAbsolute() {
		countRoot++
	}
	return strings.Count(string(p), PathSeparator) - countRoot + strings.Count(string(p), "[")
}

// isValidFieldName returns true if the field name provided is valid, false if not
func isValidFieldName(fieldName string) bool {
	forbiddenRunes := "[]" + PathSeparator
	return !strings.ContainsAny(fieldName, forbiddenRunes)
}

// EachComponent iterates through each component and calls the fieldName function if it's a named field component and calls the index function if it's a position in an index
// @return true if this ran to completion, false if methods triggered an early return
func (p Path) EachComponent(fieldName func(fieldName string) bool, index func(index int) bool) bool {
	parts := make([]Path, 0, 2)
	{
		part := p
		for !part.IsRoot() && len(string(part)) != 0 {
			parts = append(parts, part)
			part = part.Up()
		}
	}

	// need to traverse in reverse order
	for i := len(parts) - 1; i >= 0; i-- {
		if parts[i].IsArrayElement() {
			if !index(parts[i].Index()) {
				return false
			}
		} else {
			if !fieldName(parts[i].FieldName()) {
				return false
			}
		}
	}
	return true
}
