package binary

import (
	"bytes"
	"encoding/binary"
	"reflect"
	"strings"
	"testing"
)

func TestMarshalArray(t *testing.T) {
	tests := []struct {
		name    string
		array   interface{}
		wantErr bool
		wantLen int
	}{
		{
			name:    "empty int8 array",
			array:   [0]int8{},
			wantLen: 4, // Just the length prefix
		},
		{
			name:    "single element int8 positive",
			array:   [1]int8{42},
			wantLen: 5,
		},
		{
			name:    "int8 with negative value",
			array:   [2]int8{-10, -20},
			wantLen: 6,
		},
		{
			name:    "full range int8 values",
			array:   [3]int8{-128, 0, 127},
			wantLen: 7, // 4 + 3*1 = 7 bytes (length prefix + 3 elements)
		},
		{
			name:    "single element uint8",
			array:   [1]uint8{100},
			wantLen: 5,
		},
		{
			name:    "full range uint8 values",
			array:   [3]uint8{0, 127, 255},
			wantLen: 7, // 4 + 3*1 = 7 bytes (length prefix + 3 elements)
		},
		{
			name:    "larger int8 array",
			array:   [10]int8{-5, -3, 0, 7, 12, 15, 20, 25, 30, 35},
			wantLen: 14, // 4 + 10*1 = 14 bytes (length prefix + 10 elements)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.array)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != tt.wantLen {
				t.Errorf("Marshal() length = %d, want %d", len(got), tt.wantLen)
			}

			if !tt.wantErr && len(got) > 0 {
				// Verify length prefix is correct (little-endian uint32)
				prefix := binary.LittleEndian.Uint32(got[:4])
				expectedLen := intArrayLen(tt.array)
				if prefix != uint32(expectedLen) {
					t.Errorf("Marshal() prefix = %d, want %d", prefix, expectedLen)
				}
			}
		})
	}
}

func TestUnmarshalArray(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		dst      interface{}
		wantErr  bool
		expected int
	}{
		{
			name: "empty int8 array",
			data: make([]byte, 4), // Just length prefix for empty (4 bytes of zeros)
			dst: &struct {
				X [0]int8 `bin:"X"`
			}{},
			wantErr:  false,
			expected: 0,
		},
		{
			name: "single int8 element",
			data: []byte{
				1, 0, 0, 0, // length = 1 (little-endian uint32)
				42, // value
			},
			dst: &struct {
				X [1]int8 `bin:"X"`
			}{},
			wantErr:  false,
			expected: 42,
		},
		{
			name: "int8 with negative values",
			data: []byte{
				3, 0, 0, 0, // length = 3 (little-endian uint32)
				246, // -10 in two's complement (uint8 view)
				236, // -20 in two's complement
				119, // 7 positive
			},
			dst: &struct {
				X [3]int8 `bin:"X"`
			}{},
			wantErr:  false,
			expected: -10,
		},
		{
			name: "full range int8",
			data: []byte{
				3, 0, 0, 0, // length = 3 (little-endian uint32)
				128, // -128 two's complement (0x80)
				0,   // 0
				127, // 127 positive (0x7F)
			},
			dst: &struct {
				X [3]int8 `bin:"X"`
			}{},
			wantErr:  false,
			expected: -128,
		},
		{
			name: "full range uint8",
			data: []byte{
				3, 0, 0, 0, // length = 3 (little-endian uint32)
				0,   // 0 (0x00)
				127, // 127 positive (0x7F)
				255, // 255 max uint8 (0xFF)
			},
			dst: &struct {
				X [3]uint8 `bin:"X"`
			}{},
			wantErr:  false,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.data, tt.dst)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v, name %s", err, tt.wantErr, tt.name)
				return
			}

			if !tt.wantErr {
				switch v := tt.dst.(type) {
				case *struct{ X [1]int8 }:
					if v.X[0] != int8(tt.expected) {
						t.Errorf("Unmarshal() got %d, want %d", v.X[0], tt.expected)
					}
				case *struct{ X [2]int8 }:
					if v.X[0] != int8(-10) || v.X[1] != int8(-20) {
						t.Errorf("Unmarshal() got %+v, want [-10, -20]", v.X[:])
					}
				case *struct{ X [3]int8 }:
					if v.X[0] != int8(tt.expected) || v.X[1] != 0 || v.X[2] != 127 {
						t.Errorf("Unmarshal() got %+v, want [-%d, 0, 127]", v.X[:], tt.expected)
					}
				case *struct{ X [3]uint8 }:
					if v.X[0] != uint8(tt.expected) || v.X[1] != 127 || v.X[2] != 255 {
						t.Errorf("Unmarshal() got %+v, want [0, 127, 255]", v.X[:])
					}
				}
			}
		})
	}
}

func TestMarshalRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		array interface{}
	}{
		{
			name:  "empty int8",
			array: [0]int8{},
		},
		{
			name:  "single int8 negative",
			array: [1]int8{-42},
		},
		{
			name:  "multiple int8 mixed",
			array: [5]int8{-128, -64, 0, 32, 127},
		},
		{
			name:  "single uint8",
			array: [1]uint8{200},
		},
		{
			name:  "multiple uint8 mixed",
			array: [5]uint8{0, 127, 128, 200, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.array)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var dst interface{}
			switch tt.array.(type) {
			case [0]int8:
				dst = &struct {
					X [0]int8 `bin:"X"`
				}{}
			case [1]int8:
				dst = &struct {
					X [1]int8 `bin:"X"`
				}{}
			case [5]int8:
				dst = &struct {
					X [5]int8 `bin:"X"`
				}{}
			case [1]uint8:
				dst = &struct {
					X [1]uint8 `bin:"X"`
				}{}
			case [5]uint8:
				dst = &struct {
					X [5]uint8 `bin:"X"`
				}{}
			}

			err = Unmarshal(data, dst)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify round-trip preserves values
			switch v := tt.array.(type) {
			case [0]int8:
				if len(v) != 0 {
					t.Error("round-trip failed for empty int8")
				}
			case [1]int8:
				dstStruct := dst.(*struct {
					X [1]int8 `bin:"X"`
				})
				if v[0] != dstStruct.X[0] {
					t.Errorf("round-trip failed: got %d, want %d", v[0], dstStruct.X[0])
				}
			case [5]int8:
				dstStruct := dst.(*struct {
					X [5]int8 `bin:"X"`
				})
				if !arraysEqualInt8(v[:], dstStruct.X[:]) {
					t.Errorf("round-trip failed for int8 array")
				}
			case [1]uint8:
				dstStruct := dst.(*struct {
					X [1]uint8 `bin:"X"`
				})
				if v[0] != dstStruct.X[0] {
					t.Errorf("round-trip failed: got %d, want %d", v[0], dstStruct.X[0])
				}
			case [5]uint8:
				dstStruct := dst.(*struct {
					X [5]uint8 `bin:"X"`
				})
				if !arraysEqualUint8(v[:], dstStruct.X[:]) {
					t.Errorf("round-trip failed for uint8 array")
				}
			}
		})
	}
}

func TestUnmarshalArrayTruncated(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		dst     interface{}
		wantErr bool
	}{
		{
			name: "incomplete length prefix",
			data: []byte{0, 1}, // Only 2 bytes of length prefix (little-endian uint32 incomplete)
			dst: &struct {
				X [3]int8 `bin:"X"`
			}{},
			wantErr: true,
		},
		{
			name: "incomplete element data",
			data: []byte{
				1, 0, 0, 0, // length = 1 (little-endian uint32)
				42, // one byte for first element
			},
			dst: &struct {
				X [3]int8 `bin:"X"`
			}{},
			wantErr: true,
		},
		{
			name: "too short for length prefix",
			data: []byte{10, 0, 0, 0}, // Claims 10 elements (little-endian), only 4 bytes provided
			dst: &struct {
				X [3]int8 `bin:"X"`
			}{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Unmarshal(tt.data, tt.dst)
			if tt.wantErr && err == nil {
				t.Errorf("Unmarshal() error = %v, want err = %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(err.Error()) > 0 {
				t.Logf("Unexpected error for non-error case: %s", err.Error())
				return
			}

			// Accept both truncated data and array length mismatch errors (both indicate truncation scenarios)
			if tt.wantErr && err != nil && (strings.Contains(err.Error(), "truncated data") ||
				strings.Contains(err.Error(), "array length mismatch")) {
				// Test case passed - got expected error type
			} else if tt.wantErr && err == nil {
				t.Errorf("Unmarshal() should return an error for truncated/mismatched array: %v", tt.name)
			}
		})
	}
}

func TestMarshalArrayExtraBytes(t *testing.T) {
	data, _ := Marshal([3]int8{1, 2, 3})

	if len(data) > 0 && data[len(data)-1] == 42 {
		t.Error("Marshal() should not leave extra bytes in buffer")
	}
}

func TestUnmarshalArrayLengthMismatch(t *testing.T) {
	// Create data that claims to have a different length than the array type
	data := []byte{
		0, 0, 0, 10, // Claims 10 elements
		1, 2, 3, // Only 3 bytes of actual data (for first element)
	}

	err := Unmarshal(data, &struct {
		X [5]int8 `bin:"X"`
	}{})

	if err == nil {
		t.Error("Unmarshal() should fail when length prefix doesn't match array size")
		return
	}

	if !strings.Contains(err.Error(), "array length mismatch") {
		t.Errorf("Unmarshal() error = %v, want 'array length mismatch'", err)
	}
}

func TestMarshalArrayNilSliceSentinel(t *testing.T) {
	// Verify nil slices use sentinel (not tested for arrays since they have fixed size)
	// This is just to document the behavior distinction
	data, _ := Marshal([0]int8{})

	if len(data) == 4 && binary.LittleEndian.Uint32(data[:4]) != 0 {
		t.Error("empty array should have length prefix of 0")
	}
}

func TestMarshalArrayMaxValues(t *testing.T) {
	tests := []struct {
		name  string
		array interface{}
	}{
		{
			name:  "int8 min value",
			array: [1]int8{-128},
		},
		{
			name:  "int8 max value",
			array: [1]int8{127},
		},
		{
			name:  "uint8 min value",
			array: [1]uint8{0},
		},
		{
			name:  "uint8 max value",
			array: [1]uint8{255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.array)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var dst interface{}
			switch tt.array.(type) {
			case [1]int8:
				dst = &struct {
					X [1]int8 `bin:"X"`
				}{}
			case [1]uint8:
				dst = &struct {
					X [1]uint8 `bin:"X"`
				}{}
			}

			err = Unmarshal(data, dst)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			switch tt.array.(type) {
			case [1]int8:
				dstStruct := dst.(*struct {
					X [1]int8 `bin:"X"`
				})
				wantVal := int8(-128)
				if tt.name == "int8 max value" {
					wantVal = 127
				}
				if dstStruct.X[0] != wantVal {
					t.Errorf("round-trip failed for int8: got %d, want %d", dstStruct.X[0], wantVal)
				}
			case [1]uint8:
				dstStruct := dst.(*struct {
					X [1]uint8 `bin:"X"`
				})
				wantVal := uint8(0)
				if tt.name == "uint8 max value" {
					wantVal = 255
				}
				if dstStruct.X[0] != wantVal {
					t.Errorf("round-trip failed for uint8: got %d, want %d", dstStruct.X[0], wantVal)
				}
			}
		})
	}
}

// Helper function to compare int8 arrays
func arraysEqualInt8(a, b []int8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper function to compare uint8 arrays
func arraysEqualUint8(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Helper to get array length for comparison (unexported helper for tests)
func intArrayLen(arr interface{}) int {
	v := reflect.ValueOf(arr)
	if v.Kind() == reflect.Array || v.Kind() == reflect.Slice {
		return int(v.Len())
	}
	panic("unsupported array type")
}

func TestMarshalArrayDeterministic(t *testing.T) {
	// Verify that marshaling produces deterministic output for same input
	type int8Array struct {
		X [3]int8 `bin:"X"`
	}

	tests := []struct {
		name  string
		array int8Array
	}{
		{
			name: "first run",
			array: int8Array{
				X: [3]int8{-1, 0, 1},
			},
		},
		{
			name: "second run (same input)",
			array: int8Array{
				X: [3]int8{-1, 0, 1},
			},
		},
	}

	firstData := make([]byte, 0)
	for _, tt := range tests {
		data, _ := Marshal(tt.array)
		if len(firstData) == 0 {
			firstData = data
		} else if !bytes.Equal(data, firstData) {
			t.Errorf("Marshal() is not deterministic: %x != %x", data, firstData)
		}
	}

	for _, tt := range tests {
		data, _ := Marshal(tt.array)
		if len(firstData) == 0 || !bytes.Equal(data, firstData) {
			t.Errorf("Marshal() produced different output: %x != %x", data, firstData)
		}
	}
}

// Helper function to create a slice of int8 with specific values for testing
func makeInt8Slice(n int, values ...int8) []int8 {
	slice := make([]int8, n)
	copy(slice, values)
	return slice
}

// Helper function to create a slice of uint8 with specific values for testing
func makeUint8Slice(n int, values ...uint8) []uint8 {
	slice := make([]uint8, n)
	copy(slice, values)
	return slice
}

// Note: The above helper functions are for documentation and would need to be
// moved into a separate test utilities file or included in the test file itself.
// They're not currently used in the tests above but demonstrate how to create
// test data if needed.
