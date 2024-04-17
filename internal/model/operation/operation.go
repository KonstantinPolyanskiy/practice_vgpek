package operation

const (
	DataLayer    = "База данных"
	ServiceLayer = "Бизнес логика"
	HTTPLayer    = "REST API эндпоинты"
)

// Логгирование методов DAO
const (
	SaveActionDAO        = "сохранение действия в базу данных"
	SelectActionById     = "получение действия из базы данных по id"
	SoftDeleteActionById = "мягкое удаление действия по id"
	Insert               = "вставка"
	Select               = "получение"
	Update               = "обновление"
	ExecuteError         = "ошибка выполнения запроса"
	CollectError         = "ошибка приведения к структуре"
	SuccessfullyRecorded = "успешно сохранено в базу данных"
	SuccessfullyReceived = "успешно получено из базы данных"
	SuccessfullyUpdated  = "запись успешно обновлеа в базе данных"
)

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
	GetIssuedPracticeInfoByParams = "получение по параметрам информации по практическими заданиям"
	DownloadIssuedPractice        = "получение ссылки для загрузки практического задания"
)

const (
	UploadSolvedPracticeOperation = "добавление выполненной практической работы"
	GetSolvedPracticeInfoById     = "получение по id информации по выполненной практической работе"
	SetMarkSolvedPractice         = "выставление оценки выполненному практическому заданию"
)
