package mdlusers

// initialization Пользовательские зпапросы
func init() {
	sql.GetListFilter = `
	SELECT * FROM users WHERE updated_at < NOW() - INTERVAL 1 HOUR LIMIT %d 
	`
}

type sqlTyp struct {
	GetListFilter string
}

// хранилище пользовательских запросов
var sql = new(sqlTyp)
