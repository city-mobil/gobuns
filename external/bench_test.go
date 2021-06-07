package external

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func BenchmarkGetBody(b *testing.B) {
	req, err := http.NewRequest(http.MethodPost, "http://city-mobil.ru", strings.NewReader(strings.Repeat("some_long_string", 4242)))
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := req.GetBody()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func newReader(str []byte) io.ReadCloser {
	return ioutil.NopCloser(bytes.NewReader(str))
}

func BenchmarkNewBody(b *testing.B) {
	body := []byte(strings.Repeat("some_long_string", 4242))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = newReader(body)
	}
}
