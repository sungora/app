// Контроллер главной страницы
package api

import (
	"strconv"

	"github.com/sungora/app/cmd/app/config"
	"github.com/sungora/app/cmd/app/model"
	"github.com/sungora/app/cmd/app/model/mdlusers"
	"github.com/sungora/app/core"
	"github.com/sungora/app/lg"
)

func NewControlModel() core.Controller {
	return new(ControlModel)
}

type ControlModel struct {
	core.ControllerApi
}

// GET действие по умолчанию
func (control *ControlModel) GET() (err error) {
	// сессия
	var count int
	if control.Session.Get("count") != nil {
		count, _ = control.Session.Get("count").(int)
	}
	count += 1
	control.Session.Set("count", count)

	// работа с моделью
	u := mdlusers.New(10)
	control.Data = u

	model.DB.AutoMigrate(u)

	if err = u.Load(); err != nil {
		lg.Error(err)
	}

	invoiceID := uint64(1234)
	u.InvoiceID = &invoiceID
	u.Age = count
	u.IsOnline = true
	u.Status = "Пассив"

	if err = u.Save(); err != nil {
		return
	}

	return
}

// POST работа с моделью создание
func (control *ControlModel) POST() (err error) {

	u := mdlusers.New(0)
	control.Data = u

	invoiceID := uint64(8697)
	u.InvoiceID = &invoiceID
	name := "Вася пупкин"
	u.Nam = &name
	u.Age = 67
	json := `{"ServiceID": 24,"ClientID": 24,"InvoiceID": 24,"BillingCode": "NL","LocationCode": "NL"}`
	u.SampleJson = &json

	if err = u.Save(); err != nil {
		return
	}

	return
}

// PUT работа с моделью изменение
func (control *ControlModel) PUT() (err error) {
	u := mdlusers.New(0)
	control.Data = u

	if err = control.RW.RequestBodyDecodeJson(u); err != nil {
		return
	}

	if err = u.Save(); err != nil {
		return
	}

	return
}

// DELETE работа с моделью удаление
func (control *ControlModel) DELETE() (err error) {
	var ID uint64
	ID, _ = strconv.ParseUint(control.RW.RequestParams["userID"][0], 10, 64)
	u := mdlusers.New(ID)
	control.Data = u

	if err = u.Delete(); err != nil {
		return
	}

	return
}

// OPTIONS получение различных опций для текущего контроллера
func (control *ControlModel) OPTIONS() (err error) {

	// Выводимые данные
	control.Data = map[string]interface{}{
		"scenario form": mdlusers.Scenario.AdminForm, // сценарии для виджетов
		"config access": config.ServiceAccess,        // какая-то конфигурация
	}
	return
}
