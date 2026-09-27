package binary

import (
	"errors"
	"math"
	"testing"
)

// ============================================================================
// Marshal Tests - Float Types Only
// ============================================================================

func TestMarshal_Float32_PositiveValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int // expected output length in bytes
	}{
		{"float32 positive", float32(100.5), 4},
		{"float32 small", float32(0.5), 4},
		{"float32 large", float32(1e6), 4},
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

func TestMarshal_Float64_PositiveValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int // expected output length in bytes
	}{
		{"float64 positive", float64(100.5), 8},
		{"float64 small", float64(0.5), 8},
		{"float64 large", float64(1e6), 8},
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

func TestMarshal_Float_ZeroValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"float32 zero", float32(0), 4},
		{"float64 zero", float64(0), 8},
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

func TestMarshal_Float_NegativeValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"float32 -1", float32(-1), 4},
		{"float64 -1", float64(-1), 8},
		{"float32 -0.5", float32(-0.5), 4},
		{"float64 -0.5", float64(-0.5), 8},
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

func TestMarshal_Float_MaxMinValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"float32 max", float32(math.MaxFloat32), 4},
		{"float64 max", float64(math.MaxFloat64), 8},
		{"float32 min", float32(-math.MaxFloat32), 4},
		{"float64 min", float64(-math.MaxFloat64), 8},
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

func TestMarshal_Float_SpecialValues(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		wantLen int
	}{
		{"float32 inf", float32(math.Inf(1)), 4},
		{"float64 inf", float64(math.Inf(1)), 8},
		{"float32 -inf", float32(math.Inf(-1)), 4},
		{"float64 -inf", float64(math.Inf(-1)), 8},
		{"float32 nan", float32(math.NaN()), 4},
		{"float64 nan", float64(math.NaN()), 8},
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
// Unmarshal Tests - Float Types Only
// ============================================================================

func TestUnmarshal_Float32_PositiveValues(t *testing.T) {
	t.Run("float32 positive", func(t *testing.T) {
		var v float32 = 100.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 100.5 {
			t.Errorf("Expected 100.5, got %f", v)
		}
	})

	t.Run("float32 small", func(t *testing.T) {
		var v float32 = 0.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0.5 {
			t.Errorf("Expected 0.5, got %f", v)
		}
	})

	t.Run("float32 large", func(t *testing.T) {
		var v float32 = 1e6
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 1e6 {
			t.Errorf("Expected 1000000, got %f", v)
		}
	})
}

func TestUnmarshal_Float64_PositiveValues(t *testing.T) {
	t.Run("float64 positive", func(t *testing.T) {
		var v float64 = 100.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 100.5 {
			t.Errorf("Expected 100.5, got %f", v)
		}
	})

	t.Run("float64 small", func(t *testing.T) {
		var v float64 = 0.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0.5 {
			t.Errorf("Expected 0.5, got %f", v)
		}
	})

	t.Run("float64 large", func(t *testing.T) {
		var v float64 = 1e6
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 1e6 {
			t.Errorf("Expected 1000000, got %f", v)
		}
	})
}

func TestUnmarshal_Float32_NegativeValues(t *testing.T) {
	t.Run("float32 -1", func(t *testing.T) {
		var v float32 = -1
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -1 {
			t.Errorf("Expected -1, got %f", v)
		}
	})

	t.Run("float32 -0.5", func(t *testing.T) {
		var v float32 = -0.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -0.5 {
			t.Errorf("Expected -0.5, got %f", v)
		}
	})
}

func TestUnmarshal_Float64_NegativeValues(t *testing.T) {
	t.Run("float64 -1", func(t *testing.T) {
		var v float64 = -1
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -1 {
			t.Errorf("Expected -1, got %f", v)
		}
	})

	t.Run("float64 -0.5", func(t *testing.T) {
		var v float64 = -0.5
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -0.5 {
			t.Errorf("Expected -0.5, got %f", v)
		}
	})
}

func TestUnmarshal_Float_ZeroValues(t *testing.T) {
	t.Run("float32 zero", func(t *testing.T) {
		var v float32 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %f", v)
		}
	})

	t.Run("float64 zero", func(t *testing.T) {
		var v float64 = 0
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != 0 {
			t.Errorf("Expected 0, got %f", v)
		}
	})
}

func TestUnmarshal_Float_MaxMinValues(t *testing.T) {
	t.Run("float32 max", func(t *testing.T) {
		var v float32 = math.MaxFloat32
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != math.MaxFloat32 {
			t.Errorf("Expected MaxFloat32, got %f", v)
		}
	})

	t.Run("float64 max", func(t *testing.T) {
		var v float64 = math.MaxFloat64
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != math.MaxFloat64 {
			t.Errorf("Expected MaxFloat64, got %f", v)
		}
	})

	t.Run("float32 min", func(t *testing.T) {
		var v float32 = -math.MaxFloat32
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -math.MaxFloat32 {
			t.Errorf("Expected -MaxFloat32, got %f", v)
		}
	})

	t.Run("float64 min", func(t *testing.T) {
		var v float64 = -math.MaxFloat64
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if v != -math.MaxFloat64 {
			t.Errorf("Expected -MaxFloat64, got %f", v)
		}
	})
}

func TestUnmarshal_Float_SpecialValues(t *testing.T) {
	t.Run("float32 inf", func(t *testing.T) {
		var v float32 = float32(math.Inf(1))
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsInf(float64(v), 1) {
			t.Errorf("Expected positive Inf, got %f (isInf=%v)", v, math.IsInf(float64(v), 1))
		}
	})

	t.Run("float64 inf", func(t *testing.T) {
		var v float64 = math.Inf(1)
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsInf(v, 1) {
			t.Errorf("Expected positive Inf, got %f (isInf=%v)", v, math.IsInf(v, 1))
		}
	})

	t.Run("float32 -inf", func(t *testing.T) {
		var v float32 = float32(math.Inf(-1))
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsInf(float64(v), -1) {
			t.Errorf("Expected negative Inf, got %f (isInf=%v)", v, math.IsInf(float64(v), -1))
		}
	})

	t.Run("float64 -inf", func(t *testing.T) {
		var v float64 = math.Inf(-1)
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsInf(v, -1) {
			t.Errorf("Expected negative Inf, got %f (isInf=%v)", v, math.IsInf(v, -1))
		}
	})

	t.Run("float32 nan", func(t *testing.T) {
		var v float32 = float32(math.NaN())
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsNaN(float64(v)) {
			t.Errorf("Expected NaN, got %f (isNaN=%v)", v, math.IsNaN(float64(v)))
		}
	})

	t.Run("float64 nan", func(t *testing.T) {
		var v float64 = math.NaN()
		data, err := Marshal(v)
		if err != nil {
			t.Fatalf("Marshal failed: %v", err)
		}

		err = Unmarshal(data, &v)
		if err != nil {
			t.Fatalf("Unmarshal failed: %v", err)
		}

		if !math.IsNaN(v) {
			t.Errorf("Expected NaN, got %f (isNaN=%v)", v, math.IsNaN(v))
		}
	})
}

func TestUnmarshal_Float32_TruncatedData(t *testing.T) {
	t.Run("float32 truncated", func(t *testing.T) {
		data := make([]byte, 3) // needs 4 bytes for float32
		err := Unmarshal(data, (*float32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})

	t.Run("float64 truncated", func(t *testing.T) {
		data := make([]byte, 7) // needs 8 bytes for float64
		err := Unmarshal(data, (*float64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrTruncatedData but got nil")
		}

		if !errors.Is(err, ErrTruncatedData) {
			t.Errorf("Unmarshal() error = %v; want ErrTruncatedData", err)
		}
	})
}

func TestUnmarshal_Float32_ExtraBytes(t *testing.T) {
	t.Run("float32 extra bytes", func(t *testing.T) {
		data := make([]byte, 5) // needs exactly 4 bytes for float32
		err := Unmarshal(data, (*float32)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})
}

func TestUnmarshal_Float64_ExtraBytes(t *testing.T) {
	t.Run("float64 extra bytes", func(t *testing.T) {
		data := make([]byte, 9) // needs exactly 8 bytes for float64
		err := Unmarshal(data, (*float64)(nil))
		if err == nil {
			t.Fatalf("Unmarshal() expected ErrExtraBytes but got nil")
		}

		if !errors.Is(err, ErrExtraBytes) {
			t.Errorf("Unmarshal() error = %v; want ErrExtraBytes", err)
		}
	})
}

func TestMarshal_Unmarshal_Float32_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   float32
		wantVal float32
	}{
		{"float32 roundtrip positive", 100.5, 100.5},
		{"float32 roundtrip zero", 0, 0},
		{"float32 roundtrip negative", -100.5, -100.5},
		{"float32 roundtrip small", 0.1, 0.1},
		{"float32 roundtrip large", float32(1e6), float32(1e6)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var dest float32
			err = Unmarshal(data, &dest)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if dest != tt.wantVal {
				t.Errorf("Roundtrip value = %f; want %f", dest, tt.wantVal)
			}
		})
	}
}

func TestMarshal_Unmarshal_Float64_RoundTrip(t *testing.T) {
	tests := []struct {
		name    string
		input   float64
		wantVal float64
	}{
		{"float64 roundtrip positive", 100.5, 100.5},
		{"float64 roundtrip zero", 0, 0},
		{"float64 roundtrip negative", -100.5, -100.5},
		{"float64 roundtrip small", 0.1, 0.1},
		{"float64 roundtrip large", float64(1e6), float64(1e6)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var dest float64
			err = Unmarshal(data, &dest)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if dest != tt.wantVal {
				t.Errorf("Roundtrip value = %f; want %f", dest, tt.wantVal)
			}
		})
	}
}

// ============================================================================
// Helper Constants for Test Values
// ============================================================================

const (
	float32Size = 4
	float64Size = 8
)

// ============================================================================
// Size Validation Tests (Separate from Type Validation)
// ============================================================================

func TestUnmarshal_Float_SizeMismatch(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		dataLen int
		wantErr bool
	}{
		{"float32", (*float32)(nil), 3, true}, // 4 bytes needed - truncated data
		{"float64", (*float64)(nil), 7, true}, // 8 bytes needed - truncated data
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := make([]byte, tt.dataLen)
			err := Unmarshal(data, tt.input)

			if err == nil && tt.wantErr {
				t.Fatalf("Unmarshal() expected error but got nil")
			}

			if !errors.Is(err, ErrTruncatedData) && err != nil {
				t.Errorf("Unmarshal() error = %v; want ErrTruncatedData for insufficient data", err)
			}
		})
	}
}

// TestMarshalUnmarshalFloat tests round-trip marshaling/unmarshaling of floating-point numbers.
// func TestMarshalUnmarshalFloat(t *testing.T) {
// 	t.Run("float32 positive", func(t *testing.T) {
// 		var v float32 = 42.0
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != 42.0 {
// 			t.Errorf("Expected 42.0, got %f", v)
// 		}
// 	})

// 	t.Run("float32 negative", func(t *testing.T) {
// 		var v float32 = -17.5
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != -17.5 {
// 			t.Errorf("Expected -17.5, got %f", v)
// 		}
// 	})

// 	t.Run("float32 zero", func(t *testing.T) {
// 		var v float32 = 0.0
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != 0.0 {
// 			t.Errorf("Expected 0.0, got %f", v)
// 		}
// 	})

// 	t.Run("float64 positive", func(t *testing.T) {
// 		var v float64 = 3.141592653589793
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != 3.141592653589793 {
// 			t.Errorf("Expected 3.141592653589793, got %f", v)
// 		}
// 	})

// 	t.Run("float64 negative", func(t *testing.T) {
// 		var v float64 = -271.8281828459045
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != -271.8281828459045 {
// 			t.Errorf("Expected -271.8281828459045, got %f", v)
// 		}
// 	})

// 	t.Run("float64 zero", func(t *testing.T) {
// 		var v float64 = 0.0
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != 0.0 {
// 			t.Errorf("Expected 0.0, got %f", v)
// 		}
// 	})

// 	t.Run("float64 NaN", func(t *testing.T) {
// 		var v float64 = math.NaN()
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if !math.IsNaN(v) {
// 			t.Errorf("Expected NaN, got %f", v)
// 		}
// 	})

// 	t.Run("float64 infinity positive", func(t *testing.T) {
// 		var v float64 = math.Inf(1)
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != math.Inf(1) {
// 			t.Errorf("Expected +Inf, got %f", v)
// 		}
// 	})

// 	t.Run("float64 infinity negative", func(t *testing.T) {
// 		var v float64 = math.Inf(-1)
// 		data, err := Marshal(v)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, &v)
// 		if err != nil {
// 			t.Fatalf("Unmarshal failed: %v", err)
// 		}

// 		if v != math.Inf(-1) {
// 			t.Errorf("Expected -Inf, got %f", v)
// 		}
// 	})
// }

// // TestMarshalFloatError tests marshaling errors for float types.
// func TestMarshalFloatError(t *testing.T) {
// 	t.Run("nil value", func(t *testing.T) {
// 		_, err := Marshal(nil)
// 		if err == nil {
// 			t.Error("Expected error for nil input")
// 		}
// 	})
// }

// // TestUnmarshalFloatError tests error conditions for float types.
// func TestUnmarshalFloatError(t *testing.T) {
// 	type TestStruct struct {
// 		Value float32 `bin:"value"`
// 	}

// 	t.Run("truncated data", func(t *testing.T) {
// 		s := TestStruct{}
// 		data, err := Marshal(s)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		truncatedData := make([]byte, len(data)-1)
// 		copy(truncatedData, data)

// 		err = Unmarshal(truncatedData, &s)

// 		if errors.Is(err, ErrTruncatedData) {
// 			t.Log("Got expected truncated data error")
// 		} else {
// 			t.Errorf("Expected ErrTruncatedData (or wrapped), got %v", err)
// 		}
// 	})

// 	t.Run("extra bytes", func(t *testing.T) {
// 		type TestStruct struct {
// 			Value float32 `bin:"value"`
// 		}

// 		s := TestStruct{}
// 		data, err := Marshal(s)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		extraData := append(data, 1, 2, 3)
// 		err = Unmarshal(extraData, &s)
// 		if err == nil {
// 			t.Error("Expected error for extra bytes")
// 		} else if !containsError(err, "extra") {
// 			t.Errorf("Expected ErrExtraBytes or similar, got %v", err)
// 		}
// 	})

// 	t.Run("nil destination", func(t *testing.T) {
// 		type TestStruct struct {
// 			Value float32 `bin:"value"`
// 		}

// 		s := TestStruct{}
// 		data, err := Marshal(s)
// 		if err != nil {
// 			t.Fatalf("Marshal failed: %v", err)
// 		}

// 		err = Unmarshal(data, nil)
// 		if err == nil {
// 			t.Error("Expected error for nil destination")
// 		}
// 	})
// }

// // TestMarshalFloatTypeSizes tests that float type sizes are preserved correctly.
// func TestMarshalFloatTypeSizes(t *testing.T) {
// 	type TestCase struct {
// 		name    string
// 		v       interface{}
// 		wantLen int
// 	}

// 	tests := []TestCase{
// 		{"float32", float32(42.0), 4}, // just the float bits
// 		{"float64", float64(42.0), 8}, // just the float bits
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			data, err := Marshal(tc.v)
// 			if err != nil {
// 				t.Fatalf("Marshal failed: %v", err)
// 			}

// 			if len(data) != tc.wantLen {
// 				t.Errorf("Expected length %d, got %d for type %s", tc.wantLen, len(data), tc.name)
// 			}
// 		})
// 	}
// }
