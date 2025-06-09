package api_errors_msg

const (
	AuthInvalidRefreshToken        = "Невалидный refresh токен"
	AuthInvalidAccessToken         = "Невалидный access токен"
	AuthInvalidParams              = "Невалидный логин или пароль"
	JwtTokenFailed                 = "Не удалось получить токен"
	PermissionsDenied              = "Нет прав"
	MetadataNotProvided            = "Не предоставлена метадата"
	AuthorizationHeaderNotProvided = "Заголовок авторизации не передан"
	AuthorizationHeaderInvalid     = "Не валидный заголовок авторизации, не содержит префикс Bearer"
)
