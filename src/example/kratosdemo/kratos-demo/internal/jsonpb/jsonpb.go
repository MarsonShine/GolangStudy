package jsonpb

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/pkg/errors"
)

var jsonContentType = []string{"application/json; charset=utf-8"}

// JSON common json struct.
type PBJSON struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	TTL     int           `json:"ttl"`
	Data    proto.Message `json:"data,omitempty"`
}

func (m PBJSON) Reset()         { m = PBJSON{} }
func (m PBJSON) String() string { return proto.CompactTextString(m) }
func (PBJSON) ProtoMessage()    {}
func (pb PBJSON) MarshalJSONPB(*jsonpb.Marshaler) ([]byte, error) {
	return json.Marshal(pb)
}
func (r PBJSON) Render(w http.ResponseWriter) error {
	// FIXME(zhoujiahui): the TTL field will be configurable in the future
	if r.TTL <= 0 {
		r.TTL = 1
	}
	return writePBJSON(w, r)
}

// WriteContentType write json ContentType.
func (r PBJSON) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, jsonContentType)
}

func writePBJSON(w http.ResponseWriter, pb proto.Message) (err error) {
	writeContentType(w, jsonContentType)
	m := jsonpb.Marshaler{}
	var buf bytes.Buffer
	if err = m.Marshal(&buf, pb); err != nil {
		err = errors.WithStack(err)
		return
	}
	if _, err = w.Write(buf.Bytes()); err != nil {
		err = errors.WithStack(err)
	}
	return
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
