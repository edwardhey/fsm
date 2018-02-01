# fsm

## 使用
```golang
package fsm

import (
	"testing"
)

type Status int

const (
	StatusOffline Status = 0
	StatusOnline  Status = 10
)

type Mission struct {
	Status Status
}

func ChangeStatus(from State, to State, args ...interface{}) error {
	return nil
}

func init() {

	// 	m := fsm.Machine(s1)
	// 	m.Goto(s2)
	// 	err := m.Goto(s3)
	// 	fmt.Println(err)
	// fsm.Start(s1)
}

func TestFSM(t *testing.T) {
	var _fsm *FSM
	_fsm = NewFSM()
	exitCall := false
	enterCall := false
	_fsm.SetStateFuncs(StatusOnline, func(args ...interface{}) {
		exitCall = true
	}, nil)
	_fsm.SetStateFuncs(StatusOffline, nil, func(args ...interface{}) {
		enterCall = true
	})
	// _fsm.Special(StatusOnline)
	_fsm.From(StatusOnline).To(StatusOffline).Then(ChangeStatus)

	mission := &Mission{}
	mission.Status = StatusOnline

	m := _fsm.Machine(mission.Status)
	m.Goto(StatusOffline)
	if m.State != StatusOffline {
		t.Error("状态设置成为 offline 失败")
	}

	if exitCall == false {
		t.Error("exit call failed")
	}
	if enterCall == false {
		t.Error("enter call failed")
	}
	m.Goto(StatusOnline)
	if m.State != StatusOffline {
		t.Error("由于为注册 offline -> online，状态不应该变更")
	}

	_fsm.Special(StatusOnline)
	m.Goto(StatusOnline)
	if m.State != StatusOnline {
		t.Error("设置了online为特殊状态后，任意状态 -> online 失败")
	}
}

```
