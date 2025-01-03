package fsm

import (
	"fmt"
	"strings"
)

type Blueprint struct {
	transitions list
	start       string
}

// New creates a new finite state machine blueprint.
func New() *Blueprint {
	return &Blueprint{
		transitions: make(list, 0),
	}
}

// From returns a new transition for the blueprint.
// The transition will be added to the blueprint automatically when it has both
// "from" and "to" values.
func (b *Blueprint) From(start string) *Transition {
	return (&Transition{blueprint: b}).From(start)
}

// Add adds a complete transition to the blueprint.
func (b *Blueprint) Add(t *Transition) {
	idx := b.transitions.InsertPos(t)
	trans := make(list, len(b.transitions)+1)

	copy(trans, b.transitions[:idx])
	copy(trans[idx+1:], b.transitions[idx:])
	trans[idx] = t
	b.transitions = trans
}

// Start sets the start state for the machine.
func (b *Blueprint) Start(state string) {
	b.start = state
}

// Machine returns a new machine created from the blueprint.
func (b *Blueprint) Machine() *Machine {
	fsm := &Machine{
		state:       b.start,
		transitions: b.transitions,
	}

	return fsm
}

// Print print the transition list.
func (b *Blueprint) Print() {
	if len(b.transitions) == 0 {
		fmt.Println("<-")
		return
	}

	builder := strings.Builder{}
	for i, t := range b.transitions {
		if i < len(b.transitions)-1 {
			builder.WriteString(fmt.Sprintf("(%v -> %v) -> ", t.from, t.to))
		} else {
			builder.WriteString(fmt.Sprintf("(%v -> %v)", t.from, t.to))
		}
	}

	fmt.Println(builder.String())
}
