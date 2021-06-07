package promlib

import (
	"fmt"
	"sync"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type eventSuite struct {
	suite.Suite
}

func TestEventSuite(t *testing.T) {
	suite.Run(t, new(eventSuite))
}

func (s *eventSuite) SetupTest() {
	reg := prometheus.NewRegistry()
	prometheus.DefaultGatherer = reg
	prometheus.DefaultRegisterer = reg
}

func (s *eventSuite) TestIncCnt() {
	t := s.T()
	eventName := "inc_cnt"
	SetGlobalNamespace("apptest")
	IncCnt(eventName)

	fqn := fmt.Sprintf("%s_%s", globalNS, eventName)
	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 1`, fqn))

	// increment concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			IncCnt(eventName)
			wg.Done()
		}()
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 101`, fqn))
}

func (s *eventSuite) TestIncCntEvent() {
	t := s.T()
	e := &Event{
		Name:      "inc_cnt_event",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	IncCntEvent(e)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 1`, fqn))

	// increment concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			IncCntEvent(e)
			wg.Done()
		}()
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 101`, fqn))
}

func (s *eventSuite) TestIncCntEventWithLabels() {
	t := s.T()
	e := &Event{
		Name:      "inc_cnt_event_labeled",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	IncCntEventWithLabels(e, Labels{
		"ab": "1",
	})

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 1`, fqn))

	// increment concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(iter int) {
			ab := "0"
			if iter%2 == 0 {
				ab = "1"
			}
			IncCntEventWithLabels(e, Labels{
				"ab": ab,
			})
			wg.Done()
		}(i)
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="0"} 50`, fqn))
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 51`, fqn))
}

func (s *eventSuite) TestIncCntEventWithLabels_MultipleLabels() {
	t := s.T()
	e := &Event{
		Name:      "inc_cnt_event_multi_labeled",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	IncCntEventWithLabels(e, Labels{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	})

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{a="a",b="b",c="c",d="d"} 1`, fqn))
}

func (s *eventSuite) TestIncCntEventWithLabels_DiffLabels_ShouldPanic() {
	t := s.T()
	e := &Event{
		Name:      "inc_cnt_event_labeled_panic",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	IncCntEventWithLabels(e, Labels{
		"ab": "1",
	})

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 1`, fqn))

	assert.Panics(t, func() {
		IncCntEventWithLabels(e, Labels{
			"ab":     "1",
			"method": "POST",
		})
	})
}

func (s *eventSuite) TestIncCntWithLabels() {
	t := s.T()
	eventName := "inc_cnt_labeled"
	SetGlobalNamespace("apptest")

	IncCntWithLabels(eventName, Labels{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	})

	fqn := fmt.Sprintf("%s_%s", globalNS, eventName)
	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{a="a",b="b",c="c",d="d"} 1`, fqn))
}

func (s *eventSuite) TestAddCnt() {
	t := s.T()
	eventName := "add_cnt"
	SetGlobalNamespace("test")
	AddCnt(eventName, 2.0)

	fqn := fmt.Sprintf("%s_%s", globalNS, eventName)
	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 2`, fqn))

	// add concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			AddCnt(eventName, 3.0)
			wg.Done()
		}()
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 302`, fqn))
}

func (s *eventSuite) TestAddCntEvent() {
	t := s.T()
	e := &Event{
		Name:      "add_cnt_event",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	AddCntEvent(e, 2.0)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 2`, fqn))

	// add concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func() {
			AddCntEvent(e, 3.0)
			wg.Done()
		}()
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s 302`, fqn))
}

func (s *eventSuite) TestAddCntEventWithLabels() {
	t := s.T()
	e := &Event{
		Name:      "add_cnt_event_labeled",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	AddCntEventWithLabels(e, Labels{
		"ab": "1",
	}, 2.0)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 2`, fqn))

	// increment concurrently
	var wg sync.WaitGroup
	for i := 1; i <= 100; i++ {
		wg.Add(1)
		go func(iter int) {
			ab := "0"
			if iter%2 == 0 {
				ab = "1"
			}
			AddCntEventWithLabels(e, Labels{
				"ab": ab,
			}, 2.0)
			wg.Done()
		}(i)
	}
	wg.Wait()

	metrics = dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="0"} 100`, fqn))
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 102`, fqn))
}

func (s *eventSuite) TestAddCntEventWithLabels_MultipleLabels() {
	t := s.T()
	e := &Event{
		Name:      "add_cnt_event_multi_labeled",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	AddCntEventWithLabels(e, Labels{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	}, 2.0)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{a="a",b="b",c="c",d="d"} 2`, fqn))
}

func (s *eventSuite) TestAddCntEventWithLabels_DiffLabels_ShouldPanic() {
	t := s.T()
	e := &Event{
		Name:      "add_cnt_event_labeled_panic",
		Namespace: "go",
		Subsystem: "test",
	}
	fqn := buildEventFQName(e)

	AddCntEventWithLabels(e, Labels{
		"ab": "1",
	}, 2.0)

	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{ab="1"} 2`, fqn))

	assert.Panics(t, func() {
		AddCntEventWithLabels(e, Labels{
			"ab":     "1",
			"method": "POST",
		}, 1.0)
	})
}

func (s *eventSuite) TestAddCntWithLabels() {
	t := s.T()
	eventName := "add_cnt_labeled"
	SetGlobalNamespace("apptest")

	AddCntWithLabels(eventName, Labels{
		"a": "a",
		"b": "b",
		"c": "c",
		"d": "d",
	}, 2.0)

	fqn := fmt.Sprintf("%s_%s", globalNS, eventName)
	metrics := dump(prometheus.DefaultGatherer)
	assert.Contains(t, metrics, fmt.Sprintf(`%s{a="a",b="b",c="c",d="d"} 2`, fqn))
}
