package interfaces

type Filter interface {
	GetIsOr() (isOr bool)
	SetIsOr(isOr bool)
	GetConditions() (conditions []FilterCondition)
	SetConditions(conditions []FilterCondition)
	IsNil() (ok bool)
}
