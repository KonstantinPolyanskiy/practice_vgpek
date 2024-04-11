package operation

// Операции с ключем
const (
	NewKeyOperation         = "создание нового ключа"
	InvalidateKeyOperation  = "удаление ключа"
	GetKeysOperation        = "получение ключей"
	GetKeyByBodyOperation   = "получение ключа по телу"
	GetKeyByIdOperation     = "получение ключа по id"
	IncCountUsagesOperation = "увелечение счетчика регистраций"
)

// Операции с ролями
const (
	AddRoleOperation  = "добавление роли"
	GetRoleOperation  = "получение роли"
	GetRolesOperation = "получение ролей"
)

// Операции с доступами
const (
	AddPermissionOperation    = "добавление права действия в системе"
	GetPermissionOperation    = "получение доступов у роли"
	DeletePermissionOperation = "удаление права действия в системе"
)

// Операции с объектами действия
const (
	AddObjectOperation  = "добавление объекта действия в системе"
	GetObjectOperation  = "получение объекта действия"
	GetObjectsOperation = "получение объектов действий"
)

// Операции с действиями
const (
	AddActionOperation  = "добавление права действия в системе"
	GetActionOperation  = "получение права действия по id"
	GetActionsOperation = "получение прав действий по параметрам"
)

// Операции с авторизацией
const (
	RegistrationOperation = "регистрация пользователя"
	LoginOperation        = "вход пользователя"
)

// Операции с пользователем
const (
	NewPersonOperation = "добавление пользователя"
)

// Операции с аккаунтом
const (
	NewAccountOperation = "создание нового аккаунта"
	GetAccountOperation = "получение аккаунта по id"
)

// Операции с практическими заданиями
const (
	UploadIssuedPracticeOperation = "добавление практического задания"
	GetIssuedPracticeInfoById     = "получение по id информации по практическому заданию"
	DownloadIssuedPractice        = "получение ссылки для загрузки практического задания"
)
