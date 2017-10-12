package regularExpression

const (
	AddLettersToUserRegExp        = `((<@\w+>\s*)+|(^\s*))(\.add (-)?\d* letters \w+ \w+\s*$)`
	GetTransactionCountRegExp     = `((<@\w+>\s*)+|(^\s*))(\.get transaction count \w+ \w+\s*$)`
	GetCustomersCountRegExp       = `((<@\w+>\s*)+|(^\s*))(\.get customers count \w+ \w+\s*$)`
	UpdateSendgridEmailRegExp     = `((<@\w+>\s*)+|(^\s*))(\.set sendgrid email (<mailto:(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})\|(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})>) to \w+ \w+\s*$)`
	CreateScenarioByCompainRegExp = `((<@\w+>\s*)+|(^\s*))(\.create scenario \w+ <.*>\s*$)`
	ShowDbRegExp                  = `\.show`
	DeleteDbRegExp                = `\.del`
	HelpRegExp                    = `((<@\w+>\s*)+|(^\s*))(\.help)`
	LockUserExp                   = `((<@\w+>\s*)+|(^\s*))(\.lock \w+ \w+\s*$)`
	UnlockUserExp                 = `((<@\w+>\s*)+|(^\s*))(\.unlock \w+ \w+\s*$)`
	AllRegExp                     = `.*`
)
