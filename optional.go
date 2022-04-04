package easyslice

type IOptional interface {
	IsPresent() bool
	Get() interface{}
	OrElse(other interface{}) interface{}
	OrElseGet(fn func() interface{}) interface{}
	Map(fn func(interface{}) interface{}) interface{}
}

type Optional struct {
	value interface{}
	found bool
}

func (o Optional) IsPresent() bool {
	return o.found
}

func (o Optional) Get() interface{} {
	return o.value
}

func (o Optional) OrElse(other interface{}) interface{} {
	if o.found {
		return o.Get()
	}
	return other
}

func (o Optional) OrElseGet(fn func() interface{}) interface{} {
	if o.found {
		return o.Get()
	}
	return fn()
}

func (o Optional) Map(fn func(interface{}) interface{}) interface{} {
	if o.found {
		return fn(o.Get())
	}
	return nil
}
