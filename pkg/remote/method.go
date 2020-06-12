package remote

import "github.com/fullpipe/gaston/pkg/converter"

type Method struct {
	Host             string
	Version          string
	Name             string
	Rename           string
	Roles            []string
	ParamConverters  []converter.Converter
	ResultConverters []converter.Converter
}

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
