package promlib

import (
	"net/http"
	"strings"
)

type checkFn func(int, string) bool

// URLNumberNormalizer iterates over the URL chunks (split by backslash)
// replacing all only number parts by the given templates.
//
// If no chunks are found, returns the source URL.
//
// If number of the templates is less than number of the found chunks,
// normalizer will use the last template.
//
// 	 req := httptest.NewRequest("GET", "/api/v0/department/100500/user/12", nil)
//	 got := URLNumberNormalizer(req, []string{":did", ":uid"}) # /api/v0/department/:did/user/:uid
//
func URLNumberNormalizer(req *http.Request, templates []string) string {
	return normalizeURL(req.URL.Path, templates, isNumber)
}

// URLHEXNumberNormalizer does the same job as the URLNumberNormalizer method, but in hexadecimal notation.
//
// Note: since the hexadecimal number system includes the characters A, B, C, D, F,
// there should not be lines between the slashes that contain only these characters and numbers,
// because this will be treated as numbers.
//
// 	 req := httptest.NewRequest("GET", "/api/v0/department/18c86399/user/12", nil)
//	 got := URLHEXNumberNormalizer(req, []string{":did", ":uid"}) # /api/v0/department/:did/user/:uid
//
// False replacement:
// 	 req := httptest.NewRequest("GET", "/add/v0/department/18c86399/user/12", nil)
//	 got := URLHEXNumberNormalizer(req, []string{":did", ":uid"}) # /:did/v0/department/:uid/user/:uid
//
func URLHEXNumberNormalizer(req *http.Request, templates []string) string {
	return normalizeURL(req.URL.Path, templates, isHEXNumber)
}

// URLIndicesNormalizer replaces the URL parts with the given indexes.
// Indexing starts from zero.
//
// 	 req := httptest.NewRequest("GET", "/api/v0/department/18c86399/user/12", nil)
//	 got := URLIndicesNormalizer(req, []string{":did", ":uid"}, []int{3, 5}) # /api/v0/department/:did/user/:uid
//
func URLIndicesNormalizer(req *http.Request, templates []string, indices []int) string {
	var checker checkFn
	if len(indices) > 0 {
		checker = func(idx int, _ string) bool {
			for _, v := range indices {
				if v == idx {
					return true
				}
			}
			return false
		}
	}

	return normalizeURL(req.URL.Path, templates, checker)
}

func normalizeURL(path string, templates []string, checker checkFn) string {
	if len(templates) == 0 || checker == nil {
		return path
	}

	var sb strings.Builder
	sb.Grow(len(path))

	from := 0
	left := 0
	tpl := 0
	ind := 0

	for {
		left = indexByteOffset(path, left, '/')
		if left == -1 {
			sb.WriteString(path[from:])
			return sb.String()
		}

		right := indexByteOffset(path, left+1, '/')
		if right == -1 {
			right = len(path)
		}

		if checker(ind, path[left+1:right]) {
			template := templates[tpl]

			sb.WriteString(path[from : left+1])
			sb.WriteString(template)
			from = right

			if tpl+1 < len(templates) {
				tpl++
			}
		}

		left = right
		ind++
	}
}

func indexByteOffset(s string, offset int, c byte) int {
	if offset >= len(s) {
		return -1
	}
	index := strings.IndexByte(s[offset:], c)
	if index == -1 {
		return index
	}
	return offset + index
}

func isNumber(_ int, s string) bool {
	for _, ch := range s {
		if ch < '0' || ch > '9' {
			return false
		}
	}
	return true
}

func isHEXNumber(_ int, s string) bool {
	for _, ch := range s {
		if !((ch >= '0' && ch <= '9') || (ch >= 'a' && ch <= 'f') || (ch >= 'A' && ch <= 'F')) {
			return false
		}
	}
	return true
}
