[![Build Status](https://travis-ci.org/go-kivik/kouch.svg?branch=master)](https://travis-ci.org/go-kivik/kouch) [![Codecov](https://img.shields.io/codecov/c/github/go-kivik/kouch.svg?style=flat)](https://codecov.io/gh/go-kivik/kouch) [![Go Report Card](https://goreportcard.com/badge/github.com/go-kivik/kouch)](https://goreportcard.com/report/github.com/go-kivik/kouch) [![GoDoc](https://godoc.org/github.com/go-kivik/kouch?status.svg)](http://godoc.org/github.com/go-kivik/kouch)

# Kouch

Kouch is a command-line interface for CouchDB, intended to facilitate ease of
scripting or manual interaction with CouchDB.

It takes great inspiration from [curl](https://curl.haxx.se/), the command-line
tool for transferring data with URLs.

Kouch aims to make CouchDB administration and scripting easier, by providing a
simple, CouchDB-centric command-line tool for performing routine administration
and debugging operations, without the cumbersome task of manually constructing
HTTP requests for use with curl.

Kouch can also output (and read input) to pretty-printed JSON or YAML, rather
than CouchDB's native JSON format, for more more human-friendly interaction with
documents.

# Example Usage

## Fetch a document

    $ kouch get doc localhost:5984/foo/bar
    {"_attachments":{"foo.txt":{"content_type":"text/plain","digest":"md5-WiGw80mG3uQuqTKfUnIZsg==","length":9,"revpos":3,"stub":true}},"_id":"bar","_rev":"3-13438fbeeac7271383a42b57511f03ea","a":"c"}

## Fetch a document, pretty JSON output

    $ kouch get doc localhost:5984/foo/bar -F json --json-indent " "
    {
     "_attachments": {
      "foo.txt": {
       "content_type": "text/plain",
       "digest": "md5-WiGw80mG3uQuqTKfUnIZsg==",
       "length": 9,
       "revpos": 3,
       "stub": true
      }
     },
     "_id": "bar",
     "_rev": "3-13438fbeeac7271383a42b57511f03ea",
     "a": "c"
    }

## Fetch a document, YAML output

    $ kouch get doc localhost:5984/foo/bar --output-format yaml
    _attachments:
      foo.txt:
        content_type: text/plain
        digest: md5-WiGw80mG3uQuqTKfUnIZsg==
        length: 9
        revpos: 3
        stub: true
    _id: bar
    _rev: 3-13438fbeeac7271383a42b57511f03ea
    a: c

## Fetch a document, showing only the headers

    $ kouch get doc localhost:5984/foo/bar -I
    Cache-Control: must-revalidate
    Content-Length: 198
    Content-Type: application/json
    Date: Sun, 26 Aug 2018 16:32:12 GMT
    Etag: "3-13438fbeeac7271383a42b57511f03ea"
    Server: CouchDB/2.1.1 (Erlang OTP/17)
    X-Couch-Request-Id: 0ff82e5498
    X-Couchdb-Body-Time: 0

# Current status

Kouch is still in the early stages of development. Most features have not yet
been implemented. But fast progress is being made, and your contributions are
also welcome!

# License

This software is released under the terms of the Apache 2.0 license. See
LICENCE.md, or read the [full license](http://www.apache.org/licenses/LICENSE-2.0).
