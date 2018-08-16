package target

import "testing"

func TestHelpText(t *testing.T) {
	for scope := Scope(0); scope < lastScope+1; scope++ {
		result := HelpText(scope)
		if result == "" {
			t.Errorf("No help text defined for %s", ScopeName(scope))
		}
	}
}

func TestScopeName(t *testing.T) {
	for scope := Scope(0); scope < lastScope+1; scope++ {
		result := ScopeName(scope)
		if result == "" {
			t.Errorf("No name defined for scope #%d", scope)
		}
	}
}
