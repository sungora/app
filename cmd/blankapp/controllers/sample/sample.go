// Контроллер главной страницы
package sample

import (
	"PKGAPPNAME/config"
	"PKGAPPNAME/models"

	"gopkg.in/sungora/app.v1/core"
)

type ControlSample struct {
	core.Controller
}

func (self *ControlSample) GET() (err error) {

	// сессия
	var count int
	if self.Session.Get("count") != nil {
		count, _ = self.Session.Get("count").(int)
	}
	count += 1
	self.Session.Set("count", count)

	// работас моделью
	u := models.NewUser(0)
	core.DB.AutoMigrate(u)
	u.Names = "Вася пупкин"
	u.Cnt = count
	u.Save()

	self.Data = []interface{}{
		u,
		config.ServiceAccess,
	}
	return
}
