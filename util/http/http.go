// Description: http
// Author: ZHU HAIHUA
// Since: 2016-03-02 20:46
package httputil

import (
	"io"
	"bytes"
	"io/ioutil"
	"net/http"
	"errors"
)

// One of the copies, say from b to r2, could be avoided by using a more
// elaborate trick where the other copy is made during Request/Response.Write.
// This would complicate things too much, given that these functions are for
// debugging only.
func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}

func DumpBodyAsReader(req *http.Request) (reader io.ReadCloser, err error) {
	if req == nil || req.Body == nil {
		return nil, errors.New("request or body is nil")
	} else {
		reader, req.Body, err = drainBody(req.Body)
	}
	return
}

func DumpBodyAsBytes(req *http.Request) (copy []byte, err error) {
	var reader io.ReadCloser
	reader, err = DumpBodyAsReader(req)
	copy, err = ioutil.ReadAll(reader)
	return
}