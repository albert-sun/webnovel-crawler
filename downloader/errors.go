package downloader

import errors "github.com/cockroachdb/errors"

var RequestError = errors.New("error performing request")
var ParseHTMLError = errors.New("error parsing html for goquery")
var ValidationError = errors.New("error during validation") // Some sort of validation error
var GoQueryError = errors.New("missing expected during querying")
