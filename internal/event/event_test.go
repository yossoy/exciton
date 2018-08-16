package event

import (
	"reflect"
	"testing"
	"time"
)

type testStruct1 struct {
	name       string
	ch         chan Value
	isError    bool
	withResult bool
}
type recvResult struct {
	source *testStruct1
	event  *Event
}

func (ts *testStruct1) eventHandler(e *Event) {
	ts.ch <- NewValue(recvResult{source: ts, event: e})
}

func (ts *testStruct1) eventHandlerWithResult(e *Event, callback ResponceCallback) {
	callback(NewValueResult(NewValue(recvResult{source: ts, event: e})))
}

func TestValue(t *testing.T) {
	val1 := NewValue(1)
	valNil := NewValue(nil)

	var di1 int
	var ds1 recvResult
	var dps1 *recvResult
	var df1 float64
	err := val1.Decode(&di1)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
	err = val1.Decode(&ds1)
	if err == nil {
		t.Errorf("Decode invalid succeeded: %v", ds1)
	} else {
		t.Logf("Decode result: %v", di1)
	}
	err = valNil.Decode(&ds1)
	if err == nil {
		t.Errorf("Decode invalid succeeded: %v", ds1)
	}
	err = valNil.Decode(&dps1)
	if err != nil {
		t.Errorf("Decode failed: %v", err)
	}
	err = val1.Decode(&df1)
	if err == nil {
		t.Errorf("Decode invalid succeeded: %v", df1)
	}
}

func TestEventMgr(t *testing.T) {
	StartEventMgr()

	ch := make(chan Value)
	val1 := NewValue(1)
	val2 := NewValue(2)
	val3 := NewValue(3)
	// Test sources
	sources := []*testStruct1{
		&testStruct1{name: "/test1", isError: false},
		&testStruct1{name: "/p/:param1", isError: false},
		&testStruct1{name: "/test3", isError: false, withResult: true},
	}
	// Test cases
	cases := []struct {
		path       string
		isError    bool
		source     *testStruct1
		args       Value
		params     map[string]string
		withResult bool
	}{
		{path: "/test1", isError: false, source: sources[0], args: val1},
		{path: "/test2", isError: true},
		{path: "/p/foo", isError: false, source: sources[1], args: val2, params: map[string]string{"param1": "foo"}},
		{path: "/test1", isError: true, withResult: true},
		{path: "/test3", isError: false, source: sources[2], args: val3, withResult: true},
	}

	for idx, s := range sources {
		var err error
		if s.withResult {
			err = AddHandlerWithResult(s.name, s.eventHandlerWithResult)
		} else {
			s.ch = ch
			err = AddHandler(s.name, s.eventHandler)
		}
		if s.isError {
			if err == nil {
				t.Errorf("[%d][%s] event regist not error\n", idx, s.name)
			}
		} else {
			if err != nil {
				t.Errorf("[%d][%s] event regist error: %v\n", idx, s.name, err)
			}
		}
	}
	for idx, c := range cases {
		t.Logf("[%d][%s] testing...\n", idx, c.path)
		//TODO: Emit is need result for testing?
		onError := false
		var r Value
		var err error
		if c.withResult {
			err = EmitWithCallback(c.path, c.args, func(result Result) {
				ch <- result.Value()
			})
		} else {
			err = Emit(c.path, c.args)
		}
		if err == nil {
			select {
			case r = <-ch:
			case <-time.After(1 * time.Second):
				onError = true
			}
		}
		if c.isError {
			if err == nil {
				t.Errorf("[%d][%s] Emit not error: %v\n", idx, c.path, r)
			}
		} else {
			if onError {
				t.Errorf("[%d][%s] Emit timeout\n", idx, c.path)
			} else {
				var rr recvResult
				err = r.Decode(&rr)

				if err != nil {
					t.Errorf("[%d][%s] Return value decode fail.: %v", idx, c.params, err)
				} else if !reflect.DeepEqual(rr.source, c.source) {
					t.Errorf("[%d][%s] Emit handler mismatch: %v vs %v\n", idx, c.path, rr.source, c.source)
				} else if rr.event.Name != c.source.name {
					t.Errorf("[%d][%s] Emit recv event name invalid: %q vs %q\n", idx, c.path, rr.event.Name, c.path)
				} else if rr.event.Argument != c.args {
					t.Errorf("[%d][%s] Emit recv argument invalid: %v vs %v\n", idx, c.path, rr.event.Argument, c.args)
				} else if !reflect.DeepEqual(rr.event.Params, c.params) {
					t.Errorf("[%d][%s] Emit recv event parameter invalid: %v vs %v\n", idx, c.path, rr.event.Params, c.params)
				}
			}
		}
	}

	StopEventMgr()
}
