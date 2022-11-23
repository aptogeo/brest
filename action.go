package brest

import "strings"

// Action type
type Action int

const (
	// Get action
	Get Action = 1 << iota
	// Post action
	Post
	// Put action
	Put
	// Patch action
	Patch
	// Delete action
	Delete
)

// All actions
const All Action = Get + Post + Put + Patch + Delete

// None action
const None Action = 0

func (a Action) String() string {
	strs := make([]string, 0)
	if a&Get > 0 {
		strs = append(strs, "Get")
	}
	if a&Post > 0 {
		strs = append(strs, "Post")
	}
	if a&Put > 0 {
		strs = append(strs, "Put")
	}
	if a&Patch > 0 {
		strs = append(strs, "Patch")
	}
	if a&Delete > 0 {
		strs = append(strs, "Delete")
	}
	if len(strs) == 0 {
		return "None"
	}
	if len(strs) == 5 {
		return "All"
	}
	return strings.Join(strs, "|")
}
