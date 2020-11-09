package downloader

import errors "github.com/cockroachdb/errors"

var ErrRequest = errors.New("error performing request")
var ErrParseHTML = errors.New("error parsing html for goquery")
var ErrGoQuery = errors.New("missing expected during querying")
