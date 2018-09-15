package kouch

// TargetHelpText returns the help text to describe the valid target format(s)
// for the specified scope, or "" if the scope isn't defined.
func TargetHelpText(scope TargetScope) string {
	return helpText[scope]
}

var helpText = map[TargetScope]string{
	TargetRoot: `[target] should be the URL to the root of the CouchDB server. Examples:

  - http://localhost:5984/
  - example.com:5000
  - foo.com/couchdb/
`,
	TargetDatabase: `[target] may be a full or relative URL to the database. Examples:

  - foo                          -- Database 'foo', relative to the Root URL
  - http://localhost:5984/_users -- The '_users' database on localhost
  - example.com:5000/root/foo    -- The 'foo' database on example.com, with CouchDB served at the 'root/' path.

Any slashes in the database name, must be URL-encoded.
`,
	TargetDocument: `[target] may be a full or relative URL to the document. Examples:

  - bar                           -- Document 'bar' in the default database and Root URL
  - foo/bar                       -- Document 'bar' in the database 'foo' at the default Root URL
  - _design/bar                   -- Relative URL to a design document in the current database
  - _local/bar                    -- Relative URL to a non-replicating document in the current database
  - foo/_design/bar               -- The 'bar' design doc in the 'foo' database at the current Root URL
  - http://localhost:5984/foo/bar -- Full URL

Except for _design/ and _local/ documents, any slashes in a database name or document ID must be URL-encoded.
`,
	TargetAttachment: `[target] may be a full or relative URL to the attachment. Examples:

  - baz.txt                          -- Attachment 'baz.txt' from the current document
  - bar/baz.jpg                      -- Attachment 'baz.jpg' from the 'bar' document in the current database
  - foo/bar/baz.png                  -- Attachment 'baz.png' from the 'bar' doc in the 'foo' database at the current Root URL
  - foo/_design/bar/baz.html         -- Attachment 'baz.html' from the 'bar' design doc in the 'foo' databasae at the current Root URL
  - http://host.com/foo/bar/baz.html -- Full URL

  Except for _design/ and _local/ documents, any slashes in a database name, document id, or filename must be URL-encoded.
`,
}
