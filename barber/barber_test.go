package barber

import (
	"testing"
	"time"

	"github.com/city-mobil/gobuns/config"
)

func TestHostFails(t *testing.T) {
	cfg := &Config{
		Threshold: 1,
		MaxFails:  20,
	}
	barb := NewBarber([]int{0}, cfg)

	var testData = []struct {
		hostname          int
		addCount          int
		testTime          time.Time
		expectedAvailable bool
		testName          string
	}{
		{
			hostname:          0,
			addCount:          30,
			testTime:          time.Now(),
			expectedAvailable: false,
			testName:          "existing_hostname",
		},
		{
			hostname:          42,
			addCount:          1,
			testTime:          time.Now(),
			expectedAvailable: false,
			testName:          "non_existing_hostname",
		},
		{
			hostname:          0,
			addCount:          0,
			testTime:          time.Now().Add(time.Second),
			expectedAvailable: true,
			testName:          "existing_hostname_after_cooldown",
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			for i := 0; i < v.addCount; i++ {
				barb.AddError(v.hostname, v.testTime)
			}
			isAvailable := barb.IsAvailable(v.hostname, v.testTime)
			if v.expectedAvailable != isAvailable {
				t.Errorf("want %v, got %v", v.expectedAvailable, isAvailable)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	cfg := NewConfig("")
	_ = config.InitOnce()

	conf := cfg()
	if conf.Threshold != defaultThreshold {
		t.Errorf("expected threshold %v, got %v", defaultThreshold, conf.Threshold)
	}
}

func TestNilConfig(t *testing.T) {
	var cfg *Config
	cfg = cfg.withDefaults()
	if cfg.Threshold != defaultThreshold {
		t.Errorf("expected threshold %v, got %v", defaultThreshold, cfg.Threshold)
	}
}

func TestConfig(t *testing.T) {
	var testData = []struct {
		threshold         uint32
		expectedThreshold uint32
		testName          string
	}{
		{
			testName:          "valid config",
			threshold:         defaultThreshold,
			expectedThreshold: defaultThreshold,
		},
		{
			testName:          "invalid threshold",
			expectedThreshold: defaultThreshold,
		},
	}

	for _, v := range testData {
		v := v
		t.Run(v.testName, func(t *testing.T) {
			cfg := &Config{
				Threshold: v.threshold,
			}
			cfg = cfg.withDefaults()
			if cfg.Threshold != v.expectedThreshold {
				t.Errorf("expected threshold %v, got %v", v.expectedThreshold, cfg.Threshold)
			}
		})
	}
}

func TestStats(t *testing.T) {
	barber := NewBarber([]int{1}, &Config{
		Threshold: 42,
	})
	barber.AddError(1, time.Now())
	barber.AddError(1, time.Now())
	barber.AddError(1, time.Now())

	st := barber.Stats()
	if len(st.Hosts) != 1 {
		t.Errorf("got hosts len %d, expected 1", len(st.Hosts))
	}

	h := st.Hosts[0]
	if h.ServerID != 1 {
		t.Errorf("got server id %d, expected 1", h.ServerID)
	}
	if h.FailsCount != 3 {
		t.Errorf("got fails count %d, expected 2", h.FailsCount)
	}
}

func BenchmarkBarberAddError(b *testing.B) {
	// NOTE(a.petrukhin): circuit_breaker package has 4000 ns/op and 700B/op
	// This barber is faster on inserts, but little slowlier for select because of the
	// absence of the background goroutine.
	barber := NewBarber([]int{1}, &Config{
		Threshold: 42,
	})
	tm := time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				barber.IsAvailable(1, tm)
			}
		}()
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			barber.AddError(1, tm)
		}
	})
}

func BenchmarkBarberIsAvailable(b *testing.B) {
	// NOTE(a.petrukhin): circuit_breaker package has 4000 ns/op and 700B/op per insert
	// This barber is faster on inserts, but 100ns slowlier for select because of the
	// absence of the background goroutine.
	barber := NewBarber([]int{1}, &Config{
		Threshold: 42,
	})
	tm := time.Now()
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				time.Sleep(10 * time.Millisecond)
				barber.AddError(1, tm)
			}
		}()
	}
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			barber.IsAvailable(1, tm)
		}
	})
}
