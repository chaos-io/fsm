package fsm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	A uint8 = iota
	B
	C
	D
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
	assert.Equal(t, "can't transition from state 2 to 1", err.Error())
}

func TestAddsHandler(t *testing.T) {
	bp := New()
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
		0 -> 1 -> 2 -> 4
		0 -> 1 -> 3 -> 4
	*/
	bp := New()
	bp.Start(0)
	bp.Print()
	bp.From(0).To(1)
	bp.Print()
	bp.From(1).To(2).Then(func(m *Machine) { fmt.Println("from 1 to 2") })
	bp.Print()
	bp.From(1).To(3).Then(func(m *Machine) { fmt.Println("from 1 to 3") })
	bp.Print()
	bp.From(2).To(4).Then(func(m *Machine) { fmt.Println("from 2 to 4") })
	bp.Print()
	bp.From(3).To(4).Then(func(m *Machine) { fmt.Println("from 3 to 4") })
	bp.Print()

	m := bp.Machine()
	assert.NoError(t, m.Goto(1))
	assert.NoError(t, m.Goto(2))
	assert.NoError(t, m.Goto(4))

	m = bp.Machine()
	assert.NoError(t, m.Goto(1))
	assert.NoError(t, m.Goto(3))
	assert.NoError(t, m.Goto(4))

	m = bp.Machine()
	assert.NoError(t, m.Goto(1))
	assert.NoError(t, m.Goto(2))
	assert.Error(t, m.Goto(3))
}
