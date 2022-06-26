package mgd

type CreateHook interface {
	Create() interface{}
}

type UpdateHook interface {
	Update(interface{}) interface{}
}

func Create(model interface{}) {
	var t interface{}
	if hook, ok := model.(CreateHook); ok {
		t = hook.Create()
	}
	if hook, ok := model.(UpdateHook); ok {
		hook.Update(t)
	}
}
func Update(model Model) {
	if hook, ok := model.(UpdateHook); ok {
		hook.Update(nil)
	}
}
