package scan

import . "time"

type File struct {
	Path string
	FQDN string
	Sum  string
	Time Time
}
