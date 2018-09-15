package kouch

import "testing"

func TestHelpText(t *testing.T) {
	for scope := TargetScope(0); scope < TargetLastScope+1; scope++ {
		result := TargetHelpText(scope)
		if result == "" {
			t.Errorf("No help text defined for %s", TargetScopeName(scope))
		}
	}
}
