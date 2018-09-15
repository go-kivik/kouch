package target

import (
	"testing"

	"github.com/go-kivik/kouch"
)

func TestHelpText(t *testing.T) {
	for scope := kouch.TargetScope(0); scope < kouch.TargetLastScope+1; scope++ {
		result := HelpText(scope)
		if result == "" {
			t.Errorf("No help text defined for %s", kouch.TargetScopeName(scope))
		}
	}
}
