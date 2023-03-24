package disk

var setupapi Setupapi

func init() {
	err := setupapi.Unmarshal("setupapi.dll")
	if err != nil {
		panic(err)
	}
}
