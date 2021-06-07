package promlib

var (
	// globalNS describes a namespace used in prometheus metrics by default.
	globalNS = ""
)

// SetGlobalNamespace sets the default metrics namespace.
func SetGlobalNamespace(ns string) {
	globalNS = ns
}

// GetGlobalNamespace return the default metrics namespace.
func GetGlobalNamespace() string {
	// NOTE(a.petrukhin): possible by-design race.
	return globalNS
}
