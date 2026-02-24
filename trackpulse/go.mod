module trackpulse

go 1.19

require (
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/tarm/serial v0.0.0-20180830185346-98f6abe2eb07
	golang.org/x/crypto v0.48.0
)

require golang.org/x/sys v0.41.0 // indirect

replace golang.org/x/sys v0.41.0 => golang.org/x/sys v0.15.0
