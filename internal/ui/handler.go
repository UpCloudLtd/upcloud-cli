package ui

import (
	"fmt"
	"reflect"
	"sync"
)

type HandleContext struct {
	Requests      interface{}
	RequestId     func(interface{}) string
	ResultUuid    func(interface{}) string
	MessageFn	 		func(interface{}) string
	ActionMsg     string
	Action        func(interface{}) (interface{}, error)
	MaxActions    int
	InteractiveUi bool
	WaitMsg       string
	WaitFn        func(uuid string, waitMsg string, err error) (interface{}, error)
}

func (c HandleContext) HandleAction() (interface{}, error) {
	var (
		mu      sync.Mutex
		numOk   int
		results []interface{}
	)

	var elems []interface{}
	if reflect.TypeOf(c.Requests).Kind() == reflect.Slice {
		is := reflect.ValueOf(c.Requests)
		for i := 0; i < is.Len(); i++ {
			elems = append(elems, is.Index(i).Interface())
		}
	} else {
		elems = append(elems, c.Requests)
	}

	handler := func(idx int, e *LogEntry) {
		request := elems[idx]
		var msg string
		if c.MessageFn != nil {
			msg = c.MessageFn(request)
		} else if c.RequestId != nil && c.ActionMsg != "" {
			msg = fmt.Sprintf("%s %s", c.ActionMsg, c.RequestId(request))
		}
		e.SetMessage(msg)
		e.Start()

		var details interface{}
		var err error
		details, err = c.Action(request)

		var detailsUuid string
		if c.ResultUuid != nil && details != nil && !reflect.ValueOf(details).IsNil() {
			detailsUuid = c.ResultUuid(details)
		} else if c.RequestId != nil{
			detailsUuid = c.RequestId(request)
		}

		if c.WaitFn != nil {
			e.SetMessage(fmt.Sprintf("%s: %s", msg, c.WaitMsg))
			details, err = c.WaitFn(detailsUuid, c.WaitMsg, err)
		}
		if err != nil {
			e.SetMessage(LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
			e.SetDetails(err.Error(), "error: ")
		} else {
			e.SetMessage(fmt.Sprintf("%s: done", msg))
			if c.ResultUuid != nil {
				e.SetDetails(detailsUuid, "UUID: ")
			}
			mu.Lock()
			numOk++
			results = append(results, details)
			mu.Unlock()
		}
	}

	StartWorkQueue(WorkQueueConfig{
		NumTasks:           len(elems),
		MaxConcurrentTasks: c.MaxActions,
		EnableUI:           c.InteractiveUi,
	}, handler)

	if numOk != len(elems) {
		return nil, fmt.Errorf("number of operations failed: %d", len(elems)-numOk)
	}

	return results, nil
}



func (c HandleContext) Handle(requests []interface{}) (interface{}, error) {
	var (
		mu      sync.Mutex
		numOk   int
		results []interface{}
	)

	handler := func(idx int, e *LogEntry) {
		request := requests[idx]
		var msg string
		if c.MessageFn != nil {
			msg = c.MessageFn(request)
		} else if c.RequestId != nil && c.ActionMsg != "" {
			msg = fmt.Sprintf("%s %s", c.ActionMsg, c.RequestId(request))
		}
		e.SetMessage(msg)
		e.Start()

		var details interface{}
		var err error
		details, err = c.Action(request)

		var detailsUuid string
		if c.ResultUuid != nil && details != nil && !reflect.ValueOf(details).IsNil() {
			detailsUuid = c.ResultUuid(details)
		} else if c.RequestId != nil{
			detailsUuid = c.RequestId(request)
		}

		if c.WaitFn != nil {
			e.SetMessage(fmt.Sprintf("%s: %s", msg, c.WaitMsg))
			details, err = c.WaitFn(detailsUuid, c.WaitMsg, err)
		}
		if err != nil {
			e.SetMessage(LiveLogEntryErrorColours.Sprintf("%s: failed", msg))
			e.SetDetails(err.Error(), "error: ")
		} else {
			e.SetMessage(fmt.Sprintf("%s: done", msg))
			if c.ResultUuid != nil {
				e.SetDetails(detailsUuid, "UUID: ")
			}
			mu.Lock()
			numOk++
			results = append(results, details)
			mu.Unlock()
		}
	}

	StartWorkQueue(WorkQueueConfig{
		NumTasks:           len(requests),
		MaxConcurrentTasks: c.MaxActions,
		EnableUI:           c.InteractiveUi,
	}, handler)

	if numOk != len(requests) {
		return nil, fmt.Errorf("number of operations failed: %d", len(requests)-numOk)
	}

	return results, nil
}
