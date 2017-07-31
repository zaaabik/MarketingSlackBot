package textConstants

const (
	RequestErrorText     = "Ooops! something went wrong, check internet connection and base api url"
	UserDoesNotExistText = "User does not exist!"
	ServerErrorText      = "Server error"
	CanceledEventText    = "cancel"
	ApproveEventText     = "added"
	UnknownCommand       = "unknown command\nwrite .help "
	Help                 = `.add [letters Count] letters [host_id] [provider]
	.get transaction count [host_id] [provider]
	.get customers count [host_id] [provider]
	example .add 1234 letters 13 radario`

	HostIdKey       = "host_id"
	ProviderKey     = "provider"
	LettersCountKey = "letters_count"

	AddUserLetterCountMethod      = "user/letters_count"
	GetCustomersCountMethod       = "customers/count"
	GetCustomersTransactionMethod = "customer_transactions/count"
	UpdateSendgridEmail           = "user/sendgrid"
)
