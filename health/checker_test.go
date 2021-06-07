package health

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testCase struct {
	testName  string
	callbacks []callback
	expected  *CheckResponse
}

type testOrCase struct {
	testName  string
	callbacks []callback
	expected  []string
}

func runCase(t *testing.T, cas testCase) {
	defaultOpts := CheckerOptions{
		Version:   "42",
		ReleaseID: "42",
		ServiceID: "42",
	}
	ch := NewChecker(defaultOpts)

	t.Run(cas.testName, func(t *testing.T) {
		ctx := context.Background()
		for _, cb := range cas.callbacks {
			ch.AddCallback(cb.name, cb.cb)
		}
		res := ch.CheckContext(ctx)
		if !reflect.DeepEqual(res, cas.expected) {
			t.Errorf("got %+v, expected %+v", res, cas.expected)
		}
	})
}

func TestConcurrentCase(t *testing.T) {
	const (
		sleepPerStage = time.Millisecond * 200
		fault         = sleepPerStage / 2
	)

	testData := testOrCase{
		testName: "any fail output response",
		callbacks: []callback{
			{
				cb: CheckCallback(func(ctx context.Context) *CheckResult {
					time.Sleep(sleepPerStage)
					return &CheckResult{
						Error: &FailError{
							Message: "failure",
						},
					}
				}),
				name: "fail",
			},
			{
				cb: CheckCallback(func(ctx context.Context) *CheckResult {
					time.Sleep(sleepPerStage)
					return &CheckResult{
						Error: &FailError{
							Message: "failure1",
						},
					}
				}),
				name: "fail1",
			},
			{
				cb: CheckCallback(func(ctx context.Context) *CheckResult {
					time.Sleep(sleepPerStage)
					return &CheckResult{
						Error: &FailError{
							Message: "failure2",
						},
					}
				}),
				name: "fail2",
			},
		},
		expected: []string{"failure2", "failure", "failure1"},
	}

	defaultOpts := CheckerOptions{
		Version:   "42",
		ReleaseID: "42",
		ServiceID: "42",
	}
	ch := NewChecker(defaultOpts)

	t.Run("Concurrent run", func(t *testing.T) {
		start := time.Now()
		for _, cb := range testData.callbacks {
			ch.AddCallback(cb.name, cb.cb)
		}
		_ = ch.Check()
		elapsed := time.Since(start)
		require.Less(t, int64(elapsed), (int64(sleepPerStage)*int64(len(testData.callbacks)))-int64(fault))
	})

	t.Run("Any fail output response", func(t *testing.T) {
		for _, cb := range testData.callbacks {
			ch.AddCallback(cb.name, cb.cb)
		}
		res := ch.Check()
		require.Subset(t, testData.expected, []string{res.Output})
	})
}

func TestCheckerAddMultipleCallbacks(t *testing.T) {
	var ch Checker
	defaultOpts := CheckerOptions{
		Version:   "42",
		ReleaseID: "42",
		ServiceID: "42",
	}
	ch1 := &checker{
		opts: defaultOpts,
	}

	ch = ch1

	ch.AddMultipleCallbacks([]string{
		"a",
		"b",
		"c",
	}, []CheckCallback{
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
	})

	if len(ch1.callbacks) != 3 {
		t.Errorf("got callbacks len %d, expected 3", len(ch1.callbacks))
	}
}

func TestCheckerAddMultipleCallbacksPanicking(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic, but got <nil>")
		}
	}()

	var ch Checker

	defaultOpts := CheckerOptions{
		Version:   "42",
		ReleaseID: "42",
		ServiceID: "42",
	}

	ch1 := &checker{
		opts: defaultOpts,
	}
	ch = ch1

	ch.AddMultipleCallbacks([]string{
		"a",
		"b",
	}, []CheckCallback{
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
		CheckCallback(func(ctx context.Context) *CheckResult {
			return nil
		}),
	})
}

func TestChecker(t *testing.T) {
	var testData = []testCase{
		{
			testName: "valid check with empty callbacks",
			expected: &CheckResponse{
				Status:    CheckStatusPass,
				Version:   "42",
				ReleaseID: "42",
				ServiceID: "42",
				Checks:    make(map[string]*CheckResult),
			},
		},
		{
			testName: "valid check with nil-returning callback",
			callbacks: []callback{
				{
					name: "",
					cb: CheckCallback(func(ctx context.Context) *CheckResult {
						return nil
					}),
				},
			},
			expected: &CheckResponse{
				Status:    CheckStatusPass,
				Version:   "42",
				ReleaseID: "42",
				ServiceID: "42",
				Checks:    make(map[string]*CheckResult),
			},
		},
		{
			testName: "panicking callback",
			callbacks: []callback{
				{
					name: "panicking callback",
					cb: CheckCallback(func(ctx context.Context) *CheckResult {
						panic("error message")
					}),
				},
			},
			expected: &CheckResponse{
				Status:    CheckStatusFail,
				Version:   "42",
				ReleaseID: "42",
				ServiceID: "42",
				Checks: map[string]*CheckResult{
					"panicking callback": {
						Output: "error message",
						Status: CheckStatusFail,
					},
				},
				Output: "error message",
			},
		},
		{
			testName: "valid check with non-empty callback",
			callbacks: []callback{
				{
					name: "name",
					cb: CheckCallback(func(ctx context.Context) *CheckResult {
						return &CheckResult{
							ObservedUnit:  "ms",
							ObservedValue: 42,
						}
					}),
				},
			},
			expected: &CheckResponse{
				Status: CheckStatusPass,
				Checks: map[string]*CheckResult{
					"name": {
						ObservedUnit:  "ms",
						ObservedValue: 42,
					},
				},
				ReleaseID: "42",
				ServiceID: "42",
				Version:   "42",
			},
		},
		{
			testName: "warning callback",
			callbacks: []callback{
				{
					cb: CheckCallback(func(ctx context.Context) *CheckResult {
						return &CheckResult{
							Error: &WarnError{
								Message: "warn",
							},
						}
					}),
					name: "name",
				},
			},
			expected: &CheckResponse{
				Status: CheckStatusPass,
				Checks: map[string]*CheckResult{
					"name": {
						Status: CheckStatusWarn,
						Output: "warn",
						Error: &WarnError{
							Message: "warn",
						},
					},
				},
				ReleaseID: "42",
				ServiceID: "42",
				Version:   "42",
			},
		},
	}

	for _, v := range testData {
		v := v
		runCase(t, v)
	}
}
