package api_errors_msg

const (
	UserNotFound       = "Пользователь не найден"
	UserCreationFailed = "При создании пользователя произошла ошибка"
	UserBlockFailed    = "При блокировке пользователя произошла ошибка"
	UserUpdateFailed   = "При обновлении пользователя произошла ошибка"
	UserDeleteFailed   = "При удалении пользователя произошла ошибка"

	UsernameInvalid            = "Невалидное имя пользователя"
	UserEmailInvalid           = "Невалидный email пользователя"
	UserPasswordInvalid        = "Невалидный пароль пользователя"
	UserPasswordConfirmInvalid = "Пароли не совпадают"
)
