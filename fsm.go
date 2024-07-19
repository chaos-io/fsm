package fsm

import (
	"fmt"
)

// Machine represents a finite state machine.
// Instances of this struct should not be constructed by hand and instead
// should be created using a blueprint.
type Machine struct {
	// A list of available transitions.
	transitions list
	// The current machine state.
	state string
}

// isLegal returns whether or not the specified transition from state a to b
// is legal.
func (f *Machine) isLegal(a string, b string) bool {
	return f.transitions.Search(serialize(a, b)) != nil
}

// Allows returns whether or not this machine can transition to the state b.
func (f *Machine) Allows(b string) bool {
	return f.isLegal(f.state, b)
}

// Disallows returns whether or not this machine can't transition to the state
// b.
func (f *Machine) Disallows(b string) bool {
	return !f.Allows(b)
}

// State returns the current state.
func (f *Machine) State() string {
	return f.state
}

// Goto moves the machine to the specified state. An error is returned if the
// transition is not valid.
func (f *Machine) Goto(state string) error {
	t := f.transitions.Search(serialize(f.state, state))
	if t == nil {
		return fmt.Errorf("can't transition from state %s to %s", f.state, state)
	}

	f.state = state
	if t.fn != nil {
		t.fn(f)
	}

	return nil
}
