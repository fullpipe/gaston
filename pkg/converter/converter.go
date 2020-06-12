package converter

type Converter interface {
	Convert(jsonData string) (string, error)
}
