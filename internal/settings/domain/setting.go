package domain

// Setting is a tenant-scoped, DB-backed key/value pair used for
// parametrizing layout, palette and other system-level configuration —
// e.g. key "theme" with value {"palette": "navy-coral"}.
type Setting struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}
