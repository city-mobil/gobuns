package health

import (
	"context"
	"errors"
	"testing"
	"time"

	mock_health "github.com/city-mobil/gobuns/mocks/health"
	"github.com/golang/mock/gomock"
)

func testCheckResultEqual(got, expected *CheckResult) bool {
	if got.ComponentID != expected.ComponentID {
		return false
	}
	if got.ComponentType != expected.ComponentType {
		return false
	}
	if got.ObservedUnit != expected.ObservedUnit {
		return false
	}
	if got.ObservedValue.(int64) >= expected.ObservedValue.(int64) {
		return false
	}
	if got.ObservedValue.(int64)+1 >= expected.ObservedValue.(int64) {
		return false
	}
	if got.Output != expected.Output {
		return false
	}

	return got.Status == expected.Status
}

const (
	testComponentType = "test_component_type"
)

func testCheckables(ctrl *gomock.Controller) (okCheckable, longResponseCheckable, erroringCheckable Checkable) {
	ctx := context.Background()
	okCheckableM := mock_health.NewMockCheckable(ctrl)
	okCheckableM.EXPECT().Ping(ctx).Return(nil).Do(func(ctx context.Context) {
		time.Sleep(4 * time.Millisecond)
	}).AnyTimes()
	okCheckableM.EXPECT().Name().Return("ok").AnyTimes()
	okCheckableM.EXPECT().ComponentID().Return("ok").AnyTimes()
	okCheckableM.EXPECT().ComponentType().Return(testComponentType).AnyTimes()

	longResponseCheckableM := mock_health.NewMockCheckable(ctrl)
	longResponseCheckableM.EXPECT().Ping(ctx).Return(nil).Do(func(ctx context.Context) {
		time.Sleep(healthCriticalResponseTime)
	}).AnyTimes()
	longResponseCheckableM.EXPECT().Name().Return("long_response").AnyTimes()
	longResponseCheckableM.EXPECT().ComponentID().Return("long_response").AnyTimes()
	longResponseCheckableM.EXPECT().ComponentType().Return(testComponentType).AnyTimes()

	erroringCheckableM := mock_health.NewMockCheckable(ctrl)
	erroringCheckableM.EXPECT().Ping(ctx).Return(errors.New("got some error")).AnyTimes()
	erroringCheckableM.EXPECT().Name().Return("erroring").AnyTimes()
	erroringCheckableM.EXPECT().ComponentID().Return("erroring").AnyTimes()
	erroringCheckableM.EXPECT().ComponentType().Return(testComponentType).AnyTimes()

	okCheckable = okCheckableM
	longResponseCheckable = longResponseCheckableM
	erroringCheckable = erroringCheckableM

	return
}

func TestResponseTimeCallback(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	okCheckable, longResponseCheckable, erroringCheckable := testCheckables(ctrl)

	var testData = []struct {
		conn     Checkable
		expected *CheckResult
		testName string
		isSlave  bool
	}{
		{
			conn: okCheckable,
			expected: &CheckResult{
				Status:        CheckStatusPass,
				ComponentID:   "ok",
				ComponentType: testComponentType,
				ObservedUnit:  healthTimeCheckUnit,
				ObservedValue: int64(10),
			},
			isSlave:  false,
			testName: "passing callback",
		},
		{
			conn: longResponseCheckable,
			expected: &CheckResult{
				Status:        CheckStatusWarn,
				ComponentID:   "long_response",
				ComponentType: testComponentType,
				ObservedUnit:  healthTimeCheckUnit,
				ObservedValue: int64(210),
			},
			isSlave:  false,
			testName: "warning callback: long response time",
		},
		{
			conn: erroringCheckable,
			expected: &CheckResult{
				Status:        CheckStatusWarn,
				ComponentID:   "erroring",
				ComponentType: testComponentType,
				Output:        "got some error",
				ObservedUnit:  healthTimeCheckUnit,
				ObservedValue: int64(10),
			},
			isSlave:  true,
			testName: "erroring_slave",
		},
		{
			conn: erroringCheckable,
			expected: &CheckResult{
				Status:        CheckStatusFail,
				ComponentID:   "erroring",
				ComponentType: testComponentType,
				Output:        "got some error",
				ObservedUnit:  healthTimeCheckUnit,
				ObservedValue: int64(10),
			},
			testName: "erroring_master",
		},
	}

	for _, v := range testData {
		v := v
		ctx := context.Background()
		t.Run(v.testName, func(t *testing.T) {
			cb := NewResponseTimeCheckCallback(v.conn, v.isSlave)
			got := cb(ctx)

			if !testCheckResultEqual(got, v.expected) {
				t.Errorf("got %v, expected %v", got, v.expected)
			}
		})
	}
}
