package interfaces

type Permission interface {
	ModelWithKey
	ModelWithNameDescription
	GetType() (t string)
	SetType(t string)
	GetTarget() (target []string)
	SetTarget(target []string)
	GetAllow() (allow []string)
	SetAllow(allow []string)
	GetDeny() (deny []string)
	SetDeny(deny []string)
}
