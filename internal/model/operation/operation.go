package operation

const Operation = "операция"

const ParseToken = "парсинг JWT токена"

// Логирование методов REST
const (
	EndpointAddr    = "адрес"
	RegistrationReq = "запрос на регистрацию пользователя"
	DecodeError     = "ошибка парсинга json"
	ValidateError   = "ошибка валидации"
)

// Логирование методов DAO заданных практических
const (
	SaveIssuedPracticeDAO = "сохранение заданного практического задания в базе данных"
)

// Логирование методов DAO решенных практических
const (
	SaveSolvedPracticeDAO        = "сохранение решенного практического задания в базе данных"
	GetSolvedPracticeInfoByIdDAO = "получение решенного практического задания по id в базе данных"
	UpdateSolvedPracticeDAO      = "обновление решенного практического задания в базе данных"
)

// Логирование методов DAO доступов
const (
	SavePermissionsDAO    = "сохранение доступа в базе данных"
	SelectPermByRoleIdDAO = "получение доступов по id роли из базы данных"
)

// Логирование методов DAO пользователя
const (
	SavePersonDAO          = "сохранение пользователя в базу данных"
	SelectPersonByUIIDDAO  = "получение пользователя из базы данных по uuid"
	SelectPersonByAccIdDAO = "получение пользователя из базы данных по id аккаунта"
	SoftDeletePersonByUUID = "мягкое удаление пользователя по uuid"
)

// Логирование методов DAO аккаунта
const (
	SaveAccountDAO            = "сохранение аккаунта в базу данных"
	SelectAccountByIdDAO      = "получение аккаунта по id из базы данных"
	SelectAccountsByParamsDAO = "получение аккаунтов по параметрам из базы данных"
	SelectAccountByLoginDAO   = "получение аккаунта по логину из базы данных"
	SoftDeleteAccountByIdDAO  = "мягкое удаление аккаунта по id"
	HardDeleteAccountByIdDAO  = "жесткое удаление аккаунта по id"
)

// Логирование методов DAO действий
const (
	SaveActionDAO        = "сохранение действия в базу данных"
	SelectActionById     = "получение действия из базы данных по id"
	SelectActionByParams = "получение действий из базы данных по параметрам"
	SoftDeleteActionById = "мягкое удаление действия по id"
)

// Логирование методов DAO ключей
const (
	SaveKeyDAO         = "сохрание ключа в базу данных"
	SelectKeyById      = "получение ключа из базы данных по id"
	SelectKeyByBody    = "получение ключа из базы данных по телу"
	SelectKeysByParams = "получение ключей из базы данных по параметрам"
	UpdateKeyDAO       = "обновление полей ключа в базе данных"
	SoftDeleteKeyById  = "мягкое удаление ключа по id"
)

// Логирование методов Service ключей
const (
	InvalidateKey = "инвалидирование ключа регистрации"
	IncrementKey  = "инкрементирование регистраций"
)

// Логирование методов DAO ролей
const (
	SaveRoleDAO        = "сохранение роли в базу данных"
	SelectRoleById     = "получение роли из базы данных по id"
	SelectRoleByParams = "получение ролей из базы данных по параметрам"
	SoftDeleteRoleById = "мягкое удаление роли по id"
)

// Логирование методов DAO объектов
const (
	SaveObjectDAO        = "сохранение объекта в базу данных"
	SelectObjectById     = "получение объекта из базы данных по id"
	SoftDeleteObjectById = "мягкое удаление объекта по id"
	SelectObjectByParams = "получение объектов из базы данных по параметрам"
)

// Логирование методов DAO (общее)
const (
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
	AddPermissionOperation      = "добавление права действия в системе"
	GetPermissionOperation      = "получение доступов у роли"
	DeletePermissionOperation   = "удаление права действия в системе"
	GetPermByAccountIdOperation = "получение доступов по id аккаунта"
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
	NewAccountOperation          = "создание нового аккаунта"
	GetAccountOperation          = "получение аккаунта по id"
	GetAccountsByParamsOperation = "получение аккаунтов по параметрам"
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
