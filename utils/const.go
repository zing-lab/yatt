package utils

type ConfigKey string

const (
	CurrentRow       ConfigKey = "currentRow"
	CurrentNoteSheet ConfigKey = "currentNoteSheet"
	MarkedOnly       ConfigKey = "marked_only"
	PerPage          ConfigKey = "per_page"
)
