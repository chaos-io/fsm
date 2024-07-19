package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	A = "A"
	B = "B"
	C = "C"
	D = "D"
)

func makeMachine() *Machine {
	bp := New()
	bp.Start(A)
	bp.From(A).To(B)
	bp.From(B).To(C)
	bp.From(B).To(B)
	return bp.Machine()
}

func TestWorksNormally(t *testing.T) {
	m := makeMachine()
	assert.Equal(t, A, m.State())
	assert.NoError(t, m.Goto(B))
	assert.Equal(t, B, m.State())
	assert.NoError(t, m.Goto(C))
	assert.Equal(t, C, m.State())
	err := m.Goto(B)
	assert.Error(t, err)
	assert.Equal(t, C, m.State())
	assert.Equal(t, "can't transition from state C to B", err.Error())
}

func TestAddsHandler(t *testing.T) {
	bp := New()
	bp.Start(A)
	out := []uint8{}
	bp.From(A).To(B).Then(func(m *Machine) { out = append(out, 1) })
	bp.From(B).To(C)
	bp.From(C).To(D).Then(func(m *Machine) { out = append(out, 2) })
	m := bp.Machine()

	assert.Equal(t, []uint8{}, out)
	m.Goto(B)
	assert.Equal(t, []uint8{1}, out)
	m.Goto(C)
	assert.Equal(t, []uint8{1}, out)
	m.Goto(D)
	assert.Equal(t, []uint8{1, 2}, out)
}

func BenchmarkTransitions(b *testing.B) {
	m := makeMachine()
	m.Goto(B)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		m.Goto(B)
	}
}

func BenchmarkAllows(b *testing.B) {
	m := makeMachine()
	for n := 0; n < b.N; n++ {
		m.Allows(B)
	}
}

func BenchmarkGetState(b *testing.B) {
	m := makeMachine()
	for n := 0; n < b.N; n++ {
		m.State()
	}
}

func BenchmarkBuildMachine(b *testing.B) {
	for n := 0; n < b.N; n++ {
		makeMachine()
	}
}

func TestBlueprint(t *testing.T) {
	/*
		A -> B -> C
		A -> B <-> D
	*/
	bp := New()
	bp.Start(A)
	bp.Print()
	bp.From(A).To(B)
	bp.Print()
	bp.From(B).To(C).Then(func(m *Machine) { fmt.Println("from B to C") })
	bp.Print()
	bp.From(B).To(D).Then(func(m *Machine) { fmt.Println("from B to D") })
	bp.Print()
	bp.From(D).To(B).Then(func(m *Machine) { fmt.Println("from D to B") })
	bp.Print()

	m := bp.Machine()
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(C))

	m = bp.Machine()
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(D))
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(D))

	m = bp.Machine()
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(D))
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(C))

	m = bp.Machine()
	assert.NoError(t, m.Goto(B))
	assert.NoError(t, m.Goto(C))
	assert.Error(t, m.Goto(D))
}
