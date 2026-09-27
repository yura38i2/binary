package binary

import (
	"bytes"
	"testing"
)

func TestMarshalBool(t *testing.T) {
	tests := []struct {
		name     string
		input    bool
		expected []byte
		wantErr  bool
	}{
		{"marshal true", true, []byte{0x01}, false},
		{"marshal false", false, []byte{0x00}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if !bytes.Equal(data, tt.expected) {
					t.Errorf("Marshal(%v) = %v, expected %v", tt.input, data, tt.expected)
				}
			}
		})
	}
}

func TestMarshalBoolNil(t *testing.T) {
	var nilInput interface{} = (*bool)(nil)
	_, err := Marshal(nilInput)
	if err == nil {
		t.Errorf("Marshal(nil bool) should return error")
		return
	}
	if !bytes.Contains([]byte(err.Error()), []byte("cannot marshal nil")) {
		t.Errorf("Marshal() error = %v, expected nil error message", err)
	}

	var zeroInput interface{} = (*bool)(nil) // pointer to nil bool
	_, err = Marshal(zeroInput)
	if err == nil {
		t.Errorf("Marshal(nil *bool) should return error")
		return
	}
}

func TestUnmarshalBool(t *testing.T) {
	tests := []struct {
		name          string
		data          []byte
		wantTrue      bool
		wantErr       bool
		expectedError string
	}{
		{"unmarshal true from [1]", []byte{0x01}, true, false, ""},
		{"unmarshal false from [0]", []byte{0x00}, false, false, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val bool
			err := Unmarshal(tt.data, &val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && val != tt.wantTrue {
				t.Errorf("Unmarshal(%v) = %v, expected %v", tt.data, val, tt.wantTrue)
			}
		})
	}
}

func TestUnmarshalBoolTruncated(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"empty buffer", []byte{}, true},
		{"partial data", []byte{0x12}, true}, // 2 bytes is not enough for bool
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val bool
			err := Unmarshal(tt.data, &val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Error("Expected error for truncated data")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
			}
		})
	}
}

func TestUnmarshalBoolInvalidValue(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"invalid value 1 (2)", []byte{0x02}, true},
		{"invalid value 3", []byte{0x03}, true},
		{"invalid value 15", []byte{0x0F}, true},
		{"invalid value 255 (-1 as byte)", []byte{0xFF}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val bool
			err := Unmarshal(tt.data, &val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Error("Expected error for invalid boolean value")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
			}

			// Verify the value was not set to a garbage value even on error
			if tt.wantErr && val { // true means 1, which is valid but shouldn't happen for invalid inputs
				t.Log("Warning: bool may have been partially set")
			}
		})
	}
}

func TestUnmarshalBoolExtraBytes(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"extra byte after valid true", []byte{0x01, 0x02}, true},
		{"extra bytes after valid false", []byte{0x00, 0x03, 0x04}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var val bool
			err := Unmarshal(tt.data, &val)

			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err == nil {
				t.Error("Expected ErrExtraBytes for extra data")
			} else if !tt.wantErr && err != nil {
				t.Errorf("Unmarshal() unexpected error: %v", err)
			}
		})
	}
}

func TestMarshalBoolRoundTrip(t *testing.T) {
	tests := []struct {
		name string
		val  bool
	}{
		{"true roundtrip", true},
		{"false roundtrip", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := Marshal(tt.val)
			if err != nil {
				t.Fatalf("Marshal() error: %v", err)
			}

			var result bool
			err = Unmarshal(data, &result)
			if err != nil {
				t.Errorf("Unmarshal() after Marshal() error: %v", err)
				return
			}

			if result != tt.val {
				t.Errorf("Roundtrip failed: Marshal(%v) -> Unmarshal() = %v, expected %v", tt.val, result, tt.val)
			}
		})
	}
}
