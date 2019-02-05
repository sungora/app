// Контроллер главной страницы
package page

import (
	"github.com/sungora/app/cmd/app/model/mdlusers"
	"github.com/sungora/app/core"
)

func NewControlSample() core.Controller {
	return new(ControlSample)
}

type ControlSample struct {
	core.ControllerHtml
}

// GET действие по умолчанию
func (control *ControlSample) GET() (err error) {
	// сессия
	var count int
	if control.Session.Get("count") != nil {
		count, _ = control.Session.Get("count").(int)
	}
	count += 1
	control.Session.Set("count", count)

	// работа с моделью
	u := mdlusers.New(0)
	invoiceID := uint64(8697)
	u.InvoiceID = &invoiceID
	name := "Вася пупкин"
	u.Nam = &name
	u.Age = count

	control.Variables["Header"] = "Head ControlSample"
	control.Variables["User"] = u
	control.TplController += "/page/sample.html"

	return
}
