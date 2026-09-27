package binary

import (
	"errors"
	"testing"
)

// ============================================================================
// Marshal Tests - Integer Types Only
// ============================================================================

func TestMarshal_Int_PositiveValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int // expected output length in bytes
	}{
		{"int positive", int(100), 8},
		{"int8 positive", int8(50), 1},
		{"int16 positive", int16(300), 2},
		{"int32 positive", int32(10000), 4},
		{"int64 positive", int64(1000000), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("Marshal() length = %d; want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestMarshal_Int_ZeroValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"int zero", int(0), 8},
		{"int8 zero", int8(0), 1},
		{"int16 zero", int16(0), 2},
		{"int32 zero", int32(0), 4},
		{"int64 zero", int64(0), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("Marshal() length = %d; want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestMarshal_Int_NegativeValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"int -1", int(-1), 8},
		{"int8 -50", int8(-50), 1},
		{"int16 -32768", int16(-32768), 2},
		{"int32 -2147483648", int32(-2147483648), 4},
		{"int64 -9223372036854775808", int64(-9223372036854775808), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("Marshal() length = %d; want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestMarshal_Int_MaxMinValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"int max", int(^uint(0) >> 1), 8},
		{"int min", int64(minInt64), 8},
		{"int8 max", int8(maxInt8), 1},
		{"int8 min", int8(minInt8), 1},
		{"int16 max", int16(maxInt16), 2},
		{"int16 min", int16(minInt16), 2},
		{"int32 max", int32(maxInt32), 4},
		{"int32 min", int32(minInt32), 4},
		{"int64 max", int64(maxInt64), 8},
		{"int64 min", int64(minInt64), 8},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if len(got) != tt.wantLen {
				t.Errorf("Marshal() length = %d; want %d", len(got), tt.wantLen)
			}
		})
	}
}

// func TestMarshal_UnsupportedTypes(t *testing.T) {
// 	tests := []struct {
// 		name    string
// 		input   any
// 		wantErr bool
// 	}{
// 		{"map", map[string]int{}, true},
// 		{"chan", make(chan int), true},
// 		{"func", func() {}, true},
// 		{"ptr to nil", (*int)(nil), true}, // pointer itself is supported, but nil check happens earlier
// 		{"interface{}", any(nil), true},   // interface with concrete value works
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := Marshal(tt.input)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("Marshal() error = %v; wantErr = %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// ============================================================================
// Unmarshal Tests - Integer Types Only
// ============================================================================

// TestUnmarshal_Int_PositiveValues - Tests positive value preservation across all integer types
func TestUnmarshal_Int_PositiveValues(t *testing.T) {
	t.Run("int positive", func(t *testing.T) {
		var v int = 100
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 100 {
			t.Errorf("Expected 100, got %d", v)
		}
	})

	t.Run("int8 positive", func(t *testing.T) {
		var v int8 = 50
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 50 {
			t.Errorf("Expected 50, got %d", v)
		}
	})

	t.Run("int16 positive", func(t *testing.T) {
		var v int16 = 300
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 300 {
			t.Errorf("Expected 300, got %d", v)
		}
	})

	t.Run("int32 positive", func(t *testing.T) {
		var v int32 = 10000
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 10000 {
			t.Errorf("Expected 10000, got %d", v)
		}
	})

	t.Run("int64 positive", func(t *testing.T) {
		var v int64 = 1000000
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 1000000 {
			t.Errorf("Expected 1000000, got %d", v)
		}
	})
}

func TestUnmarshal_Int_NegativeValues(t *testing.T) {
	t.Run("int -1", func(t *testing.T) {
		var v int = -1
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -1 {
			t.Errorf("Expected -1, got %d", v)
		}
	})

	t.Run("int8 -50", func(t *testing.T) {
		var v int8 = -50
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -50 {
			t.Errorf("Expected -50, got %d", v)
		}
	})

	t.Run("int16 -32768", func(t *testing.T) {
		var v int16 = -32768
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -32768 {
			t.Errorf("Expected -32768, got %d", v)
		}
	})

	t.Run("int32 -2147483648", func(t *testing.T) {
		var v int32 = -2147483648
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -2147483648 {
			t.Errorf("Expected -2147483648, got %d", v)
		}
	})

	t.Run("int64 -9223372036854775808", func(t *testing.T) {
		var v int64 = -9223372036854775808
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -9223372036854775808 {
			t.Errorf("Expected -9223372036854775808, got %d", v)
		}
	})
}

func TestUnmarshal_Int_ZeroValues(t *testing.T) {
	t.Run("int zero", func(t *testing.T) {
		var v int = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %d", v)
		}
	})

	t.Run("int8 zero", func(t *testing.T) {
		var v int8 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %d", v)
		}
	})

	t.Run("int16 zero", func(t *testing.T) {
		var v int16 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %d", v)
		}
	})

	t.Run("int32 zero", func(t *testing.T) {
		var v int32 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %d", v)
		}
	})

	t.Run("int64 zero", func(t *testing.T) {
		var v int64 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %d", v)
		}
	})
}

func TestUnmarshal_Int_MaxMinValues(t *testing.T) {
	t.Run("int max", func(t *testing.T) {
		var v int = 2147483647
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 2147483647 {
			t.Errorf("Expected 2147483647, got %d", v)
		}
	})

	t.Run("int min", func(t *testing.T) {
		var v int = -2147483648
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -2147483648 {
			t.Errorf("Expected -2147483648, got %d", v)
		}
	})

	t.Run("int8 max", func(t *testing.T) {
		var v int8 = 127
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 127 {
			t.Errorf("Expected 127, got %d", v)
		}
	})

	t.Run("int8 min", func(t *testing.T) {
		var v int8 = -128
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -128 {
			t.Errorf("Expected -128, got %d", v)
		}
	})

	t.Run("int16 max", func(t *testing.T) {
		var v int16 = 32767
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 32767 {
			t.Errorf("Expected 32767, got %d", v)
		}
	})

	t.Run("int16 min", func(t *testing.T) {
		var v int16 = -32768
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -32768 {
			t.Errorf("Expected -32768, got %d", v)
		}
	})

	t.Run("int32 max", func(t *testing.T) {
		var v int32 = 2147483647
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 2147483647 {
			t.Errorf("Expected 2147483647, got %d", v)
		}
	})

	t.Run("int32 min", func(t *testing.T) {
		var v int32 = -2147483648
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -2147483648 {
			t.Errorf("Expected -2147483648, got %d", v)
		}
	})

	t.Run("int64 max", func(t *testing.T) {
		var v int64 = 9223372036854775807
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 9223372036854775807 {
			t.Errorf("Expected 9223372036854775807, got %d", v)
		}
	})

	t.Run("int64 min", func(t *testing.T) {
		var v int64 = -9223372036854775808
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -9223372036854775808 {
			t.Errorf("Expected -9223372036854775808, got %d", v)
		}
	})
}

func TestUnmarshal_Int_TruncatedData(t *testing.T) {
	t.Run("int truncated", func(t *testing.T) {
		data := make([]byte, 4) // needs 8 bytes for int
		err := Unmarshal(data, (*int)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("int8 truncated", func(t *testing.T) {
		data := make([]byte, 0) // needs 1 byte for int8
		err := Unmarshal(data, (*int8)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("int16 truncated", func(t *testing.T) {
		data := make([]byte, 1) // needs 2 bytes for int16
		err := Unmarshal(data, (*int16)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("int32 truncated", func(t *testing.T) {
		data := make([]byte, 3) // needs 4 bytes for int32
		err := Unmarshal(data, (*int32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("int64 truncated", func(t *testing.T) {
		data := make([]byte, 7) // needs 8 bytes for int64
		err := Unmarshal(data, (*int64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})
}

func TestUnmarshal_Int_ExtraBytes(t *testing.T) {
	t.Run("int extra bytes", func(t *testing.T) {
		data := make([]byte, 9) // needs exactly 8 bytes for int
		err := Unmarshal(data, (*int)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("int8 extra bytes", func(t *testing.T) {
		data := make([]byte, 2) // needs exactly 1 byte for int8
		err := Unmarshal(data, (*int8)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("int16 extra bytes", func(t *testing.T) {
		data := make([]byte, 3) // needs exactly 2 bytes for int16
		err := Unmarshal(data, (*int16)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("int32 extra bytes", func(t *testing.T) {
		data := make([]byte, 5) // needs exactly 4 bytes for int32
		err := Unmarshal(data, (*int32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("int64 extra bytes", func(t *testing.T) {
		data := make([]byte, 9) // needs exactly 8 bytes for int64
		err := Unmarshal(data, (*int64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})
}

func TestUnmarshal_Int_NilDestination(t *testing.T) {
	t.Run("int nil pointer", func(t *testing.T) {
		var dest int = 0
		data := make([]byte, 8)
		err := Unmarshal(data, &dest)

		if err != nil {
			t.Fatalf("Unmarshal() error = %v; want success with zero value", err)
		}

		if dest != 0 {
			t.Errorf("Expected 0, got %d", dest)
		}
	})

	t.Run("int8 nil pointer", func(t *testing.T) {
		var dest int8 = 0
		data := make([]byte, 1)
		err := Unmarshal(data, &dest)

		if err != nil {
			t.Fatalf("Unmarshal() error = %v; want success with zero value", err)
		}

		if dest != 0 {
			t.Errorf("Expected 0, got %d", dest)
		}
	})

	t.Run("int16 nil pointer", func(t *testing.T) {
		var dest int16 = 0
		data := make([]byte, 2)
		err := Unmarshal(data, &dest)

		if err != nil {
			t.Fatalf("Unmarshal() error = %v; want success with zero value", err)
		}

		if dest != 0 {
			t.Errorf("Expected 0, got %d", dest)
		}
	})

	t.Run("int32 nil pointer", func(t *testing.T) {
		var dest int32 = 0
		data := make([]byte, 4)
		err := Unmarshal(data, &dest)

		if err != nil {
			t.Fatalf("Unmarshal() error = %v; want success with zero value", err)
		}

		if dest != 0 {
			t.Errorf("Expected 0, got %d", dest)
		}
	})

	t.Run("int64 nil pointer", func(t *testing.T) {
		var dest int64 = 0
		data := make([]byte, 8)
		err := Unmarshal(data, &dest)

		if err != nil {
			t.Fatalf("Unmarshal() error = %v; want success with zero value", err)
		}

		if dest != 0 {
			t.Errorf("Expected 0, got %d", dest)
		}
	})
}

func TestMarshal_Unmarshal_Int_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"int roundtrip", int(100), 8},   // Add explicit type
		{"uint roundtrip", uint(100), 8}, // Add unsigned types
		{"int8 roundtrip", int8(50), 1},
		{"uint8 roundtrip", uint8(255), 1}, // Add unsigned variants
		{"int16 roundtrip", int16(300), 2},
		{"uint16 roundtrip", uint16(400), 2}, // Add unsigned variants
		{"int32 roundtrip", int32(10000), 4},
		{"uint32 roundtrip", uint32(50000), 4}, // Add unsigned variants
		{"int64 roundtrip", int64(1000000), 8},
		{"uint64 roundtrip", uint64(9223372036854775800), 8}, // Add unsigned variants
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			if len(data) != tt.wantLen {
				t.Errorf("Marshal() length = %d; want %d", len(data), tt.wantLen)
			}

			var dest any
			switch tt.input.(type) {
			case int:
				dest = new(int)
			case uint:
				dest = new(uint)
			case int8:
				dest = new(int8)
			case uint8:
				dest = new(uint8)
			case int16:
				dest = new(int16)
			case uint16:
				dest = new(uint16)
			case int32:
				dest = new(int32)
			case uint32:
				dest = new(uint32)
			case int64:
				dest = new(int64)
			case uint64:
				dest = new(uint64)
			}

			err = Unmarshal(data, dest)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify roundtrip by checking the unmarshaled value matches original
			switch v := tt.input.(type) {
			case int:
				if *dest.(*int) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*int), v)
				}
			case uint:
				if *dest.(*uint) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint), v)
				}
			case int8:
				if *dest.(*int8) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*int8), v)
				}
			case uint8:
				if *dest.(*uint8) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint8), v)
				}
			case int16:
				if *dest.(*int16) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*int16), v)
				}
			case uint16:
				if *dest.(*uint16) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint16), v)
				}
			case int32:
				if *dest.(*int32) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*int32), v)
				}
			case uint32:
				if *dest.(*uint32) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint32), v)
				}
			case int64:
				if *dest.(*int64) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*int64), v)
				}
			case uint64:
				if *dest.(*uint64) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint64), v)
				}
			}
		})
	}
}

// ============================================================================
// Helper Constants for Test Values
// ============================================================================

const (
	maxInt8  = int8(127)
	minInt8  = int8(-128)
	maxInt16 = int16(32767)
	minInt16 = int16(-32768)
	maxInt32 = int32(2147483647)
	minInt32 = int32(-2147483648)
	maxInt64 = int64(9223372036854775807)
	minInt64 = int64(-9223372036854775808)
)

// ============================================================================
// Size Validation Tests (Separate from Type Validation)
// ============================================================================

func TestUnmarshal_Int_SizeMismatch(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		dataLen int
		wantErr bool
	}{
		{"int", (*int)(nil), 4, true},      // 8 bytes needed - truncated data
		{"int8", (*int8)(nil), 0, true},    // 1 byte needed - truncated data
		{"int16", (*int16)(nil), 3, false}, // 2 bytes needed - extra bytes (not truncated)
		{"int32", (*int32)(nil), 7, false}, // 4 bytes needed - extra bytes (not truncated)
		{"int64", (*int64)(nil), 7, true},  // 8 bytes needed - truncated data
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataLen)
			err := Unmarshal(data, tt.input)

			if err == nil && tt.wantErr {
				t.Fatalf("Unmarshal() expected error but got nil")
			}

			switch tt.name {
			case "int", "int8", "int64":
				if !errors.Is(err, ErrTruncatedData) && err != nil {
					t.Errorf("Unmarshal() error = %v; want ErrTruncatedData for insufficient data", err)
				}
			case "int16", "int32":
				if !errors.Is(err, ErrExtraBytes) && err != nil {
					t.Errorf("Unmarshal() error = %v; want ErrExtraBytes for extra bytes", err)
				}
			}
		})
	}
}
