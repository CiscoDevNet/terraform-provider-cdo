package goutil

// AsPointer convert interface{} to *interface{}, if input is not nil
func AsPointer(obj interface{}) *interface{} {
	var ptr *interface{} = nil
	if obj != nil {
		ptr = &obj
	}
	return ptr
}

func NewBoolPointer(value bool) *bool {
	return &value
}
