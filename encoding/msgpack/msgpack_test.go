package msgpack

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRoundTrip(t *testing.T) {
	type payload struct {
		ID    int64  `msgpack:"id"`
		Name  string `msgpack:"name"`
		Flags []bool `msgpack:"flags"`
	}

	want := payload{ID: 42, Name: "due", Flags: []bool{true, false, true}}
	encoded, err := Marshal(want)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	const wantWire = "83a269642aa46e616d65a3647565a5666c61677393c3c2c3"
	if got := fmt.Sprintf("%x", encoded); got != wantWire {
		t.Fatalf("wire bytes = %s, want compatibility bytes %s", got, wantWire)
	}
	var got payload
	if err := Unmarshal(encoded, &got); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip = %#v, want %#v", got, want)
	}
}

func TestUnmarshalRejectsTruncatedFixExtWithoutPanic(t *testing.T) {
	for _, code := range []byte{0xd4, 0xd5, 0xd6, 0xd7, 0xd8} {
		code := code
		t.Run(fmt.Sprintf("0x%02x", code), func(t *testing.T) {
			// A fixext value requires an extension type byte and a fixed-size
			// payload. The historical v2 decoder could panic when the payload
			// was truncated after the type byte (CVE-2026-32284).
			var value any
			if err := Unmarshal([]byte{code, 0}, &value); err == nil {
				t.Fatal("Unmarshal truncated fixext returned nil error")
			}
		})
	}
}
