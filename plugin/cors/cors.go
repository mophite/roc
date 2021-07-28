package cors

import (
	"net/http"
)

const (
	corsAllowOriginHeader      string = "Access-Control-Allow-Origin"
	corsExposeHeadersHeader    string = "Access-Control-Expose-Headers"
	corsMaxAgeHeader           string = "Access-Control-Max-Age"
	corsAllowMethodsHeader     string = "Access-Control-Allow-Methods"
	corsAllowHeadersHeader     string = "Access-Control-Allow-Headers"
	corsAllowCredentialsHeader string = "Access-Control-Allow-Credentials"
	corsRequestMethodHeader    string = "Access-Control-Request-Method"
	corsRequestHeadersHeader   string = "Access-Control-Request-Headers"
	corsOriginHeader           string = "Origin"
	corsVaryHeader             string = "Vary"
	corsOriginMatchAll         string = "*"
)

func Cors(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set(corsAllowOriginHeader, corsOriginMatchAll)
	w.Header().Add(corsAllowMethodsHeader, "POST,DELETE")
	w.Header().Add(corsAllowHeadersHeader,
		"Origin, Content-Length, Content-Type")
	return nil
}
