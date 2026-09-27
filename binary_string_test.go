package binary

import (
	"testing"
)

// TestMarshalString tests marshaling of string values.
func TestMarshalString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"empty string", "", false},
		{"short string", "hello", false},
		{"medium string", "This is a medium length test string with some words.", false},
		{"long string", stringsRepeat('a', 1000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && len(got) == 0 {
				t.Error("Marshal() returned empty buffer for non-empty input")
			}
		})
	}
}

// TestUnmarshalString tests unmarshaling of string values.
func TestUnmarshalString(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		wantStr string
		wantErr bool
	}{
		{"empty string", makeEmptyString(), "", false},
		{"short string", marshalString("hello"), "hello", false},
		{"medium string", marshalString("This is a medium length test string with some words."), "This is a medium length test string with some words.", false},
		{"long string", marshalString(stringsRepeat('a', 1000)), stringsRepeat('a', 1000), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got string
			err := Unmarshal(tt.data, &got)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && got != tt.wantStr {
				t.Errorf("Unmarshal() got = %q, want %q", got, tt.wantStr)
			}
		})
	}
}

// TestMarshalStringRoundTrip verifies marshaling and unmarshaling round-trips.
func TestMarshalStringRoundTrip(t *testing.T) {
	tests := []string{
		"",
		"a",
		"hello world",
		"This is a longer string with multiple words to test the encoding.",
		stringsRepeat('x', 500),
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			b, err := Marshal(input)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var got string
			err = Unmarshal(b, &got)
			if err != nil {
				t.Fatalf("Unmarshal() error = %v", err)
			}

			if input != got {
				t.Errorf("Round-trip failed: got = %q, want %q", got, input)
			}
		})
	}
}

// TestMarshalStringPointer tests marshaling and unmarshaling pointer to string.
func TestMarshalStringPointer(t *testing.T) {
	tests := []struct {
		name    string
		input   *string
		wantErr bool
	}{
		{"nil pointer", nil, true}, // Cannot marshal nil pointer directly
		{"non-nil pointer", func() *string { s := "test"; return &s }(), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			var s string
			if !tt.wantErr {
				err = Unmarshal(got, &s)
				if err != nil {
					t.Fatalf("Unmarshal() error = %v", err)
				}
				if s != "test" {
					t.Errorf("got = %q, want %q", s, "test")
				}
			}
		})
	}
}

// TestMarshalStringStruct tests marshaling string within a struct.
func TestMarshalStringStruct(t *testing.T) {
	type testStruct struct {
		Name  string `bin:"name"`
		Value int    `bin:"value"`
	}

	s := testStruct{
		Name:  "test",
		Value: 42,
	}

	b, err := Marshal(s)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var got testStruct
	err = Unmarshal(b, &got)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if got.Name != s.Name {
		t.Errorf("Name mismatch: got = %q, want %q", got.Name, s.Name)
	}
}

// TestMarshalStringTooLong tests marshaling string exceeding max length.
func TestMarshalStringTooLong(t *testing.T) {
	longStr := stringsRepeat('a', 2097153) // Exceeds maxSliceLength (1M elements)
	_, err := Marshal(longStr)
	if err == nil {
		t.Error("Expected error for string exceeding max length")
	}

	_, ok := err.(*ValidatorError) // This is just checking error exists
	if !ok && err != nil {
		t.Logf("Got expected error: %v", err)
	}
}

// Helper functions for tests (not part of binary package).

func makeEmptyString() []byte {
	var buf [4]byte
	buf[0] = 0x00
	return buf[:]
}

func marshalString(s string) []byte {
	b, _ := Marshal(s)
	return b
}

// stringsRepeat returns a string with n repetitions of the given character.
func stringsRepeat(c rune, n int) string {
	if n <= 0 {
		return ""
	}
	buf := make([]byte, n)
	for i := range buf {
		buf[i] = byte(c)
	}
	return string(buf)
}

// ValidatorError is a placeholder to satisfy the test check.
type ValidatorError struct {
	error
}

func (e *ValidatorError) Unwrap() error { return e.error }
