package textConstants

const (
	RequestErrorText     = "Ooops! something went wrong, check internet connection and base api url"
	UserDoesNotExistText = "User does not exist!"
	ServerErrorText      = "Server error"
	CanceledEventText    = "cancel"
	ApproveEventText     = "added"
	UnknownCommand       = "unknown command\nwrite .help "
	EmailChanged         = "email changed!"
	Help                 = `command list
	.add [letters Count] letters [host id] [provider]
	.get transaction count [host id] [provider]
	.get customers count [host id] [provider]
	.set sendgrid email [new email] to [host id] [provider]
	.create scenario [scenario name] [link to company]
	.lock [host id] [provider]
	.unlock [host id] [provider]

	example .add 1234 letters 13 radario`

	HostIdKey       = "host_id"
	ProviderKey     = "provider"
	LettersCountKey = "letters_count"
	EmailKey        = "email"
	ScenarioName    = "name"
	CampaignId      = "id"
	Lock            = "lock"

	AddUserLetterCountMethod       = "user/letters_count"
	GetCustomersCountMethod        = "customers/count"
	GetCustomersTransactionMethod  = "customer_transactions/count"
	UpdateSendgridEmailMethod      = "user/sendgrid"
	CreateScenarioByCampaignMethod = "scenario"
	LockUserMethod                 = "user/lock_state"
	UnlockUserMethod               = "user/lock_state"
)
