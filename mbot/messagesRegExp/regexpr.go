package messagesRegExp

const (
	AddLettersToUserRegExp    = `((<@\w+>\s*)+|(^\s*))(\.add (-)?\d* letters \w+ \w+\s*$)`
	GetTransactionCountRegExp = `((<@\w+>\s*)+|(^\s*))(\.get transaction count \w+ \w+\s*$)`
	GetCustomersCountRegExp   = `((<@\w+>\s*)+|(^\s*))(\.get customers count \w+ \w+\s*$)`
	UpdateSendgridEmailRegExp = `((<@\w+>\s*)+|(^\s*))(\.set \w+ \w+ sendgrid \w+\s*$)`
	ShowDbRegExp              = `\.show`
	DeleteDbRegExp            = `\.del`
	HelpRegExp                = `((<@\w+>\s*)+|(^\s*))\.help)`
	AllRegExp                 = `.*`
)
