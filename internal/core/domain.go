package core

import "time"

type Version struct {
	ID        string
	Content   string
	Timestamp time.Time
}

type Snippet struct {
	ID      string
	Title   string
	OwnerID string
}

// Citation uses this Pattern: "Relationship as an Entity"
// It tracks why one snippet refers to another.
type Citation struct {
	ID        string
	SourceID  string    // ID of the Snippet that is citing
	TargetID  string    // ID of the Snippet being cited
	Context   string    // The "Reason" for the citation
	Timestamp time.Time
}