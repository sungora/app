package message

import (
	"fmt"
)

func GetMessage(code int, params ...interface{}) (message string) {
	var ok bool
	if message, ok = messages[code]; ok == true {
		message = fmt.Sprintf(message, params...)
	} else if 0 < len(params) {
		if s, ok := params[0].(string); ok == true {
			message = fmt.Sprintf(s, params[1:]...)
		}
	}
	return
}

func SetMessage(msg map[int]string) {
	for k, v := range msg {
		messages[k] = v
	}
}

var messages = map[int]string{
	100: `Установка сессионной куки: [%s] = [%s]`,
	101: `Установка постоянной куки: [%s] = [%s]`,
	102: `Шаблон статуса ошибки [%d] не найден: %v`,
	103: `Неудалось прочитать JSON данные запроса [%s] [%s]`,
	104: `Ошибка разбора JSON: [%s] : %v`,
	105: `Ошибка формирования JSON: %v`,
	106: `+++ Status [%d] [%s]`,
	107: `+++ Выполнение редиректа на URL: [%s]`,
	108: `-- Некорректный формат свойства для БД: [%s].[%s]`,
	109: `Ошибка выполнения запроса [%s] : %v`,
	110: `Ошибка компиляции запроса [%s] : %v`,
	111: `Невозможно соединиться с БД Mysql: %v`,
	112: `Отсутсвует конфигурация Mysql № [%d]`,
	113: `Ключевое поле (свойство) [%s] отсутствует в структуре [%s]`,
	114: `Ошибка выполнения пакетного запроса (QueryByte)`,
	115: `Нет прав на метод запроса [%s] [%s]`,
	116: `Нет прав на метод запроса [%s] [%s], переход на авторизацию`,
	117: `Запоминаем куда пользователю вернуться [%s]`,
	118: `Ошибка удаления временного файла контента по ключу [%s] : %v`,
	119: `Ошибка общей проверки для сценария [%s] %v`,
	120: `Ошибка копирования типа в модели: %v`,
	121: `Сценарий [%s] для источника [%s] не найден`,
	123: `Найден ранее загруженный файл контента [%s]`,
	124: `Объекты типа [%s] отсутсвуют в памяти`,
	125: `Объекты типа [%s] должны храниться в срезе`,
	126: `Запуск контроллера [%d] [%s]`,
	127: `Ошибка контроллера [%d] [%s]`,
	128: `Ошибка выполнения (парсинг) шаблона [%s] %v`,
	129: `Статика: [%s]`,
	130: `Метод [%s] не поодерживается Uri [%s]`,
	131: `Найден URI: [id:%d]:[%s]`,
	132: `Токен [%s]`,
	134: `Пользователь определен как [%s] : [%d]`,
	135: `Ошибка разбора сегмента Uri [%s][%s] %v`,
	136: `Ошибка разбора сегмента Uri 'relationid' [%s] %v`,
	137: `Не указан шаблон почтового сообщения`,
	138: `Отсутствует файл шаблона [%s] %v`,
	139: `Сегмент uri отсутствует [%s][%s]`,
	143: `Завершение работы приложения`,
	144: `Ошибка получения запроса по индексу [%s]`,
	145: `Ошибка обновления контента контроллера [%s] %v`,
	146: `Ошибка обновления контента uri [%s] %v`,
	147: `Ошибка обновления дефолтового контент-шаблона uri [%s] [%s]`,
	148: `Ошибка открытия порта для сервера [%s]`,
	149: `Сервер запущен по адресу [%s]`,
	150: `Сервер остановлен по адресу [%s]`,
	151: `Сервер [%s] не был запущен`,
	152: `Объект отдающий данные должен быть передан по ссылке: [%s]`,
	153: `Объект принимающий данные должен быть передан по ссылке 'CopyTyp': [%s]`,
	154: `Объект отдающий данные не инициализирован [%s]`,
	155: `Объект принимающий данные не инициализирован [%s]`,
	156: `Объект не найден в БД [%s]`,
	157: `Объект принимающий данные должен быть передан по ссылке 'mysql.SelectSlice': [%s] [%s]`,
	158: `Неопределенный тип хранения данных 'mysql.LoadData' [%s]`,
	159: `Ошибка загрузки. Объект не найден в памяти. [%s] [%d]`,
	160: `Сценарий не найден: [%s] -> [%s]`,
	161: `Id пользователя разработчик инициализировано не верно`,
	162: `Email [%s] уже занят`,
	163: `Логин [%s] занят`,
	164: `Пароли [%s] != [%s] не совпадают`,
	165: `Id пользователя гость инициализировано не верно`,
	166: `Id группы разработчик инициализировано не верно`,
	167: `Id группы гость инициализировано не верно`,
	168: `Использование БД отключено либо не реализовано [%s]`,
	169: `Объект принимающий данные должен быть срезом 'mysql.SelectSlice': [%s] [%s]`,
	170: `Ошибочное свойство для загрузки из БД 'mysql.SelectSlice': [%s] [%s]`,
	171: `Ошибочное поле в запросе для загрузки из БД 'mysql.SelectSlice': [%s] [%s]`,
	172: `Неверный путь до контроллера: [%s]`,
	173: `Контроллер [%s/%s] отсутсвует`,
	174: `Контроллер [%s/%s] не имеет метода [%s]`,
	175: `Удаление куки: [%s]`,
	176: `Объект принимающий данные должен быть передан по ссылке 'mysql.Select': [%s] [%s]`,
	177: `Ошибочное свойство для загрузки из БД 'mysql.Select': [%s] [%s]`,
	178: `Ошибочное поле в запросе для загрузки из БД 'mysql.Select': [%s] [%s]`,
	179: `Пароль изменен для [%s]`,
	180: `Объект принимающий данные должен быть инициализирован 'mysql.Select': [%s] [%s]`,
	181: `Ошибочное свойство для загрузки из БД 'mysql.SelectData': [%s]`,
	182: `Тип свойства для загрузки таблиц целиком не поддерживается [%s]`,
	183: `Ошибочное поле в запросе для загрузки из БД 'mysql.SelectMap': [%s] [%s]`,
	301: `Переадресация на [%s]`,
	404: `Документ не найден`,
	403: `Доступ запрещен`,
	500: `Ошибка сервера`,
	300: `Объект отдающий данные должен быть передан по ссылке: [%s]`,
	381: `Объект принимающий данные должен быть передан по ссылке: [%s]`,
	302: `Объект отдающий данные не инициализирован [%s]`,
	303: `Объект принимающий данные не инициализирован [%s]`,
	304: `Сценарий не найден: [%s] -> [%s]`,
	305: `Неверный путь до контроллера: [%s]`,
	306: `Контроллер [%s/%s] отсутсвует`,
	307: `Контроллер [%s/%s] не имеет метода [%s]`,
	308: "У объекта [%s] отсутсвует свойство [%s]",
	309: "Cвойство [%s].[%s] обязательно для заполнения",
	310: `Найден ранее загруженный файл контента [%s]`,
	311: `Ошибка удаления временного файла контента по ключу [%s] : %v`,
	312: `Ошибка копирования типа в модели: %v`,
	313: `Сценарий [%s] для источника [%s] не найден`,
	314: `Ошибка общей проверки для сценария [%s] %v`,
	315: "Cвойство [%s].[%s] не может быть изменено",
	316: `Объекты типа [%s] отсутсвуют в памяти`,
	317: `Объекты типа [%s] должны храниться в срезе`,
	318: `Ошибка загрузки. Объект не найден в памяти. [%s] [%v]`,
	319: `Найден URI: [%d] : [%s]`,
	320: `Токен пользователя [%s]`,
	321: `Пользователь определен как [%s] : [%d]`,
	322: `Запуск контроллера [%d] [%s]`,
	323: `Ошибка контроллера [%d] [%s]`,
	324: `Ошибка выполнения шаблона (парсинг) [%s] %v`,
	325: `Ошибка обновления дефолтового контент-шаблона uri [%s] [%s]`,
	326: `Ошибка обновления контента контроллера [%s] %v`,
	327: `Ошибка открытия порта для сервера [%s]`,
	328: `Сервер остановлен по адресу [%s]`,
	329: `Сервер [%s] не был запущен`,
	330: `Сервер запущен по адресу [%s]`,
	331: `Id пользователя разработчик инициализировано не верно`,
	332: `Id пользователя гость инициализировано не верно`,
	333: `Id группы разработчик инициализировано не верно`,
	334: `Id группы гость инициализировано не верно`,
	335: `Несуществующее свойство '[%s.%s]'`,
	336: `Ошибка обновления контента uri [%s] %v`,
	337: `Токен не передан`,
	338: `Токен [%s] не верный`,
	339: `Пользователь не найден`,
	340: `Ошибка получения капчи %v`,
	341: `Ошибка отправки email: %v`,
	342: `Ошибка отправки регистрационного письма [%s]`,
	359: `Слишком частые запросы капчи`,
	360: `Требуется ввод капчи. Капча не указана`,
	361: `Капча не верна [%s]`,
	362: `Email [%s] уже занят`,
	363: `Логин [%s] уже занят`,
	364: `Пароли [%s] != [%s] не совпадают`,
	365: `Логин не указан`,
	367: `Пароль не верен`,
	368: `Пользователь заблокирован`,
	369: `Пользователь удален`,
	370: `Ошибка сохранения пользователя при авторизации`,
	376: `Email не указан`,
	378: `Ошибка изменение пароля для пользователя`,
	380: `Ошибка отправки письма с восстановлением [%s]`,
	382: `Ошибка копирование срезов (copy) [%s]`,
	393: `Объект не найден`,
	200: `+++ Status [%d] [%s]`,
	203: `Тип получаемых опций не реализован [%s]`,
	250: `Установка постоянной куки: [%s] = [%s]`,
	251: `Установка сессионной куки: [%s] = [%s]`,
	252: `Удаление куки: [%s]`,
	253: `Сегмент uri отсутствует [%s][%s]`,
	254: `Ошибка разбора сегмента Uri [%s][%s] %v`,
	450: `Неудалось прочитать JSON данные запроса [%s] [%s]`,
	451: `Ошибка разбора JSON: [%s] : %v`,
	452: `Ошибка формирования JSON: %v`,
	510: `Ошибка получение входных данных запроса`,
	520: `Не реализованный метод запроса внутри контроллера [%s]`,
	530: `Ошибка проверки по сценарию [%s]`,
	540: `Ошибка сохранения [%s] [%d]`,
	550: `Ошибка удаления [%s] [%d]`,
	560: `Ошибка добавления [%s]`,
	570: `Опции [%s] не найдены`,
	575: `Тип получаемых опций не указан`,
	580: `Ошибка сортировки [%s]`,
	590: `Ошибка загрузки файла %v`,
	800: `Отсутсвует конфигурация Mysql № [%d]`,
	801: `Ошибка компиляции запроса [%s] : %v`,
	802: `Ошибка выполнения запроса [%s] : %v`,
	804: `Ошибочное свойство для загрузки из БД: [%s] [%s]`,
	805: `Объект не найден в БД [%s]`,
	806: `Невозможно соединиться с БД Mysql конфиг № [%d] : %v`,
	807: `Объект принимающий данные должен быть срезом 'mysql.SelectSlice': [%s] [%s]`,
	808: `Ключевое поле (свойство) [%s] отсутствует в структуре [%s]`,
	809: `Ошибка выполнения пакетного запроса (QueryByte)`,
	810: `Ошибка получения запроса по индексу [%s]`,
	811: `Использование БД отключено либо не реализовано [web.Config.Main.UseDb = %d]`,
	812: `Неопределенный тип хранения данных 'mysql.SelectData' [%s]`,
	813: `Объект принимающий данные должен быть хешом 'mysql.SelectMap': [%s] [%s]`,
	814: `Проверка работоспособности с БД mysql завершилась с ошибкой: %v`,
}
