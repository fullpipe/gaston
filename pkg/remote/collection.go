package remote

// MethodCollection keeps collection of Methods
type MethodCollection struct {
	Methods []Method
}

// NewMethodCollection returns new collection
func NewMethodCollection() MethodCollection {
	return MethodCollection{
		Methods: []Method{},
	}
}

// Find returns method from collection by its name
func (c *MethodCollection) Find(methodName string) *Method {
	for _, method := range c.Methods {
		if method.Name == methodName {
			return &method
		}
	}

	return nil
}

// Add adds methos to collection
func (c *MethodCollection) Add(m Method) {
	c.Methods = append(c.Methods, m)
}

// Merge combines collection to a new one
func (c *MethodCollection) Merge(c2 MethodCollection) MethodCollection {
	return MethodCollection{
		Methods: append(c.Methods, c2.Methods...),
	}
}
