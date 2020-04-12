package method

type MethodCollection struct {
	Methods []Method
}

func (c *MethodCollection) Find(methodName string, version string) *Method {
	for _, method := range c.Methods {
		if method.Name == methodName && method.Version == version {
			return &method
		}
	}

	return nil
}
