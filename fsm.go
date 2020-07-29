package fsm

import (
	"context"
	"fmt"
)

// type Event uint64
type State interface{}

type Machine struct {
	state  State
	fsm    *FSM
	object IMachine
}

type IMachine interface {
	GetState() State
	SetState(context.Context, State) error
	// BeforeSetState(context.Context) error
	OnInitWithMachine(*FSM)
}

func (m *Machine) Goto(s State, ctx context.Context, args ...interface{}) error {
	if s == m.state {
		return nil
	}

	fn, ok := m.fsm.GetHandleFunc(m.state, s)
	isSpecial := m.fsm.IsSpecial(s)
	// fmt.Println(m.state, s, fn, ok, isSpecial)
	if !ok && !isSpecial { //如果没有，并且不是特殊的函数
		return fmt.Errorf("Transition %v to %v not permitted", m.state, s)
	}

	// err := m.object.BeforeSetState(ctx)
	// if err != nil {
	// 	return err
	// }
	{
		stateFuncs, ok := m.fsm.GetStateOnFuncs(m.state)
		if ok && stateFuncs.onExit != nil {
			err := stateFuncs.onExit(m.object, ctx, args...)
			if err != nil {
				return err
			}
		}
	}
	{
		if fn != nil {
			err := fn(m.object, ctx, m.state, s, args...)
			if err != nil {
				return err
			}
		}
	}
	{
		stateFuncs, ok := m.fsm.GetStateOnFuncs(s)
		if ok && stateFuncs.onEnter != nil {
			err := stateFuncs.onEnter(m.object, ctx, args...)
			if err != nil {
				return err
			}
		}
	}

	if m.object.GetState() != s {
		err := m.object.SetState(ctx, m.object.GetState())
		if err != nil {
			return err
		}
		m.state = m.object.GetState()
	}else{
		err := m.object.SetState(ctx, s)
		if err != nil {
			return err
		}
		m.state = s
	}
	return nil
}

type OnEnterFunc func(IMachine, context.Context, ...interface{}) error
type OnExitFunc func(IMachine, context.Context, ...interface{}) error

type FSMState struct {
	onEnter OnEnterFunc
	onExit  OnExitFunc
}

func (ft *FSMState) SetOnEnter(fn OnEnterFunc) *FSMState {
	ft.onEnter = fn
	return ft
}

func (ft *FSMState) SetOnExit(fn OnExitFunc) *FSMState {
	ft.onExit = fn
	return ft
}

type HandleFunc func(IMachine, context.Context, State, State, ...interface{}) error

type FSM struct {
	// State State
	rules         map[State]map[State]HandleFunc
	currentState  State
	toState       State
	states        map[State]*FSMState
	specialStates map[State]bool
}

func (fsm *FSM) GetHandleFunc(from State, to State) (HandleFunc, bool) {
	if from == to {
		return nil, true
	}
	maps, ok := fsm.rules[from]
	if !ok {
		return nil, false
	}
	fn, ok := maps[to]
	return fn, ok
	// if !ok {
	// retu
	// }
	// return fn, nil
}

func NewFSM() *FSM {
	f := &FSM{
		rules:         make(map[State]map[State]HandleFunc, 10),
		states:        make(map[State]*FSMState, 10),
		specialStates: make(map[State]bool, 10),
	}
	return f
}

func (fsm *FSM) Machine(object IMachine) *Machine {
	object.OnInitWithMachine(fsm)
	return &Machine{
		state:  object.GetState(),
		fsm:    fsm,
		object: object,
	}
}

func (fsm *FSM) GetStateOnFuncs(s State) (*FSMState, bool) {
	_s, ok := fsm.states[s]
	return _s, ok
	// return nil, nil
}

func (fsm *FSM) SetStateFuncs(s State, onExit OnExitFunc, onEnter OnEnterFunc) {
	_s, ok := fsm.states[s]
	if !ok {
		_s = &FSMState{}
		fsm.states[s] = _s
	}
	_s.onEnter = onEnter
	_s.onExit = onExit
}

func (fsm *FSM) From(s State) *FSM {
	_, ok := fsm.rules[s]
	if !ok {
		fsm.rules[s] = make(map[State]HandleFunc, 10)
	}
	fsm.currentState = s
	fsm.toState = s
	return fsm
}

func (fsm *FSM) Special(s State) {
	fsm.specialStates[s] = true
}

func (fsm *FSM) IsSpecial(s State) bool {
	_, ok := fsm.specialStates[s]
	return ok
}

func (fsm *FSM) To(s State) *FSM {
	fsm.toState = s
	fsm.rules[fsm.currentState][s] = nil
	return fsm
	// fsm.rules[fsm.currentState][s]
}

func (fsm *FSM) Then(fn HandleFunc) {
	fsm.rules[fsm.currentState][fsm.toState] = fn
	// return fn(fsm.currentState, fsm.toState)
}

// //-----------------------------------------------------------
// type S int

// type AA struct {
// 	State int
// }

// func (a *AA) Change(from State, to State, args ...interface{}) error {
// 	fmt.Println(from, to, args)
// 	return nil
// }

// func main() {

// 	a := &AA{}

// 	s1 := S(1)
// 	s2 := S(2)
// 	s3 := S(3)
// 	fsm := NewFSM()
// 	fsm.SetStateFuncs(s1, func(args ...interface{}) {
// 		fmt.Println("on exit s1")
// 	}, nil)
// 	fsm.SetStateFuncs(s2, nil, func(args ...interface{}) {
// 		fmt.Println("on enter s2")
// 	})
// 	// fsm.Start(s1)
// 	fsm.From(s1).To(s2).Then(a.Change)

// 	m := fsm.Machine(s1)
// 	m.Goto(s2)
// 	err := m.Goto(s3)
// 	fmt.Println(err)
// }




