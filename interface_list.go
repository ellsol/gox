package gox

type InterfaceList struct {
	values []interface{}
}

func NewInterfaceList() *InterfaceList {
	return &InterfaceList{
		values: make([]interface{}, 0),
	}
}

func (it *InterfaceList) Add(val interface{}) *InterfaceList {
	it.values = append(it.values, val)
	return it
}

func (it *InterfaceList) Values() []interface{} {
	return it.values
}


// String

type StringList struct {
	values []string
}

func NewStringList() *StringList {
	return &StringList{
		values: make([]string, 0),
	}
}

func (it *StringList) Add(val string) *StringList {
	it.values = append(it.values, val)
	return it
}

func (it *StringList) Values() []string {
	return it.values
}
