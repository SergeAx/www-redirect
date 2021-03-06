package main

import (
	"fmt"
	"net/http"
	"strings"
)

const serverSoftware = "www-redirect"
const serverVersion = "1.0.0"

// isDomainName checks if a string is a presentation-format domain name
// (currently restricted to hostname-compatible "preferred name" LDH labels and
// SRV-like "underscore labels"; see golang.org/issue/12421).
func isDomainName(s string) bool {
	// See RFC 1035, RFC 3696.
	// Presentation format has dots before every label except the first, and the
	// terminal empty label is optional here because we assume fully-qualified
	// (absolute) input. We must therefore reserve space for the first and last
	// labels' length octets in wire format, where they are necessary and the
	// maximum total length is 255.
	// So our _effective_ maximum is 253, but 254 is not rejected if the last
	// character is a dot.
	l := len(s)
	if l == 0 || l > 254 || l == 254 && s[l-1] != '.' {
		return false
	}

	last := byte('.')
	ok := false // Ok once we've seen a letter.
	partlen := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		default:
			return false
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_':
			ok = true
			partlen++
		case '0' <= c && c <= '9':
			// fine
			partlen++
		case c == '-':
			// Byte before dash cannot be dot.
			if last == '.' {
				return false
			}
			partlen++
		case c == '.':
			// Byte before dot cannot be dot, dash.
			if last == '.' || last == '-' {
				return false
			}
			if partlen > 63 || partlen == 0 {
				return false
			}
			partlen = 0
		}
		last = c
	}
	if last == '-' || partlen > 63 {
		return false
	}

	return ok
}

func transformDomain(name string) (string, int) {
	name = strings.ToLower(name)
	if !isDomainName(name) {
		return "", http.StatusBadRequest
	}
	if len(name) >= 4 && name[0:4] == "www." {
		return "", http.StatusNotFound
	} else {
		return "www." + name, 0
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Server", fmt.Sprintf("%s (%s)", serverSoftware, serverVersion))
	domain := r.Header.Get("Host")
	response, httpError := transformDomain(domain)
	if httpError != 0 {
		http.Error(w, "", httpError)
	} else {
		http.Redirect(w, r, "http://"+response+r.URL.Path, http.StatusMovedPermanently)
		return
	}
}

func main() {
	http.HandleFunc("/", redirectHandler)
	http.ListenAndServe(":80", nil)
}
