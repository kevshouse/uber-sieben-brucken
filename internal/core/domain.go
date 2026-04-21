package core

import "time"

type Version struct {
	ID        string
	Content   string
	Timestamp time.Time
}

type Snippet struct {
	ID        string
	Title     string
	OwnerID   string
	CreatedAt time.Time
}

type Citation struct {
	ID        string
	SourceID  string
	TargetID  string
	Context   string
	Timestamp time.Time
}