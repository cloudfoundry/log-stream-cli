package testhelpers

import (
	"fmt"
	"io"
	"sync"

	. "github.com/cloudfoundry/cli/cf/terminal"
)

func NewFakePrinter(stdout io.Writer) *FakePrinter {
	return &FakePrinter{
		PrintfStub: func(format string, a ...interface{}) (int, error) {
			return fmt.Fprintf(stdout, format, a...)
		},
	}
}

type FakePrinter struct {
	PrintStub        func(a ...interface{}) (n int, err error)
	printMutex       sync.RWMutex
	printArgsForCall []struct {
		a []interface{}
	}
	printReturns struct {
		result1 int
		result2 error
	}
	PrintfStub        func(format string, a ...interface{}) (n int, err error)
	printfMutex       sync.RWMutex
	printfArgsForCall []struct {
		format string
		a      []interface{}
	}
	printfReturns struct {
		result1 int
		result2 error
	}
	PrintlnStub        func(a ...interface{}) (n int, err error)
	printlnMutex       sync.RWMutex
	printlnArgsForCall []struct {
		a []interface{}
	}
	printlnReturns struct {
		result1 int
		result2 error
	}
	ForcePrintStub        func(a ...interface{}) (n int, err error)
	forcePrintMutex       sync.RWMutex
	forcePrintArgsForCall []struct {
		a []interface{}
	}
	forcePrintReturns struct {
		result1 int
		result2 error
	}
	ForcePrintfStub        func(format string, a ...interface{}) (n int, err error)
	forcePrintfMutex       sync.RWMutex
	forcePrintfArgsForCall []struct {
		format string
		a      []interface{}
	}
	forcePrintfReturns struct {
		result1 int
		result2 error
	}
	ForcePrintlnStub        func(a ...interface{}) (n int, err error)
	forcePrintlnMutex       sync.RWMutex
	forcePrintlnArgsForCall []struct {
		a []interface{}
	}
	forcePrintlnReturns struct {
		result1 int
		result2 error
	}
}

func (fake *FakePrinter) Print(a ...interface{}) (n int, err error) {
	fake.printMutex.Lock()
	defer fake.printMutex.Unlock()
	fake.printArgsForCall = append(fake.printArgsForCall, struct {
		a []interface{}
	}{a})
	if fake.PrintStub != nil {
		return fake.PrintStub(a)
	} else {
		return fake.printReturns.result1, fake.printReturns.result2
	}
}

func (fake *FakePrinter) PrintCallCount() int {
	fake.printMutex.RLock()
	defer fake.printMutex.RUnlock()
	return len(fake.printArgsForCall)
}

func (fake *FakePrinter) PrintArgsForCall(i int) []interface{} {
	fake.printMutex.RLock()
	defer fake.printMutex.RUnlock()
	return fake.printArgsForCall[i].a
}

func (fake *FakePrinter) PrintReturns(result1 int, result2 error) {
	fake.printReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) Printf(format string, a ...interface{}) (n int, err error) {
	fake.printfMutex.Lock()
	defer fake.printfMutex.Unlock()
	fake.printfArgsForCall = append(fake.printfArgsForCall, struct {
		format string
		a      []interface{}
	}{format, a})
	if fake.PrintfStub != nil {
		return fake.PrintfStub(format, a...)
	} else {
		return fake.printfReturns.result1, fake.printfReturns.result2
	}
}

func (fake *FakePrinter) PrintfCallCount() int {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return len(fake.printfArgsForCall)
}

func (fake *FakePrinter) PrintfArgsForCall(i int) (string, []interface{}) {
	fake.printfMutex.RLock()
	defer fake.printfMutex.RUnlock()
	return fake.printfArgsForCall[i].format, fake.printfArgsForCall[i].a
}

func (fake *FakePrinter) PrintfReturns(result1 int, result2 error) {
	fake.printfReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) Println(a ...interface{}) (n int, err error) {
	fake.printlnMutex.Lock()
	defer fake.printlnMutex.Unlock()
	fake.printlnArgsForCall = append(fake.printlnArgsForCall, struct {
		a []interface{}
	}{a})
	if fake.PrintlnStub != nil {
		return fake.PrintlnStub(a)
	} else {
		return fake.printlnReturns.result1, fake.printlnReturns.result2
	}
}

func (fake *FakePrinter) PrintlnCallCount() int {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return len(fake.printlnArgsForCall)
}

func (fake *FakePrinter) PrintlnArgsForCall(i int) []interface{} {
	fake.printlnMutex.RLock()
	defer fake.printlnMutex.RUnlock()
	return fake.printlnArgsForCall[i].a
}

func (fake *FakePrinter) PrintlnReturns(result1 int, result2 error) {
	fake.printlnReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) ForcePrint(a ...interface{}) (n int, err error) {
	fake.forcePrintMutex.Lock()
	defer fake.forcePrintMutex.Unlock()
	fake.forcePrintArgsForCall = append(fake.forcePrintArgsForCall, struct {
		a []interface{}
	}{a})
	if fake.ForcePrintStub != nil {
		return fake.ForcePrintStub(a)
	} else {
		return fake.forcePrintReturns.result1, fake.forcePrintReturns.result2
	}
}

func (fake *FakePrinter) ForcePrintCallCount() int {
	fake.forcePrintMutex.RLock()
	defer fake.forcePrintMutex.RUnlock()
	return len(fake.forcePrintArgsForCall)
}

func (fake *FakePrinter) ForcePrintArgsForCall(i int) []interface{} {
	fake.forcePrintMutex.RLock()
	defer fake.forcePrintMutex.RUnlock()
	return fake.forcePrintArgsForCall[i].a
}

func (fake *FakePrinter) ForcePrintReturns(result1 int, result2 error) {
	fake.forcePrintReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) ForcePrintf(format string, a ...interface{}) (n int, err error) {
	fake.forcePrintfMutex.Lock()
	defer fake.forcePrintfMutex.Unlock()
	fake.forcePrintfArgsForCall = append(fake.forcePrintfArgsForCall, struct {
		format string
		a      []interface{}
	}{format, a})
	if fake.ForcePrintfStub != nil {
		return fake.ForcePrintfStub(format, a...)
	} else {
		return fake.forcePrintfReturns.result1, fake.forcePrintfReturns.result2
	}
}

func (fake *FakePrinter) ForcePrintfCallCount() int {
	fake.forcePrintfMutex.RLock()
	defer fake.forcePrintfMutex.RUnlock()
	return len(fake.forcePrintfArgsForCall)
}

func (fake *FakePrinter) ForcePrintfArgsForCall(i int) (string, []interface{}) {
	fake.forcePrintfMutex.RLock()
	defer fake.forcePrintfMutex.RUnlock()
	return fake.forcePrintfArgsForCall[i].format, fake.forcePrintfArgsForCall[i].a
}

func (fake *FakePrinter) ForcePrintfReturns(result1 int, result2 error) {
	fake.forcePrintfReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

func (fake *FakePrinter) ForcePrintln(a ...interface{}) (n int, err error) {
	fake.forcePrintlnMutex.Lock()
	defer fake.forcePrintlnMutex.Unlock()
	fake.forcePrintlnArgsForCall = append(fake.forcePrintlnArgsForCall, struct {
		a []interface{}
	}{a})
	if fake.ForcePrintlnStub != nil {
		return fake.ForcePrintlnStub(a)
	} else {
		return fake.forcePrintlnReturns.result1, fake.forcePrintlnReturns.result2
	}
}

func (fake *FakePrinter) ForcePrintlnCallCount() int {
	fake.forcePrintlnMutex.RLock()
	defer fake.forcePrintlnMutex.RUnlock()
	return len(fake.forcePrintlnArgsForCall)
}

func (fake *FakePrinter) ForcePrintlnArgsForCall(i int) []interface{} {
	fake.forcePrintlnMutex.RLock()
	defer fake.forcePrintlnMutex.RUnlock()
	return fake.forcePrintlnArgsForCall[i].a
}

func (fake *FakePrinter) ForcePrintlnReturns(result1 int, result2 error) {
	fake.forcePrintlnReturns = struct {
		result1 int
		result2 error
	}{result1, result2}
}

var _ Printer = new(FakePrinter)
