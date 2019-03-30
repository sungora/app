package page

import (
	"bytes"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/sungora/app"
	"github.com/sungora/app/request"
	"github.com/sungora/app/tpl"

	"github.com/sungora/app/internal/sample/model"
)

// Main главная страница
func Main(w http.ResponseWriter, r *http.Request) {
	var err error

	var count int = 10
	// работа с моделью
	u := model.NewUser(0)
	invoiceID := uint64(8697)
	u.InvoiceID = &invoiceID
	name := "Вася пупкин"
	u.Nam = &name
	u.Age = count

	var (
		Variables     = make(map[string]interface{})
		Functions     = make(map[string]interface{})
		TplController = app.Cfg.DirWww + "/controllers/page/sample.html"
		TplLayout     = app.Cfg.DirWww + "/layout/index.html"
	)
	Variables["Header"] = "Head Control"
	Variables["User"] = u

	// шаблон контроллера
	var buf bytes.Buffer
	if buf, err = tpl.Compilation(TplController, Functions, Variables); err != nil {
		request.NewIn(w, r).Html("<H1>Internal Server Error</H1>", 500)
		return
	}
	if TplLayout == "" {
		request.NewIn(w, r).Html(buf.String())
		return
	}
	// шаблон макета
	Variables["Content"] = buf.String()
	if buf, err = tpl.Compilation(TplLayout, Functions, Variables); err != nil {
		request.NewIn(w, r).Html("<H1>Internal Server Error</H1>", 500)
		return
	}
	request.NewIn(w, r).Html(buf.String())
	return

}

// Api страница api
func Api(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).Html("IndexApi")
}

// ApiV1 страница api v1
func ApiV1(w http.ResponseWriter, r *http.Request) {
	request.NewIn(w, r).Html("PageApiV1")
}

// Sample Пример многоуровневого роутинга и GET параметры
func Sample(w http.ResponseWriter, r *http.Request) {
	testID := chi.URLParam(r, "testID")
	orderID := chi.URLParam(r, "orderID")
	pageID := chi.URLParam(r, "pageID")
	response := []string{
		testID,
		orderID,
		pageID,
	}
	request.NewIn(w, r).JsonApi200(response, 0, "OK")
}
