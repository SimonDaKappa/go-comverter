package comverter

import (
	"regexp"
)

//--------------------------------------------------------------------------------
// Boundary Families
//
// A Boundary Family is a collection of related Comment Boundaries.
// For example, the "Javadoc" family includes boundaries for
// /**, /***, */, **/, etc.
//
// When creating a Family, it is recommended to order boundaries from
// most specific to least specific. This helps in matching and prioritization.
//--------------------------------------------------------------------------------

type BoundaryFamilyName string

const ()

type BoundaryFamilyRegistry struct {
	Families map[BoundaryFamilyName][]*CommentBoundary
}

func NewBoundaryFamilyRegistry() *BoundaryFamilyRegistry {
	reg := &BoundaryFamilyRegistry{
		Families: make(map[BoundaryFamilyName][]*CommentBoundary),
	}
	reg.loadDefaultFamilies()
	return reg
}

func (r *BoundaryFamilyRegistry) loadDefaultFamilies() {
	r.Families[JavadocBoundaryFamilyName] = JavadocBoundaryFamily
	r.Families[CBlockBoundaryFamilyName] = CBlockBoundaryFamily
	r.Families[SingleBoundaryFamilyName] = SingleLineBoundaryFamily
	r.Families[DoxygenBoundaryFamilyName] = DoxygenBoundaryFamily
	r.Families[PythonBoundaryFamilyName] = PythonBoundaryFamily
	r.Families[HashBoundaryFamilyName] = HashBoundaryFamily
}

func (r *BoundaryFamilyRegistry) Family(name BoundaryFamilyName) ([]*CommentBoundary, bool) {
	boundaries, exists := r.Families[name]
	return boundaries, exists
}

func (r *BoundaryFamilyRegistry) All() map[BoundaryFamilyName][]*CommentBoundary {
	return r.Families
}

func (r *BoundaryFamilyRegistry) Register(name BoundaryFamilyName, boundaries []*CommentBoundary) {
	r.Families[name] = boundaries
}

func (r *BoundaryFamilyRegistry) Unregister(name BoundaryFamilyName) {
	delete(r.Families, name)
}

func (r *BoundaryFamilyRegistry) GetMatchingFamily(text string) (BoundaryFamilyName, []*CommentBoundary) {
	for name, boundaries := range r.Families {
		if MatchesAny(text, boundaries...) {
			return name, boundaries
		}
	}
	return "", nil
}

func (r *BoundaryFamilyRegistry) GetAllMatchingFamilies(text string) map[BoundaryFamilyName][]*CommentBoundary {
	matches := make(map[BoundaryFamilyName][]*CommentBoundary)
	for name, boundaries := range r.Families {
		if MatchesAny(text, boundaries...) {
			matches[name] = boundaries
		}
	}
	return matches
}

func (r *BoundaryFamilyRegistry) MatchesFamily(text string, family BoundaryFamilyName) bool {
	boundaries, exists := r.Families[family]
	if !exists {
		return false
	}
	return MatchesAny(text, boundaries...)
}

func (r *BoundaryFamilyRegistry) FindFirstMatchingBoundary(text string, family BoundaryFamilyName) *CommentBoundary {
	boundaries, exists := r.Families[family]
	if !exists {
		return nil
	}
	return FindFirstMatch(text, boundaries...)
}

//--------------------------------------------------------------------------------
// Comment Boundary
//
// A Comment Boundary defines a regex-based matcher for comment delimiters
// such as /*, */, //, #, etc.
//
// A boundary is consider matched IFF the entire string matches the pattern.
// Partial matches do not count.
//--------------------------------------------------------------------------------

type CommentBoundaryName string

var (
	BoundaryForwardSlashTwice    CommentBoundaryName = "ForwardSlashTwice"
	BoundaryForwardSlashMultiple CommentBoundaryName = "ForwardSlashMultiple"

	BoundaryDoxygenQtStyle     CommentBoundaryName = "DoxygenQtStyle"
	BoundaryDoxygenBangStyle   CommentBoundaryName = "DoxygenBangStyle"
	BoundaryDoxygenTripleSlash CommentBoundaryName = "DoxygenTripleSlash"
	BoundaryDoxygenBangSlash   CommentBoundaryName = "DoxygenBangSlash"

	BoundaryPythonTripleQuote CommentBoundaryName = "PythonTripleQuote"
	BoundaryPythonHashComment CommentBoundaryName = "PythonHashComment"

	BoundaryHashComment  CommentBoundaryName = "HashComment"
	BoundaryHashMultiple CommentBoundaryName = "HashMultiple"
)

// CommentBoundary defines a regex-based boundary matcher for comment delimiters
type CommentBoundary struct {
	Name    CommentBoundaryName // Human-readable name for this boundary type
	Pattern *regexp.Regexp      // Compiled regex pattern for matching
	Raw     string              // Raw regex pattern string for reference
}

// Match checks if the given string matches this boundary pattern
func (cb *CommentBoundary) Match(text string) bool {
	return cb.Pattern.MatchString(text)
}

// FindMatch returns the matching portion of the string, or empty string if no match
func (cb *CommentBoundary) FindMatch(text string) string {
	return cb.Pattern.FindString(text)
}

// NewCommentBoundary creates a new CommentBoundary with a compiled regex pattern
func NewCommentBoundary(name CommentBoundaryName, pattern string) *CommentBoundary {
	compiled := regexp.MustCompile(pattern)
	return &CommentBoundary{
		Name:    name,
		Pattern: compiled,
		Raw:     pattern,
	}
}

// MatchesAny checks if any of the provided boundaries match the given text
func MatchesAny(text string, boundaries ...*CommentBoundary) bool {
	for _, boundary := range boundaries {
		if boundary.Match(text) {
			return true
		}
	}
	return false
}

// FindFirstMatch returns the first boundary that matches the text, or nil if none match
func FindFirstMatch(text string, boundaries ...*CommentBoundary) *CommentBoundary {
	for _, boundary := range boundaries {
		if boundary.Match(text) {
			return boundary
		}
	}
	return nil
}

func FindAllMatches(text string, boundaries ...*CommentBoundary) []CommentBoundary {
	var matches []CommentBoundary
	for _, boundary := range boundaries {
		if boundary.Match(text) {
			matches = append(matches, *boundary)
		}
	}
	return matches
}

// IsSingleLineComment checks if text matches any single-line comment pattern
func IsSingleLineComment(text string) bool {
	return MatchesAny(text, ForwardSlashTwice, ForwardSlashMultiple, DoxygenTripleSlash, DoxygenBangSlash, HashComment, HashMultiple)
}

//--------------------------------------------------------------------------------
// Javadoc-style boundaries
//--------------------------------------------------------------------------------

var (
	BoundaryJavadocExactHeader      CommentBoundaryName = "JavadocExactHeader"
	BoundaryJavadocMultipleAsterisk CommentBoundaryName = "JavadocMultipleAsterisk"
	BoundaryJavadocExactFooter      CommentBoundaryName = "JavadocExactFooter"
	BoundaryJavadocMultipleFooter   CommentBoundaryName = "JavadocMultipleFooter"

	// Javadoc-style boundaries
	JavadocExactHeader      = NewCommentBoundary(BoundaryJavadocExactHeader, `^/\*\*$`)
	JavadocMultipleAsterisk = NewCommentBoundary(BoundaryJavadocMultipleAsterisk, `^/\*\*\*+$`)
	JavadocExactFooter      = NewCommentBoundary(BoundaryJavadocExactFooter, `^\*/$`)
	JavadocMultipleFooter   = NewCommentBoundary(BoundaryJavadocMultipleFooter, `^\*{2,}/$`)

	JavadocBoundaryFamilyName BoundaryFamilyName = "Javadoc"
	JavadocBoundaryFamily                        = []*CommentBoundary{
		JavadocExactHeader,
		JavadocMultipleAsterisk,
		JavadocExactFooter,
		JavadocMultipleFooter,
	}
)

//--------------------------------------------------------------------------------
// CBlock-style boundaries
//--------------------------------------------------------------------------------

var (
	CBlockBoundaryFamilyName BoundaryFamilyName = "CBlock"

	BoundaryCBlockCommentHeader CommentBoundaryName = "CBlockCommentHeader"
	BoundaryCBlockCommentFooter CommentBoundaryName = "CBlockCommentFooter"

	// C-style block comment boundaries
	CBlockCommentHeader = NewCommentBoundary(BoundaryCBlockCommentHeader, `^/\*$`)
	CBlockCommentFooter = NewCommentBoundary(BoundaryCBlockCommentFooter, `^\*/$`)

	CBlockBoundaryFamily = []*CommentBoundary{
		CBlockCommentHeader,
		CBlockCommentFooter,
	}
)

//--------------------------------------------------------------------------------
// Forward Slash boundaries
//--------------------------------------------------------------------------------

var (
	FamilyNameSingle BoundaryFamilyName = "SingleLine"

	// Single-line comment boundaries
	ForwardSlashTwice    = NewCommentBoundary(BoundaryForwardSlashTwice, `^//$`)
	ForwardSlashMultiple = NewCommentBoundary(BoundaryForwardSlashMultiple, `^/{3,}$`)

	SingleLineBoundaryFamily = []*CommentBoundary{
		ForwardSlashTwice,
		ForwardSlashMultiple,
	}
)

//--------------------------------------------------------------------------------
// Doxygen-style boundaries
//--------------------------------------------------------------------------------

var (
	FamilyNameDoxygen BoundaryFamilyName = "Doxygen"

	// Doxygen-style boundaries
	DoxygenQtStyle     = NewCommentBoundary(BoundaryDoxygenQtStyle, `^/!\*$`)
	DoxygenBangStyle   = NewCommentBoundary(BoundaryDoxygenBangStyle, `^/\*!$`)
	DoxygenTripleSlash = NewCommentBoundary(BoundaryDoxygenTripleSlash, `^///$`)
	DoxygenBangSlash   = NewCommentBoundary(BoundaryDoxygenBangSlash, `^//!$`)

	DoxygenBoundaryFamily = []*CommentBoundary{
		DoxygenQtStyle,
		DoxygenBangStyle,
		DoxygenTripleSlash,
		DoxygenBangSlash,
	}
)

//--------------------------------------------------------------------------------
// Python-style boundaries
//--------------------------------------------------------------------------------

var (
	FamilyNamePython BoundaryFamilyName = "Python"

	// Python-style boundaries
	PythonTripleQuote = NewCommentBoundary(BoundaryPythonTripleQuote, `^"""|'''$`)
	PythonHashComment = NewCommentBoundary(BoundaryPythonHashComment, `^#+$`)

	PythonBoundaryFamily = []*CommentBoundary{
		PythonTripleQuote,
		PythonHashComment,
	}
)

//--------------------------------------------------------------------------------
// Hash-style boundaries
//--------------------------------------------------------------------------------

var (
	FamilyNameHash BoundaryFamilyName = "Hash"

	// Hash-style boundaries (shell, ruby, perl, etc.)
	HashComment  = NewCommentBoundary(BoundaryHashComment, `^#$`)
	HashMultiple = NewCommentBoundary(BoundaryHashMultiple, `^#{2,}$`)

	HashBoundaryFamily = []*CommentBoundary{
		HashComment,
		HashMultiple,
	}
)
