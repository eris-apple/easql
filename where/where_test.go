package where_test

import (
	"github.com/eris-apple/easql/where"
	"reflect"
	"testing"
)

type TestStruct struct {
	Name  string `map:"name"`
	Age   int    `map:"age"`
	Email string `map:"email"`
	Phone string `map:"phone"`
}

func TestNewWhereCondition(t *testing.T) {
	tests := []struct {
		name         string
		input        interface{}
		includeEmpty bool
		expected     map[string]interface{}
		expectError  bool
	}{
		{
			name: "Include all fields",
			input: TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "johndoe@example.com",
				Phone: "+123456789",
			},
			includeEmpty: true,
			expected: map[string]interface{}{
				"name":  "John Doe",
				"age":   30,
				"email": "johndoe@example.com",
				"phone": "+123456789",
			},
			expectError: false,
		},
		{
			name: "Exclude empty fields",
			input: TestStruct{
				Name:  "John Doe",
				Age:   30,
				Email: "",
				Phone: "",
			},
			includeEmpty: false,
			expected: map[string]interface{}{
				"name": "John Doe",
				"age":  30,
			},
			expectError: false,
		},
		{
			name: "Empty struct with includeEmpty=true",
			input: TestStruct{
				Name:  "",
				Age:   0,
				Email: "",
				Phone: "",
			},
			includeEmpty: true,
			expected: map[string]interface{}{
				"name":  "",
				"age":   0,
				"email": "",
				"phone": "",
			},
			expectError: false,
		},
		{
			name: "Empty struct with includeEmpty=false",
			input: TestStruct{
				Name:  "",
				Age:   0,
				Email: "",
				Phone: "",
			},
			includeEmpty: false,
			expected:     map[string]interface{}{},
			expectError:  false,
		},
		{
			name:         "Invalid input (not a struct)",
			input:        "not a struct",
			includeEmpty: true,
			expected:     nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := where.NewWhereCondition(tt.input, tt.includeEmpty)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("expected: %#v, got: %#v", tt.expected, result)
			}
		})
	}
}

func TestIsEmptyValue(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected bool
	}{
		{"Empty string", "", true},
		{"Non-empty string", "hello", false},
		{"Zero int", 0, true},
		{"Non-zero int", 42, false},
		{"Zero float", 0.0, true},
		{"Non-zero float", 3.14, false},
		{"Nil pointer", (*int)(nil), true},
		{"Non-nil pointer", new(int), false},
		{"Empty slice", []int{}, true},
		{"Non-empty slice", []int{1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := reflect.ValueOf(tt.input)
			result := where.IsEmptyValue(v)
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestWhere_String(t *testing.T) {
	tests := []struct {
		name     string
		input    where.Where
		expected string
	}{
		{"Empty map", where.Where{}, ""},
		{"Empty value", where.Where{"key": ""}, "key="},
		{"Empty key", where.Where{"": "value"}, "=value"},
		{"Correct map", where.Where{"key": "value"}, "key=value"},
		{"Multiple", where.Where{"key": "value", "key2": "value2"}, "key=value key2=value2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.input.String()
			if result != tt.expected {
				t.Errorf("expected: %v, got: %v", tt.expected, result)
			}
		})
	}
}
