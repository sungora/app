// запуск теста
// SET GOPATH=C:\Work\zegota
// go test -v lib/database/mysql | go test -v -bench . lib/database/mysql
package mysql_test

import (
	//"fmt"
	//"os"
	//"path/filepath"
	//"sort"
	//"strconv"
	"testing"
	//"time"
	"lib/database"
	"lib/logs"
)

func TestAR(t *testing.T) {
	logs.GoStart()
	var query = database.NewAr().SelectScenario(`Users`, `All`)
	var sql = query.From(`Users as z`).Where(`AND Id < 100`).Order(`Name ASC`).Get()
	logs.Notice(0, sql)
	logs.GoClose()
}
