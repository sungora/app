// Контроллер
package users

import (
	"net/http"

	"github.com/sungora/app/request"
)

// Gets
func Gets(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// Post
func Post(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// Options
func Options(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// Get
func Get(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// Put
func Put(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// Delete
func Delete(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).JsonOk("Users", 0, "OK")
}

// // GET действие по умолчанию
// func ModelGET(ctx context.Context, rw *core.RW) (context.Context, *core.RW) {
// 	// var err error
// 	// // сессия
// 	// var count int
// 	// if control.Session.Get("count") != nil {
// 	//     count, _ = control.Session.Get("count").(int)
// 	// }
// 	// count += 1
// 	// control.Session.Set("count", count)
// 	//
// 	// // работа с моделью
// 	// u := mdlusers.New(10)
// 	// control.Data = u
// 	//
// 	// model.DB.AutoMigrate(u)
// 	//
// 	// if err = u.Load(); err != nil {
// 	//     lg.Error(err)
// 	// }
// 	//
// 	// invoiceID := uint64(1234)
// 	// u.InvoiceID = &invoiceID
// 	// u.Age = count
// 	// u.IsOnline = true
// 	// u.Status = "Пассив"
// 	//
// 	// if err = u.Save(); err != nil {
// 	//     lg.Error(err)
// 	//     return
// 	// }
// 	rw.ResponseJsonApi200("GET OK", 0, "")
// 	return ctx, rw
// }
//
// // POST работа с моделью создание
// func ModelPOST(ctx context.Context, rw *core.RW) (context.Context, *core.RW) {
// 	// var err error
// 	// u := mdlusers.New(0)
// 	// control.Data = u
// 	//
// 	// invoiceID := uint64(8697)
// 	// u.InvoiceID = &invoiceID
// 	// name := "Вася пупкин"
// 	// u.Nam = &name
// 	// u.Age = 67
// 	// json := `{"ServiceID": 24,"ClientID": 24,"InvoiceID": 24,"BillingCode": "NL","LocationCode": "NL"}`
// 	// u.SampleJson = &json
// 	//
// 	// if err = u.Save(); err != nil {
// 	//     lg.Error(err)
// 	//     return
// 	// }
// 	rw.ResponseJsonApi200("POST OK", 0, "")
// 	return ctx, rw
// }
//
// // PUT работа с моделью изменение
// func ModelPUT(ctx context.Context, rw *core.RW) (context.Context, *core.RW) {
// 	// var err error
// 	// u := mdlusers.New(0)
// 	// control.Data = u
// 	//
// 	// if err = control.RW.RequestBodyDecodeJson(u); err != nil {
// 	//     lg.Error(err)
// 	//     return
// 	// }
// 	//
// 	// if err = u.Save(); err != nil {
// 	//     lg.Error(err)
// 	//     return
// 	// }
// 	rw.ResponseJsonApi200("PUT OK", 0, "")
// 	return ctx, rw
// }
//
// // DELETE работа с моделью удаление
// func ModelDELETE(ctx context.Context, rw *core.RW) (context.Context, *core.RW) {
// 	// var err error
// 	// var ID uint64
// 	// ID, _ = strconv.ParseUint(control.RW.RequestParams["userID"][0], 10, 64)
// 	// u := mdlusers.New(ID)
// 	// control.Data = u
// 	//
// 	// if err = u.Delete(); err != nil {
// 	//     lg.Error(err)
// 	//     return
// 	// }
// 	rw.ResponseJsonApi200("DELETE OK", 0, "")
// 	return ctx, rw
// }
//
// // OPTIONS получение различных опций для текущего контроллера
// func ModelOPTIONS(ctx context.Context, rw *core.RW) (context.Context, *core.RW) {
// 	// // Выводимые данные
// 	// control.Data = map[string]interface{}{
// 	//     "scenario form": mdlusers.Scenario.AdminForm, // сценарии для виджетов
// 	//     "config access": config.SampleCustomConfig,   // какая-то конфигурация
// 	// }
// 	rw.ResponseJsonApi200("OPTIONS OK", 0, "")
// 	return ctx, rw
// }
