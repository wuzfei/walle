package field

const StatusEnable = 1
const StatusDisable = 2

type Status uint8

func (f Status) IsEnable() bool {
	return f == StatusEnable
}

func (f Status) IsDisable() bool {
	return f == StatusDisable
}
