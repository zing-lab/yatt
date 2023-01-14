package utils

type ConfigKey string

const (
	CurrentRow       ConfigKey = "currentRow"
	CurrentNoteSheet ConfigKey = "currentNoteSheet"
	UncheckedOnly    ConfigKey = "unchecked_only"
	PerPage          ConfigKey = "per_page"
	Tags             ConfigKey = "tags"
	CurrentTagIdx    ConfigKey = "current_tag_index"
)

const (
	EmptyString = ""
)
