package json

import (
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
