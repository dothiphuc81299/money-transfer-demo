package transaction

type Type int

const (
	DepositType Type = iota + 1
	WithdrawalType
	TransferType
)
