package status

//Status contains a status code and status information
type Status struct {
	StatusCode int
	StatusInfo string
}

// --- Status Codes

/*
Status codes always have 4 characters:
1. Status type
	'1': Status OK
	'2': Status in progress
	'3': Input Error
	'4': Internal Error
	'5': Authorization Error

2. Status Area
	'0': General; not further specified
	'1': File Storage
	'2': Database: General
	'3': Database: StorageNode
	'4': Database: CoordinatorNode
	'5': Networking: General
	'6': Networking: StorageNode
	'7': Networking: CoordinatorNode
	'8': JobQueue
	'9': -- unused --

3. & 4.: Status ID
*/

const OK int = 1000

const InProgress int = 2000

const GenericInputError int = 3000

const GenericInternalError int = 4000

const SettingsReadError int = 4100
const SettingsWriteError int = 4101

const DBPrepareError int = 4200
const DBWriteError int = 4201
const DBReadError int = 4202
const DBIdConflict int = 4203
const DBOpenError int = 4204
const DBCloseError int = 4205
const DBStructureError int = 4206

const SNDBPrepareError int = 4300
const SNDBWriteError int = 4301
const SNDBReadError int = 4302
const SNDBIdConflict int = 4311

const CNDBPrepareError int = 4400
const CNDBWriteError int = 4401
const CNDBReadError int = 4402
const CNDBIdConflict int = 4410

const NetworkingBadNodeType int = 4501

const SNNetworkingOutgoingRequestError int = 4601
const SNNetworkingReadingResponseError int = 4602

const CNNetworkingOutgoingRequestError int = 4701
const CNNetworkingReadingResponseError int = 4702

const JQTooManyWorkers int = 4800
const JQQueueTooLong int = 4801
