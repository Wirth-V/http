package modules

import (
	"log"
	"os"
)

// Используйте log.New() для создания логгера для записи информационных сообщений. Для этого нужно
// три параметра: место назначения для записи логов (os.Stdout), строка
// с префиксом сообщения (INFO или ERROR) и флаги, указывающие, какая
// дополнительная информация будет добавлена.
var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Создаем логгер для записи сообщений об ошибках таким же образом, но используем stderr как
// место для записи и используем флаг log.Lshortfile для включения в лог
// названия файла и номера строки где обнаружилась ошибка.
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// Включает в лог данные из тела ответов
var i int
var ResponseLog = log.New(os.Stdout, "", i)
