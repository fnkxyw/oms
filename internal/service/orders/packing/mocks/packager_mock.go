// Code generated by http://github.com/gojuno/minimock (v3.4.0). DO NOT EDIT.

package mocks

//go:generate minimock -i gitlab.ozon.dev/akugnerevich/homework.git/internal/service/orders/packing.Packager -o packager_mock_test.go -n PackagerMock -p mocks

import (
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"gitlab.ozon.dev/akugnerevich/homework.git/internal/models"
)

// PackagerMock implements mm_packing.Packager
type PackagerMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcPack          func(o *models.Order) (err error)
	funcPackOrigin    string
	inspectFuncPack   func(o *models.Order)
	afterPackCounter  uint64
	beforePackCounter uint64
	PackMock          mPackagerMockPack
}

// NewPackagerMock returnStorage a mock for mm_packing.Packager
func NewPackagerMock(t minimock.Tester) *PackagerMock {
	m := &PackagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.PackMock = mPackagerMockPack{mock: m}
	m.PackMock.callArgs = []*PackagerMockPackParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mPackagerMockPack struct {
	optional           bool
	mock               *PackagerMock
	defaultExpectation *PackagerMockPackExpectation
	expectations       []*PackagerMockPackExpectation

	callArgs []*PackagerMockPackParams
	mutex    sync.RWMutex

	expectedInvocations       uint64
	expectedInvocationsOrigin string
}

// PackagerMockPackExpectation specifies expectation struct of the Packager.Pack
type PackagerMockPackExpectation struct {
	mock               *PackagerMock
	params             *PackagerMockPackParams
	paramPtrs          *PackagerMockPackParamPtrs
	expectationOrigins PackagerMockPackExpectationOrigins
	results            *PackagerMockPackResults
	returnOrigin       string
	Counter            uint64
}

// PackagerMockPackParams contains parameters of the Packager.Pack
type PackagerMockPackParams struct {
	o *models.Order
}

// PackagerMockPackParamPtrs contains pointers to parameters of the Packager.Pack
type PackagerMockPackParamPtrs struct {
	o **models.Order
}

// PackagerMockPackResults contains results of the Packager.Pack
type PackagerMockPackResults struct {
	err error
}

// PackagerMockPackOrigins contains origins of expectations of the Packager.Pack
type PackagerMockPackExpectationOrigins struct {
	origin  string
	originO string
}

// Marks this method to be optional. The default behavior of any method with Return() is '1 or more', meaning
// the test will fail minimock's automatic final call check if the mocked method was not called at least once.
// Optional() makes method check to work in '0 or more' mode.
// It is NOT RECOMMENDED to use this option unless you really need it, as default behaviour helps to
// catch the problems when the expected method call is totally skipped during test run.
func (mmPack *mPackagerMockPack) Optional() *mPackagerMockPack {
	mmPack.optional = true
	return mmPack
}

// Expect sets up expected params for Packager.Pack
func (mmPack *mPackagerMockPack) Expect(o *models.Order) *mPackagerMockPack {
	if mmPack.mock.funcPack != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by Set")
	}

	if mmPack.defaultExpectation == nil {
		mmPack.defaultExpectation = &PackagerMockPackExpectation{}
	}

	if mmPack.defaultExpectation.paramPtrs != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by ExpectParams functions")
	}

	mmPack.defaultExpectation.params = &PackagerMockPackParams{o}
	mmPack.defaultExpectation.expectationOrigins.origin = minimock.CallerInfo(1)
	for _, e := range mmPack.expectations {
		if minimock.Equal(e.params, mmPack.defaultExpectation.params) {
			mmPack.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmPack.defaultExpectation.params)
		}
	}

	return mmPack
}

// ExpectOParam1 sets up expected param o for Packager.Pack
func (mmPack *mPackagerMockPack) ExpectOParam1(o *models.Order) *mPackagerMockPack {
	if mmPack.mock.funcPack != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by Set")
	}

	if mmPack.defaultExpectation == nil {
		mmPack.defaultExpectation = &PackagerMockPackExpectation{}
	}

	if mmPack.defaultExpectation.params != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by Expect")
	}

	if mmPack.defaultExpectation.paramPtrs == nil {
		mmPack.defaultExpectation.paramPtrs = &PackagerMockPackParamPtrs{}
	}
	mmPack.defaultExpectation.paramPtrs.o = &o
	mmPack.defaultExpectation.expectationOrigins.originO = minimock.CallerInfo(1)

	return mmPack
}

// Inspect accepts an inspector function that has same arguments as the Packager.Pack
func (mmPack *mPackagerMockPack) Inspect(f func(o *models.Order)) *mPackagerMockPack {
	if mmPack.mock.inspectFuncPack != nil {
		mmPack.mock.t.Fatalf("Inspect function is already set for PackagerMock.Pack")
	}

	mmPack.mock.inspectFuncPack = f

	return mmPack
}

// Return sets up results that will be returned by Packager.Pack
func (mmPack *mPackagerMockPack) Return(err error) *PackagerMock {
	if mmPack.mock.funcPack != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by Set")
	}

	if mmPack.defaultExpectation == nil {
		mmPack.defaultExpectation = &PackagerMockPackExpectation{mock: mmPack.mock}
	}
	mmPack.defaultExpectation.results = &PackagerMockPackResults{err}
	mmPack.defaultExpectation.returnOrigin = minimock.CallerInfo(1)
	return mmPack.mock
}

// Set uses given function f to mock the Packager.Pack method
func (mmPack *mPackagerMockPack) Set(f func(o *models.Order) (err error)) *PackagerMock {
	if mmPack.defaultExpectation != nil {
		mmPack.mock.t.Fatalf("Default expectation is already set for the Packager.Pack method")
	}

	if len(mmPack.expectations) > 0 {
		mmPack.mock.t.Fatalf("Some expectations are already set for the Packager.Pack method")
	}

	mmPack.mock.funcPack = f
	mmPack.mock.funcPackOrigin = minimock.CallerInfo(1)
	return mmPack.mock
}

// When sets expectation for the Packager.Pack which will trigger the result defined by the following
// Then helper
func (mmPack *mPackagerMockPack) When(o *models.Order) *PackagerMockPackExpectation {
	if mmPack.mock.funcPack != nil {
		mmPack.mock.t.Fatalf("PackagerMock.Pack mock is already set by Set")
	}

	expectation := &PackagerMockPackExpectation{
		mock:               mmPack.mock,
		params:             &PackagerMockPackParams{o},
		expectationOrigins: PackagerMockPackExpectationOrigins{origin: minimock.CallerInfo(1)},
	}
	mmPack.expectations = append(mmPack.expectations, expectation)
	return expectation
}

// Then sets up Packager.Pack return parameters for the expectation previously defined by the When method
func (e *PackagerMockPackExpectation) Then(err error) *PackagerMock {
	e.results = &PackagerMockPackResults{err}
	return e.mock
}

// Times sets number of times Packager.Pack should be invoked
func (mmPack *mPackagerMockPack) Times(n uint64) *mPackagerMockPack {
	if n == 0 {
		mmPack.mock.t.Fatalf("Times of PackagerMock.Pack mock can not be zero")
	}
	mm_atomic.StoreUint64(&mmPack.expectedInvocations, n)
	mmPack.expectedInvocationsOrigin = minimock.CallerInfo(1)
	return mmPack
}

func (mmPack *mPackagerMockPack) invocationsDone() bool {
	if len(mmPack.expectations) == 0 && mmPack.defaultExpectation == nil && mmPack.mock.funcPack == nil {
		return true
	}

	totalInvocations := mm_atomic.LoadUint64(&mmPack.mock.afterPackCounter)
	expectedInvocations := mm_atomic.LoadUint64(&mmPack.expectedInvocations)

	return totalInvocations > 0 && (expectedInvocations == 0 || expectedInvocations == totalInvocations)
}

// Pack implements mm_packing.Packager
func (mmPack *PackagerMock) Pack(o *models.Order) (err error) {
	mm_atomic.AddUint64(&mmPack.beforePackCounter, 1)
	defer mm_atomic.AddUint64(&mmPack.afterPackCounter, 1)

	mmPack.t.Helper()

	if mmPack.inspectFuncPack != nil {
		mmPack.inspectFuncPack(o)
	}

	mm_params := PackagerMockPackParams{o}

	// Record call args
	mmPack.PackMock.mutex.Lock()
	mmPack.PackMock.callArgs = append(mmPack.PackMock.callArgs, &mm_params)
	mmPack.PackMock.mutex.Unlock()

	for _, e := range mmPack.PackMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmPack.PackMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmPack.PackMock.defaultExpectation.Counter, 1)
		mm_want := mmPack.PackMock.defaultExpectation.params
		mm_want_ptrs := mmPack.PackMock.defaultExpectation.paramPtrs

		mm_got := PackagerMockPackParams{o}

		if mm_want_ptrs != nil {

			if mm_want_ptrs.o != nil && !minimock.Equal(*mm_want_ptrs.o, mm_got.o) {
				mmPack.t.Errorf("PackagerMock.Pack got unexpected parameter o, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
					mmPack.PackMock.defaultExpectation.expectationOrigins.originO, *mm_want_ptrs.o, mm_got.o, minimock.Diff(*mm_want_ptrs.o, mm_got.o))
			}

		} else if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmPack.t.Errorf("PackagerMock.Pack got unexpected parameters, expected at\n%s:\nwant: %#v\n got: %#v%s\n",
				mmPack.PackMock.defaultExpectation.expectationOrigins.origin, *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmPack.PackMock.defaultExpectation.results
		if mm_results == nil {
			mmPack.t.Fatal("No results are set for the PackagerMock.Pack")
		}
		return (*mm_results).err
	}
	if mmPack.funcPack != nil {
		return mmPack.funcPack(o)
	}
	mmPack.t.Fatalf("Unexpected call to PackagerMock.Pack. %v", o)
	return
}

// PackAfterCounter returnStorage a count of finished PackagerMock.Pack invocations
func (mmPack *PackagerMock) PackAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPack.afterPackCounter)
}

// PackBeforeCounter returnStorage a count of PackagerMock.Pack invocations
func (mmPack *PackagerMock) PackBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmPack.beforePackCounter)
}

// Calls returnStorage a list of arguments used in each call to PackagerMock.Pack.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmPack *mPackagerMockPack) Calls() []*PackagerMockPackParams {
	mmPack.mutex.RLock()

	argCopy := make([]*PackagerMockPackParams, len(mmPack.callArgs))
	copy(argCopy, mmPack.callArgs)

	mmPack.mutex.RUnlock()

	return argCopy
}

// MinimockPackDone returnStorage true if the count of the Pack invocations corresponds
// the number of defined expectations
func (m *PackagerMock) MinimockPackDone() bool {
	if m.PackMock.optional {
		// Optional methods provide '0 or more' call count restriction.
		return true
	}

	for _, e := range m.PackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	return m.PackMock.invocationsDone()
}

// MinimockPackInspect logs each unmet expectation
func (m *PackagerMock) MinimockPackInspect() {
	for _, e := range m.PackMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to PackagerMock.Pack at\n%s with params: %#v", e.expectationOrigins.origin, *e.params)
		}
	}

	afterPackCounter := mm_atomic.LoadUint64(&m.afterPackCounter)
	// if default expectation was set then invocations count should be greater than zero
	if m.PackMock.defaultExpectation != nil && afterPackCounter < 1 {
		if m.PackMock.defaultExpectation.params == nil {
			m.t.Errorf("Expected call to PackagerMock.Pack at\n%s", m.PackMock.defaultExpectation.returnOrigin)
		} else {
			m.t.Errorf("Expected call to PackagerMock.Pack at\n%s with params: %#v", m.PackMock.defaultExpectation.expectationOrigins.origin, *m.PackMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcPack != nil && afterPackCounter < 1 {
		m.t.Errorf("Expected call to PackagerMock.Pack at\n%s", m.funcPackOrigin)
	}

	if !m.PackMock.invocationsDone() && afterPackCounter > 0 {
		m.t.Errorf("Expected %d calls to PackagerMock.Pack at\n%s but found %d calls",
			mm_atomic.LoadUint64(&m.PackMock.expectedInvocations), m.PackMock.expectedInvocationsOrigin, afterPackCounter)
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *PackagerMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockPackInspect()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *PackagerMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *PackagerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockPackDone()
}
