package promlib

var (
	PostgreSQLQueryTime = &HistogramOpts{
		MetaOpts: MetaOpts{
			Name:      "query_time",
			Subsystem: "postgresql",
			Help:      "PostgreSQL queries response time",
		},
		Labels:  []string{"query"},
		Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
	}

	MySQLQueryTime = &HistogramOpts{
		MetaOpts: MetaOpts{
			Name:      "query_time",
			Subsystem: "mysql",
			Help:      "MySQL queries response time",
		},
		Labels:  []string{"query"},
		Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
	}

	TarantoolOpTime = &HistogramOpts{
		MetaOpts: MetaOpts{
			Name:      "operation_time",
			Subsystem: "tarantool",
			Help:      "Tarantool operations response time",
		},
		Labels:  []string{"collection", "operation"},
		Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
	}

	RedisOpTime = &HistogramOpts{
		MetaOpts: MetaOpts{
			Name:      "operation_time",
			Subsystem: "redis",
			Help:      "Redis operations response time",
		},
		Labels:  []string{"operation"},
		Buckets: []float64{.002, .005, .01, .015, .025, .05, .1, .25, .5, 1, 2, 10},
	}
)
