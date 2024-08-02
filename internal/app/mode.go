package app

type ModeStruct string

func (m *ModeStruct) set(value string) {
	*m = ModeStruct(value)
}

func (m *ModeStruct) IsTest() bool {
	return *m == "test"
}

func (m *ModeStruct) IsProd() bool {
	return *m == "prod"
}
