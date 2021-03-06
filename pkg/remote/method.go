package remote

import "github.com/fullpipe/gaston/pkg/converter"

// Method struct.
type Method struct {
	Host             string
	Version          string
	Name             string
	RemoteName       string
	Roles            []string
	ParamConverters  []converter.Converter
	ResultConverters []converter.Converter
}

// IsGranted check is it possible to use Method with provided roles
func (m *Method) IsGranted(roles []string) bool {
	if len(m.Roles) == 0 {
		return true
	}

	for _, role := range roles {
		for _, reqRole := range m.Roles {
			if reqRole == role {
				return true
			}
		}
	}

	return false
}
