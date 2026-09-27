package binary

import (
	"errors"
	"math"
	"testing"
)

// ============================================================================
// Marshal Tests - Unsigned Integer Types Only
// ============================================================================

func TestMarshal_Uint_PositiveValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int // expected output length in bytes
	}{
		{"uint positive", uint(100), 8},
		{"uint8 positive", uint8(50), 1},
		{"uint16 positive", uint16(300), 2},
		{"uint32 positive", uint32(10000), 4},
		{"uint64 positive", uint64(1000000), 8},
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

func TestMarshal_Uint_ZeroValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"uint zero", uint(0), 8},
		{"uint8 zero", uint8(0), 1},
		{"uint16 zero", uint16(0), 2},
		{"uint32 zero", uint32(0), 4},
		{"uint64 zero", uint64(0), 8},
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

func TestMarshal_Uint_MaxValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"uint max", uint(^uint(0)), 8},
		{"uint8 max", uint8(math.MaxUint8), 1},
		{"uint16 max", uint16(math.MaxUint16), 2},
		{"uint32 max", uint32(math.MaxUint32), 4},
		{"uint64 max", uint64(math.MaxUint64), 8},
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

// ============================================================================
// Unmarshal Tests - Unsigned Integer Types Only
// ============================================================================

func TestUnmarshal_Uint_PositiveValues(t *testing.T) {
	t.Run("uint positive", func(t *testing.T) {
		var v uint = 100
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

	t.Run("uint8 positive", func(t *testing.T) {
		var v uint8 = 50
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

	t.Run("uint16 positive", func(t *testing.T) {
		var v uint16 = 300
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

	t.Run("uint32 positive", func(t *testing.T) {
		var v uint32 = 10000
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

	t.Run("uint64 positive", func(t *testing.T) {
		var v uint64 = 1000000
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

func TestUnmarshal_Uint_ZeroValues(t *testing.T) {
	t.Run("uint zero", func(t *testing.T) {
		var v uint = 0
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

	t.Run("uint8 zero", func(t *testing.T) {
		var v uint8 = 0
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

	t.Run("uint16 zero", func(t *testing.T) {
		var v uint16 = 0
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

	t.Run("uint32 zero", func(t *testing.T) {
		var v uint32 = 0
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

	t.Run("uint64 zero", func(t *testing.T) {
		var v uint64 = 0
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

func TestUnmarshal_Uint_MaxValues(t *testing.T) {
	t.Run("uint max", func(t *testing.T) {
		var v uint = ^uint(0)
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != ^uint(0) {
			t.Errorf("Expected max uint, got %d", v)
		}
	})

	t.Run("uint8 max", func(t *testing.T) {
		var v uint8 = 255
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 255 {
			t.Errorf("Expected 255, got %d", v)
		}
	})

	t.Run("uint16 max", func(t *testing.T) {
		var v uint16 = 65535
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 65535 {
			t.Errorf("Expected 65535, got %d", v)
		}
	})

	t.Run("uint32 max", func(t *testing.T) {
		var v uint32 = ^uint32(0)
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != ^uint32(0) {
			t.Errorf("Expected max uint32, got %d", v)
		}
	})

	t.Run("uint64 max", func(t *testing.T) {
		var v uint64 = 18446744073709551615 // math.MaxUint64
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 18446744073709551615 {
			t.Errorf("Expected max uint64, got %d", v)
		}
	})
}

func TestUnmarshal_Uint_TruncatedData(t *testing.T) {
	t.Run("uint truncated", func(t *testing.T) {
		data := make([]byte, 4) // needs 8 bytes for uint
		err := Unmarshal(data, (*uint)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("uint8 truncated", func(t *testing.T) {
		data := make([]byte, 0) // needs 1 byte for uint8
		err := Unmarshal(data, (*uint8)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("uint16 truncated", func(t *testing.T) {
		data := make([]byte, 1) // needs 2 bytes for uint16
		err := Unmarshal(data, (*uint16)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("uint32 truncated", func(t *testing.T) {
		data := make([]byte, 3) // needs 4 bytes for uint32
		err := Unmarshal(data, (*uint32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("uint64 truncated", func(t *testing.T) {
		data := make([]byte, 7) // needs 8 bytes for uint64
		err := Unmarshal(data, (*uint64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})
}

func TestUnmarshal_Uint_ExtraBytes(t *testing.T) {
	t.Run("uint extra bytes", func(t *testing.T) {
		data := make([]byte, 9) // needs exactly 8 bytes for uint
		err := Unmarshal(data, (*uint)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("uint8 extra bytes", func(t *testing.T) {
		data := make([]byte, 2) // needs exactly 1 byte for uint8
		err := Unmarshal(data, (*uint8)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("uint16 extra bytes", func(t *testing.T) {
		data := make([]byte, 3) // needs exactly 2 bytes for uint16
		err := Unmarshal(data, (*uint16)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("uint32 extra bytes", func(t *testing.T) {
		data := make([]byte, 5) // needs exactly 4 bytes for uint32
		err := Unmarshal(data, (*uint32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})

	t.Run("uint64 extra bytes", func(t *testing.T) {
		data := make([]byte, 9) // needs exactly 8 bytes for uint64
		err := Unmarshal(data, (*uint64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})
}

func TestMarshal_Unmarshal_Uint_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"uint roundtrip", uint(100), 8},
		{"uint8 roundtrip", uint8(255), 1},
		{"uint16 roundtrip", uint16(300), 2},
		{"uint32 roundtrip", uint32(50000), 4},
		{"uint64 roundtrip", uint64(9223372036854775800), 8},
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
			case uint:
				dest = new(uint)
			case uint8:
				dest = new(uint8)
			case uint16:
				dest = new(uint16)
			case uint32:
				dest = new(uint32)
			case uint64:
				dest = new(uint64)
			}

			err = Unmarshal(data, dest)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			// Verify roundtrip by checking the unmarshaled value matches original
			switch v := tt.input.(type) {
			case uint:
				if *dest.(*uint) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint), v)
				}
			case uint8:
				if *dest.(*uint8) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint8), v)
				}
			case uint16:
				if *dest.(*uint16) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint16), v)
				}
			case uint32:
				if *dest.(*uint32) != v {
					t.Errorf("Roundtrip value = %d; want %d", *dest.(*uint32), v)
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
// Helper Constants for Test Values (Unsigned)
// ============================================================================

const (
	maxUint8  = uint8(255)
	maxUint16 = uint16(65535)
	maxUint32 = uint32(^uint32(0))
	maxUint64 = uint64(math.MaxUint64)
	minUint   = uint(0)
	minUint8  = uint8(0)
	minUint16 = uint16(0)
	minUint32 = uint32(0)
	minUint64 = uint64(0)
)

// ============================================================================
// Size Validation Tests (Separate from Type Validation)
// ============================================================================

func TestUnmarshal_Uint_SizeMismatch(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		dataLen int
		wantErr bool
	}{
		{"uint", (*uint)(nil), 4, true},      // 8 bytes needed - truncated data
		{"uint8", (*uint8)(nil), 0, true},    // 1 byte needed - truncated data
		{"uint16", (*uint16)(nil), 3, false}, // 2 bytes needed - extra bytes (not truncated)
		{"uint32", (*uint32)(nil), 7, false}, // 4 bytes needed - extra bytes (not truncated)
		{"uint64", (*uint64)(nil), 7, true},  // 8 bytes needed - truncated data
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataLen)
			err := Unmarshal(data, tt.input)

			if err == nil && tt.wantErr {
				t.Fatalf("Unmarshal() expected error but got nil")
			}

			switch tt.name {
			case "uint", "uint8", "uint64":
				if !errors.Is(err, ErrTruncatedData) && err != nil {
					t.Errorf("Unmarshal() error = %v; want ErrTruncatedData for insufficient data", err)
				}
			case "uint16", "uint32":
				if !errors.Is(err, ErrExtraBytes) && err != nil {
					t.Errorf("Unmarshal() error = %v; want ErrExtraBytes for extra bytes", err)
				}
			}
		})
	}
}
