package hub

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"h12.me/sej"
	"h12.me/sej/hub/proto"
)

func TestMarshal(t *testing.T) {
	req := proto.Request{
		Title: proto.RequestTitle{
			Verb:     uint8(proto.PUT),
			ClientID: "b",
		},
		Header: &proto.Put{
			JournalDir: "c.3.4",
		},
		Messages: []sej.Message{
			{
				Timestamp: time.Now().UTC().Truncate(time.Millisecond),
				Value:     []byte("a"),
			},
		},
	}
	w := new(bytes.Buffer)
	if _, err := req.WriteTo(w); err != nil {
		t.Fatal(err)
	}
	var res proto.Request
	if _, err := res.ReadFrom(bytes.NewReader(w.Bytes())); err != nil {
		t.Fatal(err)
	}
	if expect, actual := js(req), js(res); expect != actual {
		t.Fatalf("expect\n%v\ngot\n%v\n", expect, actual)
	}
}

func js(v interface{}) string {
	buf, _ := json.MarshalIndent(v, "", "    ")
	return string(buf)
}
