package resources

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
