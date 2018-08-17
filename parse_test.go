package kouch

import (
	"testing"

	"github.com/flimzy/diff"
	"github.com/flimzy/testy"
	"github.com/spf13/pflag"
)

func TestFilenameFromFlags(t *testing.T) {
	filenameFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagFilename, "", "filename")
		})
	}
	tests := []struct {
		name     string
		target   *Target
		flags    *pflag.FlagSet
		expected *Target
		err      string
	}{
		{
			name:     "no flags",
			target:   &Target{},
			flags:    filenameFlagSet(),
			expected: &Target{},
		},
		{
			name:   "no flag defined",
			target: &Target{},
			flags:  flagSet(),
			err:    "flag accessed but not defined: filename",
		},
		{
			name:   "filename already set",
			target: &Target{Filename: "foo"},
			flags: func() *pflag.FlagSet {
				fs := filenameFlagSet()
				if err := fs.Set("filename", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagFilename + " and pass separate filename",
		},
		{
			name:   "filename set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := filenameFlagSet()
				if err := fs.Set("filename", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Filename: "bar"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.FilenameFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestDocumentFromFlags(t *testing.T) {
	idFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagDocument, "", "id")
		})
	}
	tests := []struct {
		name     string
		target   *Target
		flags    *pflag.FlagSet
		expected *Target
		err      string
	}{
		{
			name:     "no flags",
			target:   &Target{},
			flags:    idFlagSet(),
			expected: &Target{},
		},
		{
			name:   "no flag defined",
			target: &Target{},
			flags:  flagSet(),
			err:    "flag accessed but not defined: id",
		},
		{
			name:   "id already set",
			target: &Target{Document: "321"},
			flags: func() *pflag.FlagSet {
				fs := idFlagSet()
				if err := fs.Set("id", "123"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagDocument + " and pass document ID as part of the target",
		},
		{
			name:   "id set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := idFlagSet()
				if err := fs.Set("id", "123"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Document: "123"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.DocumentFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}

func TestDatabaseFromFlags(t *testing.T) {
	dbFlagSet := func() *pflag.FlagSet {
		return flagSet(func(pf *pflag.FlagSet) {
			pf.String(FlagDatabase, "", "db")
		})
	}
	tests := []struct {
		name     string
		target   *Target
		flags    *pflag.FlagSet
		expected *Target
		err      string
	}{
		{
			name:     "no flags",
			target:   &Target{},
			flags:    dbFlagSet(),
			expected: &Target{},
		},
		{
			name:   "no flag defined",
			target: &Target{},
			flags:  flagSet(),
			err:    "flag accessed but not defined: database",
		},
		{
			name:   "id already set",
			target: &Target{Database: "foo"},
			flags: func() *pflag.FlagSet {
				fs := dbFlagSet()
				if err := fs.Set("database", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			err: "Must not use --" + FlagDatabase + " and pass database as part of the target",
		},
		{
			name:   "id set anew",
			target: &Target{},
			flags: func() *pflag.FlagSet {
				fs := dbFlagSet()
				if err := fs.Set("database", "bar"); err != nil {
					t.Fatal(err)
				}
				return fs
			}(),
			expected: &Target{Database: "bar"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.target.DatabaseFromFlags(test.flags)
			testy.Error(t, test.err, err)
			if d := diff.Interface(test.expected, test.target); d != nil {
				t.Error(d)
			}
		})
	}
}
