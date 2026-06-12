package json

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

// testStruct is a representative production payload for benchmarking.
type testStruct struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	Age       int      `json:"age"`
	Active    bool     `json:"active"`
	Score     float64  `json:"score"`
	Tags      []string `json:"tags"`
	Address   address  `json:"address"`
	CreatedAt string   `json:"created_at"`
}

type address struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

var benchData = testStruct{
	ID:     12345,
	Name:   "张三",
	Email:  "zhangsan@example.com",
	Age:    30,
	Active: true,
	Score:  99.5,
	Tags:   []string{"go", "backend", "senior"},
	Address: address{
		Street: "中关村大街1号",
		City:   "北京",
		State:  "北京市",
		Zip:    "100080",
	},
	CreatedAt: "2025-01-01T00:00:00Z",
}

var benchJSON []byte

func init() {
	var err error
	benchJSON, err = Marshal(benchData)
	if err != nil {
		panic("failed to marshal benchmark data: " + err.Error())
	}
}

func TestMarshalUnmarshalRoundTrip(t *testing.T) {
	data, err := Marshal(benchData)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got testStruct
	if err := Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.ID != benchData.ID ||
		got.Name != benchData.Name ||
		got.Email != benchData.Email ||
		got.Age != benchData.Age ||
		got.Active != benchData.Active ||
		got.Score != benchData.Score ||
		got.CreatedAt != benchData.CreatedAt ||
		got.Address != benchData.Address ||
		len(got.Tags) != len(benchData.Tags) {
		t.Fatalf("round trip mismatch: got %+v, want %+v", got, benchData)
	}
	for i := range got.Tags {
		if got.Tags[i] != benchData.Tags[i] {
			t.Fatalf("round trip tags mismatch: got %+v, want %+v", got.Tags, benchData.Tags)
		}
	}
}

func TestMarshalToStringAndValid(t *testing.T) {
	s, err := MarshalToString(benchData)
	if err != nil {
		t.Fatalf("MarshalToString() error = %v", err)
	}
	if s == "" {
		t.Fatal("MarshalToString() returned empty string")
	}
	if !Valid([]byte(s)) {
		t.Fatalf("Valid(%q) = false, want true", s)
	}
	if Valid([]byte(`{"broken":`)) {
		t.Fatal("Valid() accepted malformed JSON")
	}
}

func TestRawMessage(t *testing.T) {
	type envelope struct {
		Data RawMessage `json:"data"`
	}

	src := envelope{Data: RawMessage(`{"nested":true}`)}
	data, err := Marshal(src)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `{"data":{"nested":true}}` {
		t.Fatalf("Marshal() = %s, want raw nested JSON", data)
	}

	var got envelope
	if err := Unmarshal(data, &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if string(got.Data) != `{"nested":true}` {
		t.Fatalf("RawMessage = %s, want nested JSON", got.Data)
	}
}

var (
	_ Marshaler   = customJSON{}
	_ Unmarshaler = (*customJSON)(nil)
)

type customJSON struct {
	Value string
}

func (c customJSON) MarshalJSON() ([]byte, error) {
	return []byte(`"custom"`), nil
}

func (c *customJSON) UnmarshalJSON(data []byte) error {
	c.Value = string(data)
	return nil
}

func TestMarshalerAliases(t *testing.T) {
	data, err := Marshal(customJSON{})
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `"custom"` {
		t.Fatalf("Marshal() = %s, want custom JSON", data)
	}

	var got customJSON
	if err := Unmarshal([]byte(`"decoded"`), &got); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got.Value != `"decoded"` {
		t.Fatalf("UnmarshalJSON() stored %q, want raw decoded JSON", got.Value)
	}
}

func TestFormatHelpers(t *testing.T) {
	var compact bytes.Buffer
	if err := Compact(&compact, []byte("{\n  \"a\": 1\n}")); err != nil {
		t.Fatalf("Compact() error = %v", err)
	}
	if compact.String() != `{"a":1}` {
		t.Fatalf("Compact() = %q, want compact JSON", compact.String())
	}

	var indented bytes.Buffer
	if err := Indent(&indented, compact.Bytes(), "", "  "); err != nil {
		t.Fatalf("Indent() error = %v", err)
	}
	if !strings.Contains(indented.String(), "\n  ") {
		t.Fatalf("Indent() = %q, want indented JSON", indented.String())
	}

	var escaped bytes.Buffer
	HTMLEscape(&escaped, []byte(`{"tag":"<script>"}`))
	if !strings.Contains(escaped.String(), `\u003cscript\u003e`) {
		t.Fatalf("HTMLEscape() = %q, want escaped HTML characters", escaped.String())
	}
}

func TestEncoderDecoder(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(struct {
		Value string `json:"value"`
	}{Value: "<tag>"}); err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
	if strings.Contains(buf.String(), `\u003c`) {
		t.Fatalf("Encode() escaped HTML despite SetEscapeHTML(false): %q", buf.String())
	}
	if !strings.Contains(buf.String(), "\n  ") {
		t.Fatalf("Encode() did not indent despite SetIndent(): %q", buf.String())
	}

	dec := NewDecoder(strings.NewReader(buf.String()))
	var got struct {
		Value string `json:"value"`
	}
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got.Value != "<tag>" {
		t.Fatalf("Decode() value = %q, want %q", got.Value, "<tag>")
	}
}

func TestDecoderBuffered(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"value":1} trailing`))
	var got struct {
		Value int `json:"value"`
	}
	if err := dec.Decode(&got); err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got.Value != 1 {
		t.Fatalf("Decode() value = %d, want 1", got.Value)
	}

	remaining, err := io.ReadAll(dec.Buffered())
	if err != nil {
		t.Fatalf("ReadAll(Buffered()) error = %v", err)
	}
	if !strings.Contains(string(remaining), "trailing") {
		t.Fatalf("Buffered() = %q, want trailing data", remaining)
	}
}

func TestDecoderOptions(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`{"known": 1, "extra": 2}`))
	dec.DisallowUnknownFields()
	var dst struct {
		Known int `json:"known"`
	}
	if err := dec.Decode(&dst); err == nil {
		t.Fatal("Decode() with DisallowUnknownFields accepted an unknown field")
	}

	dec = NewDecoder(strings.NewReader(`{"n": 9007199254740993}`))
	dec.UseNumber()
	var numbers map[string]any
	if err := dec.Decode(&numbers); err != nil {
		t.Fatalf("Decode() with UseNumber error = %v", err)
	}
	gotNumber, ok := numbers["n"].(Number)
	if !ok {
		t.Fatalf("UseNumber decoded type = %T, want json.Number", numbers["n"])
	}
	if got := gotNumber.String(); got != "9007199254740993" {
		t.Fatalf("UseNumber decoded value = %q, want %q", got, "9007199254740993")
	}
}

func TestDecoderMore(t *testing.T) {
	dec := NewDecoder(strings.NewReader(`1 2`))
	if !dec.More() {
		t.Fatal("More() before first value = false, want true")
	}

	var first int
	if err := dec.Decode(&first); err != nil {
		t.Fatalf("Decode() first value error = %v", err)
	}
	if first != 1 {
		t.Fatalf("first value = %d, want 1", first)
	}
	if !dec.More() {
		t.Fatal("More() before second value = false, want true")
	}

	var second int
	if err := dec.Decode(&second); err != nil {
		t.Fatalf("Decode() second value error = %v", err)
	}
	if second != 2 {
		t.Fatalf("second value = %d, want 2", second)
	}
	if dec.More() {
		t.Fatal("More() after final value = true, want false")
	}
}

func TestTopLevelFunctionsDelegateToAPI(t *testing.T) {
	original := API
	t.Cleanup(func() {
		API = original
	})

	API = mockCore{}

	data, err := Marshal("ignored")
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}
	if string(data) != `{"mock":true}` {
		t.Fatalf("Marshal() = %s, want mock payload", data)
	}

	unmarshalErr := Unmarshal([]byte(`{}`), nil)
	if !errors.Is(unmarshalErr, errMockUnmarshal) {
		t.Fatalf("Unmarshal() error = %v, want %v", unmarshalErr, errMockUnmarshal)
	}

	s, err := MarshalToString("ignored")
	if err != nil {
		t.Fatalf("MarshalToString() error = %v", err)
	}
	if s != `{"mock":"string"}` {
		t.Fatalf("MarshalToString() = %q, want mock string", s)
	}

	if Valid(nil) {
		t.Fatal("Valid() = true, want false from mock")
	}
}

var errMockUnmarshal = errors.New("mock unmarshal")

type mockCore struct{}

func (mockCore) Marshal(_ any) ([]byte, error) {
	return []byte(`{"mock":true}`), nil
}

func (mockCore) Unmarshal(_ []byte, _ any) error {
	return errMockUnmarshal
}

func (mockCore) MarshalIndent(_ any, _, _ string) ([]byte, error) {
	return []byte("{\n  \"mock\": true\n}"), nil
}

func (mockCore) MarshalToString(_ any) (string, error) {
	return `{"mock":"string"}`, nil
}

func (mockCore) NewEncoder(_ io.Writer) Encoder {
	return nil
}

func (mockCore) NewDecoder(_ io.Reader) Decoder {
	return nil
}

func (mockCore) Valid(_ []byte) bool {
	return false
}

// BenchmarkMarshal measures Marshal throughput via the API interface.
func BenchmarkMarshal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := Marshal(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkUnmarshal measures Unmarshal throughput via the API interface.
func BenchmarkUnmarshal(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		var v testStruct
		if err := Unmarshal(benchJSON, &v); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMarshalToString measures MarshalToString throughput.
func BenchmarkMarshalToString(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := MarshalToString(benchData)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkValid measures JSON validation throughput.
func BenchmarkValid(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		if !Valid(benchJSON) {
			b.Fatal("expected valid JSON")
		}
	}
}

// BenchmarkMarshalSmall tests with a tiny payload to highlight interface overhead.
func BenchmarkMarshalSmall(b *testing.B) {
	small := struct {
		ID int `json:"id"`
	}{ID: 1}
	b.ReportAllocs()
	for b.Loop() {
		_, err := Marshal(small)
		if err != nil {
			b.Fatal(err)
		}
	}
}
