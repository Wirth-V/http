package modules

import (
	"log"
	"os"
)

// Используйте log.New() для создания логгера для записи информационных сообщений
var InfoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

// Создаем логгер для записи сообщений об ошибках
var ErrorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

// Создаем логгер для записи сообщений об тел ответов
var i int
var ResponseLog = log.New(os.Stdout, "", i)
