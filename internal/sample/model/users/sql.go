package users

// initialization Пользовательские зпапросы
func init() {
	Sql.GetListFilter = `
	SELECT * FROM users WHERE updated_at < NOW() - INTERVAL 1 HOUR LIMIT %d 
	`
}

type sqlTyp struct {
	GetListFilter string
}

// хранилище пользовательских запросов
var Sql = new(sqlTyp)
