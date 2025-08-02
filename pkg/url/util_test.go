package url

import "testing"

func TestNormalize(t *testing.T) {
	testCases := []struct {
		name         string
		input        string
		expected     string
		errorPresent bool
	}{
		{
			name:         "F1: test case 1",
			input:        "http://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 2",
			input:        "http://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 3",
			input:        "https://www.hello.com/world",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 4",
			input:        "https://www.hello.com/world/",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 5",
			input:        "https://www.hello.com/world?unit=testing",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 6",
			input:        "https://www.hello.com/world?unit=testing#foo",
			expected:     "www.hello.com/world",
			errorPresent: false,
		},
		{
			name:         "F1: test case 7",
			input:        "com.invalid .https://",
			expected:     "",
			errorPresent: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := Normalize(testCase.input)
			if (err != nil) != testCase.errorPresent {
				t.Errorf("%s failed, unexpected error: %v", testCase.name, err)
			}
			if result != testCase.expected {
				t.Errorf("%s failed, %s != %s", testCase.name, result, testCase.expected)
			}
		})
	}
}
