package main

import (
	"errors"
	"reflect"
	"testing"
)

func Test_generateChar(t *testing.T) {
	// arrange
	// act
	actualChar := generateChar(0)

	// assert
	if actualChar == 0 {
		t.Errorf("expected not 0 but it got %d", actualChar)
	}
}

func Test_display(t *testing.T) {
	// arrange
	data := "ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"
	tcs := []struct {
		expectedResult interface{}
		format         string
	}{
		{
			expectedResult: data,
			format:         "text",
		},
		{
			expectedResult: "4142434445464748494a4b4c4d4e4f505152535455565758595a5f6162636465666768696a6b6c6d6e6f707172737475767778797a",
			format:         "hex",
		},
		{
			expectedResult: "QUJDREVGR0hJSktMTU5PUFFSU1RVVldYWVpfYWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXo=",
			format:         "base64",
		},
		{
			expectedResult: `{"length":53,"result":"ABCDEFGHIJKLMNOPQRSTUVWXYZ_abcdefghijklmnopqrstuvwxyz"}`,
			format:         "json",
		},
		{
			expectedResult: []byte(data),
			format:         "bytes",
		},
		{
			expectedResult: errors.New("unknown result format"),
			format:         "unknown",
		},
	}

	for _, tc := range tcs {
		// act
		actualResult := display(tc.format, []byte(data))

		// assert
		if !reflect.DeepEqual(tc.expectedResult, actualResult) {
			t.Errorf("expected %v but it got %v", tc.expectedResult, actualResult)
		}
	}
}

func Benchmark_Generate_10000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateBytes(10000)
	}
}

func Benchmark_Generate_1000(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateBytes(1000)
	}
}

func Benchmark_Generate_100(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateBytes(100)
	}
}

func Benchmark_Generate_10(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateBytes(10)
	}
}

func Benchmark_Generate_1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		generateBytes(1)
	}
}
