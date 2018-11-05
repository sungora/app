package users

import (
	"gopkg.in/sungora/app.v1/core"
)

// GetHourlyAfter получение списка услуг для почасовой тарификации
func (self *User) GetListFilter(params ...interface{}) (users []*User) {
	core.DB.Raw(sql.GetListFilter, params...).Scan(&users)
	return
}
