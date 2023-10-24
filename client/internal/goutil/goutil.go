package goutil

// AsPointer convert interface{} to *interface{}, if input is not nil
func AsPointer(obj interface{}) *interface{} {
	var ptr *interface{} = nil
	if obj != nil {
		ptr = &obj
	}
	return ptr
}

// NewBoolPointer return a pointer of the given boolean value, this function is needed because you cant do &true or &false in golang
func NewBoolPointer(value bool) *bool {
	return &value
}
