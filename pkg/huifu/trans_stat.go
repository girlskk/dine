package huifu

type TransStat string

const (
	TransStatProcessing TransStat = "P"
	TransStatSuccess    TransStat = "S"
	TransStatFail       TransStat = "F"
)
