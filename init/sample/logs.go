package sample

import "github.com/sungora/app/lg"

func logs() {
	lg.SetMessages(map[int]string{
		1000: "Message format Fmt from 1000: %s",
		1001: "Message format Fmt from 1001: %s",
		1002: "Message format Fmt from 1002: %s",
		1003: "Message format Fmt from 1003: %s",
		1004: "Message format Fmt from 1004: %s",
		1005: "Message format Fmt from 1005: %s",
	})
}
