package comverter

import (
	"testing"
)

func TestCommentBoundaryMatching(t *testing.T) {
	tests := []struct {
		boundary *CommentBoundary
		input    string
		expected bool
	}{
		// Javadoc exact header tests
		{JavadocExactHeader, "/**", true},
		{JavadocExactHeader, "/***", false},
		{JavadocExactHeader, "/*", false},

		// Javadoc multiple asterisk tests
		{JavadocMultipleAsterisk, "/***", true},
		{JavadocMultipleAsterisk, "/****", true},
		{JavadocMultipleAsterisk, "/*****", true},
		{JavadocMultipleAsterisk, "/**", false}, // Should not match exact /**

		// Footer tests
		{JavadocExactFooter, "*/", true},
		{JavadocExactFooter, "**/", false},
		{JavadocMultipleFooter, "**/", true},
		{JavadocMultipleFooter, "***/", true},
		{JavadocMultipleFooter, "*/", false}, // Should not match single */

		// Single line comment tests
		{ForwardSlashTwice, "//", true},
		{ForwardSlashTwice, "///", false},
		{ForwardSlashMultiple, "///", true},
		{ForwardSlashMultiple, "////", true},
		{ForwardSlashMultiple, "//", false},

		// Doxygen tests
		{DoxygenTripleSlash, "///", true},
		{DoxygenTripleSlash, "//", false},
		{DoxygenBangSlash, "//!", true},
		{DoxygenBangSlash, "//", false},

		// Python tests
		{PythonTripleQuote, `"""`, true},
		{PythonTripleQuote, `'''`, true},
		{PythonTripleQuote, `"`, false},
		{PythonHashComment, "#", true},
		{PythonHashComment, "##", true},
		{PythonHashComment, "###", true},

		// Hash comment tests
		{HashComment, "#", true},
		{HashComment, "##", false},
		{HashMultiple, "##", true},
		{HashMultiple, "###", true},
		{HashMultiple, "#", false},
	}

	for _, test := range tests {
		result := test.boundary.Match(test.input)
		if result != test.expected {
			t.Errorf("Boundary %s with input '%s': expected %v, got %v",
				test.boundary.Name, test.input, test.expected, result)
		}
	}
}

func TestCommentBoundaryFindMatch(t *testing.T) {
	tests := []struct {
		boundary *CommentBoundary
		input    string
		expected string
	}{
		{JavadocExactHeader, "/**", "/**"},
		{JavadocMultipleAsterisk, "/****", "/****"},
		{ForwardSlashMultiple, "////", "////"},
		{PythonTripleQuote, `"""`, `"""`},
		{JavadocExactHeader, "not a match", ""},
	}

	for _, test := range tests {
		result := test.boundary.FindMatch(test.input)
		if result != test.expected {
			t.Errorf("Boundary %s FindMatch with input '%s': expected '%s', got '%s'",
				test.boundary.Name, test.input, test.expected, result)
		}
	}
}

func ExampleJavadocMultipleAsterisk() {
	// This boundary matches /** followed by any number of additional asterisks
	examples := []string{"/***", "/****", "/*****"}

	for _, example := range examples {
		if JavadocMultipleAsterisk.Match(example) {
			println("Matched:", example)
		}
	}
}
