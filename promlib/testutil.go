package promlib

import (
	"bytes"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func dump(g prometheus.Gatherer) string {
	metrics, err := g.Gather()
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	enc := expfmt.NewEncoder(&buf, expfmt.FmtText)
	for _, m := range metrics {
		if err := enc.Encode(m); err != nil {
			panic(err)
		}
	}

	return buf.String()
}
