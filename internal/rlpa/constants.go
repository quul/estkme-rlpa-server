package rlpa

const (
	TagMessageBox          = 0x00
	TagManagement          = 0x01
	TagDownloadProfile     = 0x02
	TagProcessNotification = 0x03
	TagReboot              = 0xFB
	TagClose               = 0xFC
	TagAPDULock            = 0xFD
	TagAPDU                = 0xFE
	TagAPDUUnlock          = 0xFF
)

type CallbackType string

const (
	CallbackTypeConnAdd     CallbackType = "add"
	CallbackTypeConnRemove               = "remove"
	CallbackTypeConnSetType              = "setType"
)

type ConnType string

const (
	ConnTypeManagement          ConnType = "management"
	ConnTypeProcessNotification ConnType = "processNotification"
	ConnTypeDownloadProfile     ConnType = "downloadProfile"
)
