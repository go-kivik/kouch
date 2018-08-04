package kouch

// Exit statuses, borrowed from Curl. Not all Curl statuses are represented here.
const (
	// Exited with an unknown failure.
	ExitUnknownFailure = 1
	// Failed to initialize.
	ExitFailedToInitialize = 2
	// Write error. Kouch couldn't write data to a local filesystem or similar.
	ExitWriteError = 23

/*
3      URL malformed. The syntax was not correct.
5      Couldn't resolve proxy. The given proxy host could not be resolved.
6      Couldn't resolve host. The given remote host was not resolved.
7      Failed to connect to host.
8      Weird server reply. The server sent data curl couldn't parse.
18     Partial file. Only a part of the file was transferred.
22     HTTP page not retrieved. The requested url was not found or returned another error with the HTTP error code being 400 or above. This return code only appears if -f, --fail is used.
26     Read error. Various reading problems.
27     Out of memory. A memory allocation request failed.
28     Operation timeout. The specified time-out period was reached according to the conditions.
33     HTTP range error. The range "command" didn't work.
34     HTTP post error. Internal post-request generation error.
35     SSL connect error. The SSL handshaking failed.
37     FILE couldn't read file. Failed to open the file. Permissions?
43     Internal error. A function was called with a bad parameter.
45     Interface error. A specified outgoing interface could not be used.
47     Too many redirects. When following redirects, curl hit the maximum amount.
51     The peer's SSL certificate or SSH MD5 fingerprint was not OK.
52     The server didn't reply anything, which here is considered an error.
53     SSL crypto engine not found.
54     Cannot set SSL crypto engine as default.
55     Failed sending network data.
56     Failure in receiving network data.
58     Problem with the local certificate.
59     Couldn't use specified SSL cipher.
60     Peer certificate cannot be authenticated with known CA certificates.
61     Unrecognized transfer encoding.
63     Maximum file size exceeded.
65     Sending the data requires a rewind that failed.
66     Failed to initialise SSL Engine.
67     The user name, password, or similar was not accepted and curl failed to log in.
75     Character conversion failed.
76     Character conversion functions required.
77     Problem with reading the SSL CA cert (path? access rights?).
78     The resource referenced in the URL does not exist.
79     An unspecified error occurred during the SSH session.
80     Failed to shut down the SSL connection.
82     Could not load CRL file, missing or wrong format (added in 7.19.0).
83     Issuer check failed (added in 7.19.0).
85     RTSP: mismatch of CSeq numbers
86     RTSP: mismatch of Session Identifiers
89     No connection available, the session will be queued
90     SSL public key does not matched pinned public key
*/
)