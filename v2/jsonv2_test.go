//go:build jsonv2

package json

import (
	"bytes"
	stdjson "encoding/json"
	"encoding/json/jsontext"
	"io"
	"strings"
	"testing"
	"time"
)

func TestJSONV2Package(t *testing.T) {
	if Package != "encoding/json/v2" {
		t.Fatalf("Package = %q, want %q", Package, "encoding/json/v2")
	}
	if API == nil {
		t.Fatal("API is nil")
	}
}

// TestJSONV2MapKeyOrderIsDeterministic fails if marshalOpts drops
// Deterministic: encoding/json/v2 would then emit map members in Go's
// randomized iteration order.
func TestJSONV2MapKeyOrderIsDeterministic(t *testing.T) {
	m := map[string]int{"z": 1, "a": 2, "m": 3, "b": 4, "y": 5, "c": 6, "x": 7, "d": 8}
	want := `{"a":2,"b":4,"c":6,"d":8,"m":3,"x":7,"y":5,"z":1}`

	for i := range 300 {
		got, err := Marshal(m)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != want {
			t.Fatalf("Marshal() iteration %d = %s, want sorted keys %s", i, got, want)
		}

		s, err := MarshalToString(m)
		if err != nil {
			t.Fatalf("MarshalToString() error = %v", err)
		}
		if s != want {
			t.Fatalf("MarshalToString() iteration %d = %s, want sorted keys %s", i, s, want)
		}
	}

	for i := range 300 {
		got, err := MarshalIndent(m, "", "")
		if err != nil {
			t.Fatalf("MarshalIndent() error = %v", err)
		}
		if wantIndented := mustIndent(t, want, "", ""); string(got) != wantIndented {
			t.Fatalf("MarshalIndent() iteration %d = %q, want %q", i, got, wantIndented)
		}
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	for range 300 {
		buf.Reset()
		if err := enc.Encode(m); err != nil {
			t.Fatalf("Encode() error = %v", err)
		}
		if buf.String() != want+"\n" {
			t.Fatalf("Encode() = %q, want %q", buf.String(), want+"\n")
		}
	}
}

// TestJSONV2MarshalIndentArbitraryPrefix fails if MarshalIndent is routed
// through jsontext.WithIndentPrefix, which panics on anything other than
// spaces and tabs, and it pins the output to the default backend byte for byte.
func TestJSONV2MarshalIndentArbitraryPrefix(t *testing.T) {
	type inner struct {
		Z []int          `json:"z"`
		M map[string]int `json:"m"`
	}
	type outer struct {
		Name  string  `json:"name"`
		In    inner   `json:"in"`
		Score float64 `json:"score"`
		Ptr   *int    `json:"ptr"`
	}
	n := 5
	values := []any{
		outer{Name: "zhangsan", In: inner{Z: []int{1, 2}, M: map[string]int{"k": 1}}, Score: 1.5, Ptr: &n},
		map[string]any{"empty": map[string]any{}, "arr": []any{}},
		[]any{1, "two", nil, true},
		struct{}{},
	}
	indents := [][2]string{
		{"", "  "},
		{"", "\t"},
		{"  ", "\t"},
		{"", ""},
		{"--", "  "},
		{"\n", " "},
		{"prefix", "indent"},
	}

	for _, v := range values {
		for _, pair := range indents {
			prefix, indent := pair[0], pair[1]
			got, err := MarshalIndent(v, prefix, indent)
			if err != nil {
				t.Fatalf("MarshalIndent(%#v, %q, %q) error = %v", v, prefix, indent, err)
			}
			want, err := stdjson.MarshalIndent(v, prefix, indent)
			if err != nil {
				t.Fatalf("stdjson.MarshalIndent(%#v, %q, %q) error = %v", v, prefix, indent, err)
			}
			if string(got) != string(want) {
				t.Fatalf("MarshalIndent(%#v, %q, %q):\n got = %q\nwant = %q", v, prefix, indent, got, want)
			}
		}
	}
}

// TestJSONV2MarshalIndentEmptyArgsStayMultiline fails if MarshalIndent skips
// indentation when both arguments are empty. encoding/json still emits
// multiline output in that case; only Encoder.SetIndent("", "") disables it.
func TestJSONV2MarshalIndentEmptyArgsStayMultiline(t *testing.T) {
	got, err := MarshalIndent(map[string]int{"a": 1, "b": 2}, "", "")
	if err != nil {
		t.Fatalf("MarshalIndent() error = %v", err)
	}
	if !strings.Contains(string(got), "\n") {
		t.Fatalf("MarshalIndent(v, \"\", \"\") = %q, want multiline output", got)
	}
	want, err := stdjson.MarshalIndent(map[string]int{"a": 1, "b": 2}, "", "")
	if err != nil {
		t.Fatalf("stdjson.MarshalIndent() error = %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("MarshalIndent(v, \"\", \"\") = %q, want %q", got, want)
	}
}

func TestJSONV2EncoderIndent(t *testing.T) {
	value := map[string]int{"a": 1, "b": 2}

	// SetIndent("", "") disables indentation, unlike MarshalIndent(v, "", "").
	var compact bytes.Buffer
	enc := NewEncoder(&compact)
	enc.SetIndent("", "")
	if err := enc.Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if compact.String() != `{"a":1,"b":2}`+"\n" {
		t.Fatalf("Encode() after SetIndent(\"\", \"\") = %q, want compact output", compact.String())
	}

	// An arbitrary prefix must not panic and must match the default backend.
	var indented bytes.Buffer
	enc = NewEncoder(&indented)
	enc.SetIndent("--", "  ")
	if err := enc.Encode(value); err != nil {
		t.Fatalf("Encode() after SetIndent(\"--\", \"  \") error = %v", err)
	}
	want, err := stdjson.MarshalIndent(value, "--", "  ")
	if err != nil {
		t.Fatalf("stdjson.MarshalIndent() error = %v", err)
	}
	if indented.String() != string(want)+"\n" {
		t.Fatalf("Encode() = %q, want %q", indented.String(), string(want)+"\n")
	}
}

func TestJSONV2EncoderEscapeHTML(t *testing.T) {
	value := map[string]string{"tag": "<script>&</script>"}

	// This backend defaults to not escaping, matching its own Marshal.
	var def bytes.Buffer
	if err := NewEncoder(&def).Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if strings.Contains(def.String(), `\u003c`) {
		t.Fatalf("Encode() with default settings = %q, want unescaped HTML", def.String())
	}

	var on bytes.Buffer
	enc := NewEncoder(&on)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if !strings.Contains(on.String(), `\u003c`) {
		t.Fatalf("Encode() after SetEscapeHTML(true) = %q, want escaped HTML", on.String())
	}

	var off bytes.Buffer
	enc = NewEncoder(&off)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if strings.Contains(off.String(), `\u003c`) {
		t.Fatalf("Encode() after SetEscapeHTML(false) = %q, want unescaped HTML", off.String())
	}

	// Toggling back to true on the same Encoder must take effect too.
	enc.SetEscapeHTML(true)
	off.Reset()
	if err := enc.Encode(value); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if !strings.Contains(off.String(), `\u003c`) {
		t.Fatalf("Encode() after toggling SetEscapeHTML back to true = %q, want escaped HTML", off.String())
	}
}

func TestJSONV2EncoderSuccessiveEncodes(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(map[string]int{"a": 1}); err != nil {
		t.Fatalf("Encode() first error = %v", err)
	}
	if err := enc.Encode(map[string]int{"b": 2}); err != nil {
		t.Fatalf("Encode() second error = %v", err)
	}
	if got := buf.String(); got != "{\"a\":1}\n{\"b\":2}\n" {
		t.Fatalf("two Encode() calls = %q, want each value newline-terminated", got)
	}
}

func TestJSONV2EncoderWriteError(t *testing.T) {
	enc := NewEncoder(errWriter{})
	if err := enc.Encode(map[string]int{"a": 1}); err == nil {
		t.Fatal("Encode() to a failing writer returned nil error")
	}
}

func TestJSONV2MarshalErrorPaths(t *testing.T) {
	// A func value has no JSON representation, which is the reachable way to
	// make every marshal entry point fail.
	unsupported := func() {}

	if _, err := Marshal(unsupported); err == nil {
		t.Fatal("Marshal() of an unsupported type returned nil error")
	}
	if s, err := MarshalToString(unsupported); err == nil {
		t.Fatalf("MarshalToString() of an unsupported type returned %q, nil", s)
	} else if s != "" {
		t.Fatalf("MarshalToString() returned %q alongside an error, want empty string", s)
	}
	if b, err := MarshalIndent(unsupported, "", "  "); err == nil {
		t.Fatalf("MarshalIndent() of an unsupported type returned %s, nil", b)
	} else if b != nil {
		t.Fatalf("MarshalIndent() returned %s alongside an error, want nil", b)
	}

	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	if err := enc.Encode(unsupported); err == nil {
		t.Fatal("Encode() of an unsupported type returned nil error")
	}
	if buf.Len() != 0 {
		t.Fatalf("Encode() wrote %q despite failing to marshal", buf.String())
	}

	// The same path with indentation enabled, so the failure happens before
	// indentJSON rather than inside it.
	buf.Reset()
	enc = NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(unsupported); err == nil {
		t.Fatal("Encode() with indentation of an unsupported type returned nil error")
	}
	if buf.Len() != 0 {
		t.Fatalf("Encode() wrote %q despite failing to marshal", buf.String())
	}
}

type errWriter struct{}

func (errWriter) Write([]byte) (int, error) { return 0, io.ErrClosedPipe }

// TestJSONV2UseNumberIsRecursive fails if UseNumber only intercepts the
// top-level value: encoding/json/v2 has no exported option for this, so the
// interception has to reach every any inside maps and slices too.
func TestJSONV2UseNumberIsRecursive(t *testing.T) {
	const input = `{"n":9007199254740993,"nested":{"m":1.5},"arr":[1,2],"s":"x","b":true,"nil":null}`

	dec := NewDecoder(strings.NewReader(input))
	dec.UseNumber()
	var got map[string]any
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	top, ok := got["n"].(Number)
	if !ok {
		t.Fatalf("top-level number type = %T, want Number", got["n"])
	}
	if top.String() != "9007199254740993" {
		t.Fatalf("top-level number = %q, want %q without float64 rounding", top, "9007199254740993")
	}

	nested, ok := got["nested"].(map[string]any)
	if !ok {
		t.Fatalf("nested object type = %T, want map[string]any", got["nested"])
	}
	nestedNum, ok := nested["m"].(Number)
	if !ok {
		t.Fatalf("nested number type = %T, want Number", nested["m"])
	}
	if nestedNum.String() != "1.5" {
		t.Fatalf("nested number = %q, want %q", nestedNum, "1.5")
	}

	arr, ok := got["arr"].([]any)
	if !ok {
		t.Fatalf("array type = %T, want []any", got["arr"])
	}
	for i, elem := range arr {
		if _, ok := elem.(Number); !ok {
			t.Fatalf("array element %d type = %T, want Number", i, elem)
		}
	}

	// Non-numeric kinds must fall through to the default v2 behavior.
	if _, ok := got["s"].(string); !ok {
		t.Fatalf("string type = %T, want string", got["s"])
	}
	if _, ok := got["b"].(bool); !ok {
		t.Fatalf("bool type = %T, want bool", got["b"])
	}
	if got["nil"] != nil {
		t.Fatalf("null = %#v, want nil", got["nil"])
	}
}

func TestJSONV2DecoderWithoutUseNumber(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"n":1.5}`))
	var got map[string]any
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if _, ok := got["n"].(float64); !ok {
		t.Fatalf("number type without UseNumber = %T, want float64", got["n"])
	}
}

// TestJSONV2DecoderOptionsApplyToLaterDecodes fails if the options are baked
// into the jsontext.Decoder at construction instead of being passed to each
// UnmarshalDecode call.
func TestJSONV2DecoderOptionsApplyToLaterDecodes(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"known":1} {"known":1,"extra":2}`))
	var first struct {
		Known int `json:"known"`
	}
	if err := dec.Decode(&first); err != nil {
		t.Fatalf("Decode() first error = %v", err)
	}

	dec.DisallowUnknownFields()
	var second struct {
		Known int `json:"known"`
	}
	if err := dec.Decode(&second); err == nil {
		t.Fatal("Decode() after DisallowUnknownFields accepted an unknown member")
	}
}

func TestJSONV2DecoderRepeatedOptionCalls(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"known":1,"extra":2}`))
	for range 1000 {
		dec.UseNumber()
		dec.DisallowUnknownFields()
	}
	var dst struct {
		Known int `json:"known"`
	}
	if err := dec.Decode(&dst); err == nil {
		t.Fatal("Decode() with DisallowUnknownFields accepted an unknown member")
	}
}

func TestJSONV2DecoderMoreInsideContainers(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{}`))
	if !dec.More() {
		t.Fatal("More() before the only value = false, want true")
	}
	var obj map[string]any
	if err := dec.Decode(&obj); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if dec.More() {
		t.Fatal("More() after the only value = true, want false")
	}

	empty := NewDecoder(strings.NewReader(``))
	if empty.More() {
		t.Fatal("More() on empty input = true, want false")
	}
	var v int
	if err := empty.Decode(&v); err == nil {
		t.Fatal("Decode() on empty input returned nil error")
	}
}

// TestJSONV2DecoderSyntaxErrorIsSticky pins that a malformed value is never
// silently skipped: once Decode fails on syntax, later calls keep failing and
// More reports that nothing more can be read, so a for-More loop terminates.
func TestJSONV2DecoderSyntaxErrorIsSticky(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"a":  } {"a":1}`))
	var v map[string]any

	if err := dec.Decode(&v); err == nil {
		t.Fatal("Decode() of malformed input returned nil error")
	}
	for i := range 3 {
		if err := dec.Decode(&v); err == nil {
			t.Fatalf("Decode() call %d after a syntax error returned nil error", i+2)
		}
	}
	if dec.More() {
		t.Fatal("More() after a syntax error = true, want false so that a for-More loop terminates")
	}
}

func TestJSONV2EncoderIndentEdgeCases(t *testing.T) {
	value := map[string]int{"a": 1, "b": 2}
	// Only ("", "") disables indentation; a non-empty prefix alone still indents,
	// matching the default backend.
	cases := []struct {
		prefix, indent string
		wantIndented   bool
	}{
		{" ", "", true},
		{"", "\t", true},
		{"\t", "", true},
		{"", "", false},
	}
	for _, c := range cases {
		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		enc.SetIndent(c.prefix, c.indent)
		if err := enc.Encode(value); err != nil {
			t.Fatalf("Encode() after SetIndent(%q, %q) error = %v", c.prefix, c.indent, err)
		}
		if gotIndented := strings.Contains(buf.String(), "\n \"") || strings.Contains(buf.String(), "\n\t\""); gotIndented != c.wantIndented {
			t.Fatalf("Encode() after SetIndent(%q, %q) = %q, indented = %v, want %v",
				c.prefix, c.indent, buf.String(), gotIndented, c.wantIndented)
		}

		want, err := stdjson.MarshalIndent(value, c.prefix, c.indent)
		if !c.wantIndented {
			want, err = Marshal(value)
		}
		if err != nil {
			t.Fatalf("building expected output error = %v", err)
		}
		if buf.String() != string(want)+"\n" {
			t.Fatalf("Encode() after SetIndent(%q, %q) = %q, want %q", c.prefix, c.indent, buf.String(), string(want)+"\n")
		}
	}
}

func TestJSONV2EncoderEscapeHTMLWithMapOrdering(t *testing.T) {
	value := map[string]string{"z": "<z>", "a": "<a>", "m": "<m>"}
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.SetEscapeHTML(true)
	for range 50 {
		buf.Reset()
		if err := enc.Encode(value); err != nil {
			t.Fatalf("Encode() error = %v", err)
		}
		want := `{"a":"\u003ca\u003e","m":"\u003cm\u003e","z":"\u003cz\u003e"}` + "\n"
		if buf.String() != want {
			t.Fatalf("Encode() = %q, want %q", buf.String(), want)
		}
	}
}

func TestJSONV2DecoderBufferedBeforeDecode(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"a":1}`))
	b, err := io.ReadAll(dec.Buffered())
	if err != nil {
		t.Fatalf("ReadAll(Buffered()) before Decode error = %v", err)
	}
	if len(b) != 0 {
		t.Fatalf("Buffered() before Decode = %q, want empty", b)
	}
}

// TestJSONV2Duration fails if marshalOpts or unmarshalOpts drops
// FormatDurationAsNano: encoding/json/v2 alone reports "no default
// representation" for time.Duration, so every entry point below would fail
// outright while the other four backends keep working.
func TestJSONV2Duration(t *testing.T) {
	type withDuration struct {
		D time.Duration `json:"d"`
	}
	const encoded = `{"d":90000000000}`
	value := withDuration{D: 90 * time.Second}

	got, err := Marshal(value)
	if err != nil {
		t.Fatalf("Marshal() of a time.Duration error = %v", err)
	}
	if string(got) != encoded {
		t.Fatalf("Marshal() = %s, want %s", got, encoded)
	}

	s, err := MarshalToString(value)
	if err != nil {
		t.Fatalf("MarshalToString() of a time.Duration error = %v", err)
	}
	if s != encoded {
		t.Fatalf("MarshalToString() = %s, want %s", s, encoded)
	}

	indented, err := MarshalIndent(value, "", "  ")
	if err != nil {
		t.Fatalf("MarshalIndent() of a time.Duration error = %v", err)
	}
	if !strings.Contains(string(indented), "90000000000") {
		t.Fatalf("MarshalIndent() = %q, want the nanosecond count", indented)
	}

	var buf bytes.Buffer
	if err := NewEncoder(&buf).Encode(value); err != nil {
		t.Fatalf("Encode() of a time.Duration error = %v", err)
	}
	if buf.String() != encoded+"\n" {
		t.Fatalf("Encode() = %q, want %q", buf.String(), encoded+"\n")
	}

	var viaUnmarshal withDuration
	if err := Unmarshal([]byte(encoded), &viaUnmarshal); err != nil {
		t.Fatalf("Unmarshal() into a time.Duration error = %v", err)
	}
	if viaUnmarshal.D != 90*time.Second {
		t.Fatalf("Unmarshal() = %v, want %v", viaUnmarshal.D, 90*time.Second)
	}

	var viaDecode withDuration
	if err := NewDecoder(strings.NewReader(encoded)).Decode(&viaDecode); err != nil {
		t.Fatalf("Decode() into a time.Duration error = %v", err)
	}
	if viaDecode.D != 90*time.Second {
		t.Fatalf("Decode() = %v, want %v", viaDecode.D, 90*time.Second)
	}
}

// TestJSONV2OmitEmptyLegacySemantics fails if marshalOpts drops
// OmitEmptyWithLegacySemantics: v2 alone only omits values that encode as
// null, "", {} or [], so zero numbers and false would start showing up.
func TestJSONV2OmitEmpty(t *testing.T) {
	type omitEmpty struct {
		S  string         `json:"s,omitempty"`
		N  int            `json:"n,omitempty"`
		U  uint           `json:"u,omitempty"`
		B  bool           `json:"b,omitempty"`
		F  float64        `json:"f,omitempty"`
		Sl []string       `json:"sl,omitempty"`
		M  map[string]int `json:"m,omitempty"`
		P  *int           `json:"p,omitempty"`
		I  any            `json:"i,omitempty"`
	}

	got, err := Marshal(omitEmpty{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(got) != "{}" {
		t.Fatalf("Marshal() of an all-zero struct = %s, want {}", got)
	}

	var buf bytes.Buffer
	if err = NewEncoder(&buf).Encode(omitEmpty{}); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if buf.String() != "{}\n" {
		t.Fatalf("Encode() of an all-zero struct = %q, want %q", buf.String(), "{}\n")
	}

	// Non-zero values must still be present.
	n := 7
	got, err = Marshal(omitEmpty{N: 1, B: true, F: 1.5, P: &n})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	for _, want := range []string{`"n":1`, `"b":true`, `"f":1.5`, `"p":7`} {
		if !strings.Contains(string(got), want) {
			t.Fatalf("Marshal() = %s, want it to contain %s", got, want)
		}
	}
}

// TestJSONV2StandardTagSupport pins the jsonv2 column of the README matrix.
//
// These capabilities come from the shared v2 engine rather than from anything
// this backend implements, and no option switches them off: field tags take
// precedence over the equivalent Options, so neither
// OmitZeroStructFields(false) nor MatchCaseInsensitiveNames(true) makes these
// assertions fail. That makes this a regression sentinel — against toolchain
// changes and against the backend being rewired to something other than
// encoding/json/v2 — not proof of a contract this package upholds on its own.
func TestJSONV2StandardTagSupport(t *testing.T) {
	t.Run("omitzero", func(t *testing.T) {
		var v struct {
			NilSlice   []int `json:"ns,omitzero"`
			EmptySlice []int `json:"es,omitzero"`
			Zero       int   `json:"z,omitzero"`
			Keep       int   `json:"keep"`
		}
		v.EmptySlice = []int{}
		v.Keep = 1
		got, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		// omitzero omits the nil slice and the zero int, but keeps the
		// non-nil empty slice: that is the Go-zero-value definition.
		if string(got) != `{"es":[],"keep":1}` {
			t.Fatalf("Marshal() = %s, want omitzero to drop the nil slice and zero int", got)
		}
	})

	t.Run("omitzero 认 IsZero 方法", func(t *testing.T) {
		type wrapper struct {
			C    isZeroAt42 `json:"c,omitzero"`
			Keep int        `json:"keep"`
		}
		got, err := Marshal(wrapper{C: isZeroAt42{N: 42}, Keep: 1})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `{"keep":1}` {
			t.Fatalf("Marshal() = %s, want the field omitted via IsZero()", got)
		}
		got, err = Marshal(wrapper{C: isZeroAt42{N: 7}, Keep: 1})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `{"c":{"N":7},"keep":1}` {
			t.Fatalf("Marshal() = %s, want the field kept when IsZero() is false", got)
		}
	})

	t.Run("case:strict", func(t *testing.T) {
		var v struct {
			Name string `json:"name,case:strict"`
		}
		if err := Unmarshal([]byte(`{"NAME":"x"}`), &v); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if v.Name != "" {
			t.Fatalf("Name = %q, want zero value: case:strict must reject NAME", v.Name)
		}
	})

	t.Run("case:ignore", func(t *testing.T) {
		var v struct {
			Name string `json:"name,case:ignore"`
		}
		if err := Unmarshal([]byte(`{"NAME":"x"}`), &v); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if v.Name != "x" {
			t.Fatalf("Name = %q, want \"x\": case:ignore must accept NAME", v.Name)
		}
	})

	t.Run("embed 兜底字段", func(t *testing.T) {
		type fallback struct {
			A    int        `json:"a"`
			Rest RawMessage `json:",embed"`
		}
		var v fallback
		if err := Unmarshal([]byte(`{"a":1,"zzz":9}`), &v); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if v.A != 1 {
			t.Fatalf("A = %d, want 1", v.A)
		}
		if string(v.Rest) != `{"zzz":9}` {
			t.Fatalf("Rest = %s, want the unmatched member", v.Rest)
		}
		got, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `{"a":1,"zzz":9}` {
			t.Fatalf("Marshal() = %s, want the fallback members merged back", got)
		}
	})

	t.Run("MarshalerTo 与 UnmarshalerFrom", func(t *testing.T) {
		got, err := Marshal(v2Iface{})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `"marshalerTo"` {
			t.Fatalf("Marshal() = %s, want MarshalJSONTo to be called", got)
		}
		var back v2Iface
		if err := Unmarshal([]byte(`"in"`), &back); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if back.V != `unmarshalerFrom:"in"` {
			t.Fatalf("V = %q, want UnmarshalJSONFrom to be called", back.V)
		}
	})
}

type isZeroAt42 struct{ N int }

func (i isZeroAt42) IsZero() bool { return i.N == 42 }

type v2Iface struct{ V string }

func (v2Iface) MarshalJSONTo(enc *jsontext.Encoder) error {
	return enc.WriteValue([]byte(`"marshalerTo"`))
}

func (v *v2Iface) UnmarshalJSONFrom(dec *jsontext.Decoder) error {
	raw, err := dec.ReadValue()
	if err != nil {
		return err
	}
	v.V = "unmarshalerFrom:" + string(raw)
	return nil
}

// TestJSONV2NativeSemantics pins the differences from the default backend that
// this backend deliberately keeps. A failure here means the semantics shifted,
// and the README table needs to change with it.
func TestJSONV2NativeSemantics(t *testing.T) {
	t.Run("HTML 不转义", func(t *testing.T) {
		got, err := Marshal(map[string]string{"tag": "<script>"})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `{"tag":"<script>"}` {
			t.Fatalf("Marshal() = %s, want unescaped HTML", got)
		}
	})

	t.Run("拒绝重复对象键", func(t *testing.T) {
		var dst map[string]int
		if err := Unmarshal([]byte(`{"a":1,"a":2}`), &dst); err == nil {
			t.Fatal("Unmarshal() accepted duplicate object names")
		}
		if Valid([]byte(`{"a":1,"a":2}`)) {
			t.Fatal("Valid() accepted duplicate object names")
		}
	})

	t.Run("nil slice 与 nil map 编码为空容器", func(t *testing.T) {
		got, err := Marshal(struct {
			S []string          `json:"s"`
			M map[string]string `json:"m"`
		}{})
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}
		if string(got) != `{"s":[],"m":{}}` {
			t.Fatalf("Marshal() = %s, want empty containers instead of null", got)
		}
	})

	t.Run("字段名区分大小写", func(t *testing.T) {
		var dst struct {
			Name string
		}
		if err := Unmarshal([]byte(`{"name":"x"}`), &dst); err != nil {
			t.Fatalf("Unmarshal() error = %v", err)
		}
		if dst.Name != "" {
			t.Fatalf("Name = %q, want zero value: v2 matches names case-sensitively", dst.Name)
		}
	})

	t.Run("无效 UTF-8 报错", func(t *testing.T) {
		if _, err := Marshal(map[string]string{"k": "\xff"}); err == nil {
			t.Fatal("Marshal() accepted invalid UTF-8")
		}
	})
}

func TestJSONV2EdgeInputs(t *testing.T) {
	if got, err := Marshal(nil); err != nil || string(got) != "null" {
		t.Fatalf("Marshal(nil) = %s, %v, want null, nil", got, err)
	}
	if err := Unmarshal([]byte(`{}`), nil); err == nil {
		t.Fatal("Unmarshal() into nil returned nil error")
	}
	if Valid(nil) {
		t.Fatal("Valid(nil) = true, want false")
	}
	if !Valid([]byte(`  {"a":1}  `)) {
		t.Fatal("Valid() rejected JSON surrounded by whitespace")
	}
	if Valid([]byte(`{"a":1} trailing`)) {
		t.Fatal("Valid() accepted trailing data")
	}
	if got, err := MarshalToString(nil); err != nil || got != "null" {
		t.Fatalf("MarshalToString(nil) = %q, %v, want \"null\", nil", got, err)
	}
	if _, err := MarshalIndent(func() {}, "", "  "); err == nil {
		t.Fatal("MarshalIndent() of an unsupported type returned nil error")
	}
}

func mustIndent(t *testing.T, src, prefix, indent string) string {
	t.Helper()
	var buf bytes.Buffer
	if err := stdjson.Indent(&buf, []byte(src), prefix, indent); err != nil {
		t.Fatalf("stdjson.Indent() error = %v", err)
	}
	return buf.String()
}
