package controllers

func NewActionControllerDelegate(id ControllerId, actions []Action) (d *ActionControllerDelegate) {
	return &ActionControllerDelegate{
		id:      id,
		actions: actions,
	}
}

type ActionControllerDelegate struct {
	id      ControllerId
	actions []Action
}

func (ctr *ActionControllerDelegate) Actions() (actions []Action) {
	return ctr.actions
}
