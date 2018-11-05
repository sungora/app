package lg

import "gopkg.in/sungora/app.v1/lg/message"

func init() {

	message.SetMessage(map[int]string{
		1010: "smaple message log",
	})

}
