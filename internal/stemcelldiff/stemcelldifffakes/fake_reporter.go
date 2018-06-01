// Code generated by counterfeiter. DO NOT EDIT.
package stemcelldifffakes

import (
	"sync"
)

type FakeReporter struct {
	PrintReportStub        func(report string)
	printReportMutex       sync.RWMutex
	printReportArgsForCall []struct {
		report string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeReporter) PrintReport(report string) {
	fake.printReportMutex.Lock()
	fake.printReportArgsForCall = append(fake.printReportArgsForCall, struct {
		report string
	}{report})
	fake.recordInvocation("PrintReport", []interface{}{report})
	fake.printReportMutex.Unlock()
	if fake.PrintReportStub != nil {
		fake.PrintReportStub(report)
	}
}

func (fake *FakeReporter) PrintReportCallCount() int {
	fake.printReportMutex.RLock()
	defer fake.printReportMutex.RUnlock()
	return len(fake.printReportArgsForCall)
}

func (fake *FakeReporter) PrintReportArgsForCall(i int) string {
	fake.printReportMutex.RLock()
	defer fake.printReportMutex.RUnlock()
	return fake.printReportArgsForCall[i].report
}

func (fake *FakeReporter) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.printReportMutex.RLock()
	defer fake.printReportMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeReporter) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}