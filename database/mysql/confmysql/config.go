package confmysql

type Mysql struct {
	Type     string // Тип подключения (socket | tcp (tcp по умолчанию)
	Socket   string // Путь к socket файлу
	Host     string // Хост базы данных (localhost - по умолчанию)
	Port     int64  // Порт подключения по протоколу tcp/ip (3306 по умолчанию)
	Name     string // Имя базы данных
	Login    string // Логин к базе данных
	Password string // Пароль к базе данных
	Charset  string // Кодировка данных (utf-8 - по умолчанию)
	TimeOut  int64  // Таймаут использования соединения (в секундах) (5 - по умолчанию)
	Updates  string // Путь где лежат обновления БД
	CntConn  int64  // Максимальное допустимое количество коннектов (50 - по умолчанию)
}
