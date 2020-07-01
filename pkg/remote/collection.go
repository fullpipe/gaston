package remote

type MethodCollection struct {
	Methods []Method
}

func (c *MethodCollection) Find(methodName string) *Method {
	for _, method := range c.Methods {
		if method.Name == methodName {
			return &method
		}
	}

	return nil
}
func (c *MethodCollection) Add(m Method) {
	c.Methods = append(c.Methods, m)
}

func (c *MethodCollection) Merge(c2 MethodCollection) MethodCollection {
	return MethodCollection{
		Methods: append(c.Methods, c2.Methods...),
	}
}
