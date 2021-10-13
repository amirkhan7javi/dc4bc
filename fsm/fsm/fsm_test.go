package fsm

import (
	"testing"
)

const (
	testName = "fsm_test"
	// Init process from global idle state
	stateInit = StateGlobalIdle
	// Set up data
	stateStage1 = State("state_stage1")
	// Process data
	stateStage2 = State("state_stage2")
	// Canceled with internal event
	stateCanceledByInternal = State("state_canceled")
	// Canceled with external event
	stateCanceled2 = State("state_canceled2")
	// Out endpoint to switch
	stateOutToFSM2 = State("state_out_to_fsm2")

	// Events
	eventInit    = Event("event_init")
	eventCancel  = Event("event_cancel")
	eventProcess = Event("event_process")

	// Internal events
	eventInternal         = Event("event_internal")
	eventCancelByInternal = Event("event_internal_cancel")
	eventInternalOut2     = Event("event_internal_out")
)

var (
	testingFSM *FSM

	testingEvents = []EventDesc{
		// Init
		{Name: eventInit, SrcState: []State{stateInit}, DstState: stateStage1},
		{Name: eventInternal, SrcState: []State{stateStage1}, DstState: stateStage2, IsInternal: true},

		// Cancellation events
		{Name: eventCancelByInternal, SrcState: []State{stateStage2}, DstState: stateCanceledByInternal, IsInternal: true},
		{Name: eventCancel, SrcState: []State{stateStage2}, DstState: stateCanceled2},

		// Out
		{Name: eventProcess, SrcState: []State{stateStage2}, DstState: stateOutToFSM2},
		{Name: eventInternalOut2, SrcState: []State{stateStage2}, DstState: stateOutToFSM2, IsInternal: true},
	}

	testingCallbacks = Callbacks{
		eventInit: func(event Event, args ...interface{}) (Event, interface{}, error) {
			return event, nil, nil
		},
		eventInternalOut2: func(event Event, args ...interface{}) (Event, interface{}, error) {
			return event, nil, nil
		},
		eventProcess: func(event Event, args ...interface{}) (Event, interface{}, error) {
			return event, nil, nil
		},
	}
)

func init() {
	testingFSM = MustNewFSM(
		testName,
		stateInit,
		testingEvents,
		testingCallbacks,
	)
}

func compareRecoverStr(t *testing.T, r interface{}, assertion string) {
	if r == nil {
		return
	}
	msg, ok := r.(string)
	if !ok {
		t.Error("not asserted recover:", r)
	}
	if msg != assertion {
		t.Error("not asserted recover:", msg)
	}
}

func compareStatesArr(src, dst []State) bool {
	if len(src) != len(dst) {
		return false
	}
	// create a map of string -> int
	diff := make(map[State]int, len(src))
	for _, _x := range src {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range dst {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

func compareEventsArr(src, dst []Event) bool {
	if len(src) != len(dst) {
		return false
	}
	// create a map of string -> int
	diff := make(map[Event]int, len(src))
	for _, _x := range src {
		// 0 value for int is 0, so just increment a counter for the string
		diff[_x]++
	}
	for _, _y := range dst {
		// If the string _y is not in diff bail out early
		if _, ok := diff[_y]; !ok {
			return false
		}
		diff[_y] -= 1
		if diff[_y] == 0 {
			delete(diff, _y)
		}
	}
	return len(diff) == 0
}

func TestMustNewFSM_Empty_Name_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "machine name cannot be empty")
	}()
	testingFSM = MustNewFSM(
		"",
		"init_state",
		[]EventDesc{},
		nil,
	)

	t.Errorf("did not panic on empty machine name")
}

func TestMustNewFSM_Empty_Initial_State_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "initial state state cannot be empty")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"",
		[]EventDesc{},
		nil,
	)

	t.Errorf("did not panic on empty initial")
}

func TestMustNewFSM_Empty_Events_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "cannot init fsm with empty events")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"init_state",
		[]EventDesc{},
		nil,
	)

	t.Errorf("did not panic on empty events list")
}

func TestMustNewFSM_Event_Empty_Name_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "cannot init empty event")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"init_state",
		[]EventDesc{
			{Name: "", SrcState: []State{"init_state"}, DstState: StateGlobalDone},
		},
		nil,
	)

	t.Errorf("did not panic on empty event name")
}

func TestMustNewFSM_Event_Empty_Source_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "event must have minimum one source available state")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"init_state",
		[]EventDesc{
			{Name: "event", SrcState: []State{}, DstState: StateGlobalDone},
		},
		nil,
	)

	t.Errorf("did not panic on empty event sources")
}

func TestMustNewFSM_States_Min_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "machine must contain at least two states")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"init_state",
		[]EventDesc{
			{Name: "event", SrcState: []State{"init_state"}, DstState: StateGlobalDone},
		},
		nil,
	)

	t.Errorf("did not panic on less than two states")
}

func TestMustNewFSM_State_Entry_Conflict_Panic(t *testing.T) {
	defer func() {
		compareRecoverStr(t, recover(), "machine must contain at least two states")
	}()

	testingFSM = MustNewFSM(
		"fsm",
		"init_state",
		[]EventDesc{
			{Name: "event1", SrcState: []State{"init_state"}, DstState: "state"},
			{Name: "event2", SrcState: []State{"init_state"}, DstState: "state"},
		},
		nil,
	)

	t.Errorf("did not panic on initialize with conflict in entry state")
}

func TestFSM_Name(t *testing.T) {
	if testingFSM.Name() != testName {
		t.Errorf("expected machine name \"%s\"", testName)
	}
}

func TestFSM_EntryEvent(t *testing.T) {
	if testingFSM.InitialState() != stateInit {
		t.Errorf("expected initial state \"%s\"", stateInit)
	}
}

func TestFSM_EventsList(t *testing.T) {
	eventsList := []Event{
		eventInit,
		eventCancel,
		eventProcess,
	}

	if !compareEventsArr(testingFSM.EventsList(), eventsList) {
		t.Error("expected public events", eventsList)
	}

}

func TestFSM_StatesList(t *testing.T) {
	statesList := []State{
		stateInit,
		stateStage1,
		stateStage2,
	}

	if !compareStatesArr(testingFSM.StatesList(), statesList) {
		t.Error("expected states", statesList)
	}
}

func TestFSM_CopyWithState(t *testing.T) {
	testingFSM1 := MustNewFSM(
		testName,
		stateInit,
		testingEvents,
		testingCallbacks,
	)
	testingFSM2 := testingFSM.MustCopyWithState(stateStage2)

	if testingFSM1.State() != stateInit {
		t.Fatal("expect unchanged source state")
	}

	if testingFSM2.State() != stateStage2 {
		t.Fatal("expect changed source state")
	}

}
