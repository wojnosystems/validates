package tree

import (
	"container/list"
	"validates/ifaces"
)

// ErrorNode contains the errors for this node
// Storage is allocated lazily to avoid pre-allocation
type ErrorNode struct {
	// parent is the reference to the parent of this node
	parent *ErrorNode
	// NamedChildren are the fields of this node by name
	NamedChildren map[string]*ErrorNode
	// NumberedChildren are the index positions (sparse) of this node, if any
	NumberedChildren map[int]*ErrorNode
	// errs is the list of errors for this node. If there are errs, there should not be any children
	errs []ifaces.ValidateError
}

// NewErrorNode creates a new node in the tree using the
// parent as the root. Pass in nil if there is no parent
// or this is the root of the tree. Internal components
// are lazily created to conserve memory
func NewErrorNode(parent *ErrorNode) *ErrorNode {
	return &ErrorNode{
		parent:           parent,
		NamedChildren:    nil,
		NumberedChildren: nil,
		errs:             nil,
	}
}

// Up traverses to the parent of this node, if not root
func (n *ErrorNode) Up() *ErrorNode {
	if n.IsRoot() {
		return n
	}
	return n.parent
}

// IsRoot returns true if there is no parent to this node, false if there is
func (n ErrorNode) IsRoot() bool {
	return n.parent == nil
}

// DownField creates a new node or traverses into the node if missing and returns it
func (n *ErrorNode) DownField(name string) (e *ErrorNode) {
	var ok bool
	if n.NamedChildren == nil {
		n.NamedChildren = make(map[string]*ErrorNode)
	}
	if e, ok = n.NamedChildren[name]; !ok {
		e = NewErrorNode(n)
		n.NamedChildren[name] = e
	}
	return e
}

// DownIndex creates a new node or traverses into the node if missing and returns it
func (n *ErrorNode) DownIndex(index int) (e *ErrorNode) {
	var ok bool
	if n.NumberedChildren == nil {
		n.NumberedChildren = make(map[int]*ErrorNode)
	}
	if e, ok = n.NumberedChildren[index]; !ok {
		e = NewErrorNode(n)
		n.NumberedChildren[index] = e
	}
	return e
}

// HasErrors is true if there is at least 1 error in itself OR its children
func (n ErrorNode) HasErrors() bool {
	if n.errs != nil && len(n.errs) != 0 {
		return true
	}
	if n.NamedChildren != nil {
		for _, c := range n.NamedChildren {
			if c.HasErrors() {
				return true
			}
		}
	}
	if n.NumberedChildren != nil {
		for _, c := range n.NumberedChildren {
			if c.HasErrors() {
				return true
			}
		}
	}
	return false
}

// Errors returns the list of errors for ONLY this node. Does not recurse to children
func (n ErrorNode) Errors() []ifaces.ValidateError {
	return n.errs
}

// HasErrorAt traverses down the tree, looking to see if there's an error in the structure at that path
// If the path is missing, returns false immediately. Does not detect errors on child objects
// @return true if there's at least 1 error at the path, false if not
func (n *ErrorNode) HasErrorAt(path Path) bool {
	current := n.traverseTo(path)
	if current == nil {
		return false
	}
	return len(current.errs) != 0
}

// traverseTo descends the tree from the current node given the path
// If the path is missing at any point along the way, nil is returned, otherwise, return the node at the end of the path
// Does not create nodes. Does not alter the nodes in anyway.
func (n *ErrorNode) traverseTo(path Path) *ErrorNode {
	current := n
	if path.EachComponent(func(fieldName string) bool {
		if c, ok := current.NamedChildren[fieldName]; !ok {
			return false
		} else {
			current = c
			return true
		}
	}, func(index int) bool {
		if c, ok := current.NumberedChildren[index]; !ok {
			return false
		} else {
			current = c
			return true
		}
	}) {
		return current
	}
	return nil
}

// IsErrorAt returns true if the validateError provided is present at the path node, false if not or path node doesn't exist
func (n ErrorNode) IsErrorAt(path Path, validateError ifaces.ValidateError) bool {
	current := n.traverseTo(path)
	if current == nil {
		return false
	}
	for _, e := range current.errs {
		if e.IsEqual(validateError) {
			return true
		}
	}
	return false
}

// Add appends the error to this node
func (n *ErrorNode) Add(e ifaces.ValidateError) {
	if n.errs == nil {
		n.errs = make([]ifaces.ValidateError, 0, 0)
	}
	n.errs = append(n.errs, e)
}

// IsEqual attempts to ensure that the current node has the same errors and sub-errors as the provided node.
// IsEqual will attempt to compare itself first before recursing into child nodes
func (n *ErrorNode) IsEqual(o *ErrorNode) bool {
	// First, check the errors locally
	if n.errs != nil && o.errs == nil || n.errs == nil && o.errs != nil {
		return false // one had errors, but the other did not
	}

	// Compare this node's errors
	// both either have errors or do not have errors at this point
	if n.errs != nil {
		// different lengths
		if len(n.errs) != len(o.errs) {
			return false
		}

		// both have errors and both are the same length
		// Create a list of errors to "mark" them as no yet compared
		unVisitedErrors := list.New()
		for _, e := range o.errs {
			unVisitedErrors.PushFront(e)
		}

		// Both have errors, check 'em
		for _, e := range n.errs {
			for i := unVisitedErrors.Front(); i != nil; {
				if e.IsEqual(i.Value.(ifaces.ValidateError)) {
					c := i
					i = c.Next()
					// We've visited this error, mark it as visited by removing it from the list
					unVisitedErrors.Remove(c)
				} else {
					i = i.Next()
				}
			}
		}

		// We didn't remove all of the errors, that means some were not in unVisitedErrors
		if unVisitedErrors.Len() != 0 {
			return false
		}
	}

	// Numbered children
	{
		nn := n.NumberedChildren
		oo := o.NumberedChildren

		// local errors match, now we need to traverse the children, if any
		if nn != nil && oo == nil {
			return false
		}
		if oo != nil && nn == nil {
			return false
		}

		if nn != nil {
			// Both now have numbered children, go through each
			if len(nn) != len(oo) {
				return false
			}
			for i := 0; i < len(nn); i++ {
				if !nn[i].IsEqual(oo[i]) {
					return false
				}
			}
		}
	}

	// Named children
	{
		nn := n.NamedChildren
		oo := o.NamedChildren

		// local errors match, now we need to traverse the children, if any
		if nn != nil && oo == nil {
			return false
		}
		if oo != nil && nn == nil {
			return false
		}

		if nn != nil {
			// Both now have numbered children, go through each
			if len(nn) != len(oo) {
				return false
			}
			for i := range nn {
				if !nn[i].IsEqual(oo[i]) {
					return false
				}
			}
		}
	}

	return true
}
