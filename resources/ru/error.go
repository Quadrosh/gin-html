package resources

import "fmt"

// LinkIsOld Ссылка устарела или не актуальна
func LinkIsOld() string {
	return "Ссылка устарела или не актуальна"
}

// UserNotFound Пользователь не найден
func UserNotFound() string {
	return "Пользователь не найден"
}

// DatabaseSaveError - Ошибка сохранения данных в базу данных
func DatabaseSaveError() string {
	return "Ошибка сохранения данных в базу данных"
}

// CantGetFromSession - Can't get %s from seesion
func CantGetFromSession(key string) string {
	return fmt.Sprintf("Can't get %s from seesion", key)
}

// InvalidEmailOrPassword - Неверный email или пароль
func InvalidEmailOrPassword() string {
	return "Неверный email или пароль"
}

// InvalidPassword - Неверный пароль
func InvalidPassword() string {
	return "Неверный пароль"
}

// InvalidEmail - Неверный email
func InvalidEmail() string {
	return "Неверный email"
}

// Forbidden - Недостаточно полномочий для выполнения операции
func Forbidden() string {
	return "Недостаточно полномочий для выполнения операции"
}

// SystemError - Системная ошибка - обратитесь к администратору
func SystemError() string {
	return "Системная ошибка - обратитесь к администратору"
}

// InvalidID - Некорректный идентификатор
func InvalidID() string {
	return "Некорректный идентификатор"
}

// PageNotFound - Страница не найдена
func PageNotFound() string {
	return "Страница не найдена"
}

// UploadFailed - Ошибка загрузки
func UploadFailed() string {
	return "Ошибка загрузки"
}
