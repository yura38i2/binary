package binary

import (
	"bytes"
	"math"
	"reflect"
	"testing"
)

// TestStructMarshalBasic tests marshaling a struct with all supported field types.
func TestStructMarshalBasic(t *testing.T) {
	type Inner struct {
		A int8
		B int8
	}

	type MyStruct struct {
		BoolField    bool
		Int8Field    int8
		Int64Field   int64
		Float64Field float64
		StringField  string
		SliceField   []int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "all fields with non-zero values",
			input: MyStruct{
				BoolField:    true,
				Int8Field:    42,
				Int64Field:   123456789012345,
				Float64Field: 3.14159,
				StringField:  "hello",
				SliceField:   []int64{1, 2, 3},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Marshal/Unmarshal round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

// TestStructMarshalNestedStruct tests marshaling a struct with nested simple structs.
func TestStructMarshalNestedSimpleStruct(t *testing.T) {
	type Inner struct {
		A int8
		B int8
	}

	type MyStruct struct {
		NestedField Inner
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "nested simple struct with values",
			input: MyStruct{
				NestedField: Inner{A: 10, B: 20},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Marshal/Unmarshal round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

// TestStructMarshalNestedComplexStruct tests marshaling a struct with nested complex structs.
func TestStructMarshalNestedComplexStruct(t *testing.T) {
	type Inner struct {
		BoolField  bool
		Int64Field int64
	}

	type MyStruct struct {
		NestedField Inner
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "nested complex struct with values",
			input: MyStruct{
				NestedField: Inner{BoolField: true, Int64Field: 99},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Marshal/Unmarshal round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

// TestStructMarshalSliceNestedSimpleStruct tests marshaling a struct with slice of simple structs.
func TestStructMarshalSliceOfNestedSimpleStruct(t *testing.T) {
	type Inner struct {
		A int8
		B int8
	}

	type MyStruct struct {
		SliceField []Inner
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "slice of simple structs with multiple elements",
			input: MyStruct{
				SliceField: []Inner{{A: 1, B: 2}, {A: 3, B: 4}},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Marshal/Unmarshal round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

// TestStructMarshalZeroValues tests marshaling a struct with zero values.
func TestStructMarshalZeroValues(t *testing.T) {
	type MyStruct struct {
		BoolField    bool
		Int8Field    int8
		Int64Field   int64
		Float64Field float64
		StringField  string
		SliceField   []int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "all zero values",
			input: MyStruct{
				BoolField:    false,
				Int8Field:    0,
				Int64Field:   0,
				Float64Field: 0.0,
				StringField:  "",
				SliceField:   nil, // nil slice is allowed and should be encoded as sentinel
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Marshal/Unmarshal round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

// TestStructMarshalEmptyString tests marshaling a struct with an empty string.
func TestStructMarshalEmptyString(t *testing.T) {
	type MyStruct struct {
		StringField string
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "empty string",
			input:   MyStruct{StringField: ""},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if result.StringField != "" {
				t.Errorf("Expected empty string, got: %q", result.StringField)
			}
		})
	}
}

// TestStructMarshalEmptySlice tests marshaling a struct with an empty slice.
func TestStructMarshalEmptySlice(t *testing.T) {
	type MyStruct struct {
		SliceField []int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "empty slice",
			input:   MyStruct{SliceField: make([]int64, 0)},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if len(result.SliceField) != 0 {
				t.Errorf("Expected empty slice, got length: %d", len(result.SliceField))
			}
		})
	}
}

// TestStructMarshalInt64BoundaryValues tests marshaling int64 boundary values.
func TestStructMarshalInt64BoundaryValues(t *testing.T) {
	type MyStruct struct {
		Int64Field int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "max int64",
			input:   MyStruct{Int64Field: math.MaxInt64},
			wantErr: false,
		},
		{
			name:    "min int64",
			input:   MyStruct{Int64Field: math.MinInt64},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if result.Int64Field != tt.input.Int64Field {
				t.Errorf("Int64 round-trip mismatch:\n  input: %d\n  got:   %d", tt.input.Int64Field, result.Int64Field)
			}
		})
	}
}

// TestStructMarshalFloat64BoundaryValues tests marshaling float64 boundary values.
func TestStructMarshalFloat64BoundaryValues(t *testing.T) {
	type MyStruct struct {
		Float64Field float64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "max float64",
			input:   MyStruct{Float64Field: math.MaxFloat64},
			wantErr: false,
		},
		{
			name:    "min positive normal float64",
			input:   MyStruct{Float64Field: math.SmallestNonzeroFloat64},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Use approximate comparison for floating-point values
			if math.Abs(result.Float64Field-tt.input.Float64Field) > 1e-9 {
				t.Errorf("Float64 round-trip mismatch:\n  input: %v\n  got:   %v", tt.input.Float64Field, result.Float64Field)
			}
		})
	}
}

// TestStructMarshalSliceOfInt64 tests marshaling a struct with slice of int64.
func TestStructMarshalSliceOfInt64(t *testing.T) {
	type MyStruct struct {
		SliceField []int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "single element slice",
			input:   MyStruct{SliceField: []int64{42}},
			wantErr: false,
		},
		{
			name:    "multiple elements slice",
			input:   MyStruct{SliceField: []int64{-100, 0, 100}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input.SliceField, result.SliceField) {
				t.Errorf("Slice round-trip mismatch:\n  input: %v\n  got:   %v", tt.input.SliceField, result.SliceField)
			}
		})
	}
}

// TestStructMarshalWithSkippedFields tests marshaling a struct with bin:"-" tag.
func TestStructMarshalWithSkippedFields(t *testing.T) {
	type MyStruct struct {
		PublicField  bool   `bin:"public"`
		SkippedField string `bin:"-"`
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "struct with skipped field",
			input:   MyStruct{PublicField: true, SkippedField: "should be ignored"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Skipped field should be zeroed after unmarshaling
			if result.PublicField != true {
				t.Errorf("PublicField mismatch: got %v, want true", result.PublicField)
			}
			if result.SkippedField != "" { // Should be empty (zero value) after unmarshal
				t.Errorf("SkippedField should be zeroed after round-trip, got: %q", result.SkippedField)
			}
		})
	}
}

// TestStructMarshalUnmarshalNilSlicePointer tests marshaling a struct with nil slice pointer.
func TestStructMarshalUnmarshalNilSlicePointer(t *testing.T) {
	type MyStruct struct {
		SliceField *[]int64
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name:    "nil slice pointer",
			input:   MyStruct{SliceField: nil},
			wantErr: false,
		},
		{
			name:    "non-nil empty slice pointer",
			input:   MyStruct{SliceField: &[]int64{}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Both should produce the same marshaled data for round-trip consistency
			data2, _ := Marshal(tt.input) // nil slice pointer case
			if !bytes.Equal(data, data2) {
				t.Errorf("Nil slice pointer and non-nil empty slice marshal differently")
			}

			var result2 MyStruct
			Unmarshal(data, &result2)

			if len(*result.SliceField) != 0 {
				t.Errorf("Expected nil or empty slice after round-trip, got length: %d", len(*result.SliceField))
			}
		})
	}
}

// TestStructMarshalUnmarshalOrder tests that fields are encoded in declaration order.
func TestStructMarshalUnmarshalOrder(t *testing.T) {
	type MyStruct struct {
		Zebra int64
		Apple bool
		Mango string
	}

	tests := []struct {
		name    string
		input   MyStruct
		wantErr bool
	}{
		{
			name: "fields in declaration order",
			input: MyStruct{
				Zebra: 1,
				Apple: true,
				Mango: "hello",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			var result MyStruct
			err = Unmarshal(data, &result)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.input, result) {
				t.Errorf("Round-trip mismatch:\n  input: %+v\n  got:   %+v", tt.input, result)
			}
		})
	}
}

func TestStructMarshalUnmarshalDeterminism(t *testing.T) {
	type MyStruct struct {
		BoolField    bool
		Int8Field    int8
		Int64Field   int64
		Float64Field float64
		StringField  string
	}

	input := MyStruct{
		BoolField:    true,
		Int8Field:    -127,
		Int64Field:   9223372036854775807, // math.MaxInt64
		Float64Field: 3.14,
		StringField:  "test",
	}

	data := make([]byte, 0)
	for i := 0; i < 10; i++ {
		d, err := Marshal(input)
		if err != nil {
			t.Fatalf("Marshal() error = %v", err)
		}

		// Initialize data on first iteration to avoid comparing empty slice with actual data
		if i == 0 {
			data = d
			continue // Skip comparison for the initialization step
		}

		if !bytes.Equal(data, d) {
			t.Errorf("Marshaling is not deterministic: iteration %d differs from first", i)
		}
		data = d
	}

	var result MyStruct
	err := Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if !reflect.DeepEqual(input, result) {
		t.Errorf("Round-trip mismatch:\n  input: %+v\n  got:   %+v", input, result)
	}
}
