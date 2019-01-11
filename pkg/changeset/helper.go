package changeset

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
)

// md5sum will return a checksum from a string.
func md5sum(s string) string {
	h := md5.New()
	r := bytes.NewReader([]byte(s))
	_, _ = io.Copy(h, r)
	return fmt.Sprintf("%x", h.Sum(nil))
}
