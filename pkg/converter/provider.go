package converter

import (
	"errors"

	"github.com/tidwall/gjson"
)

// Provider keeps ConverterFactories
type Provider struct {
	factories map[string]ConverterFactory
}

// NewProvider returns new Provider with buildin converters
func NewProvider() *Provider {
	p := Provider{
		factories: make(map[string]ConverterFactory),
	}

	p.Add("rename", NewRename)
	p.Add("overwite", NewOverwrite)
	p.Add("snakeCase", NewSnakeCase)
	p.Add("snake_case", NewSnakeCase)
	p.Add("remove", NewRemove)
	p.Add("delete", NewRemove)
	p.Add("castNumber", NewCastNumber)

	return &p
}

// Add adds ConverterFactory to Provider
func (p *Provider) Add(converterType string, factory ConverterFactory) {
	p.factories[converterType] = factory
}

// Get creates and returns Converter by its type and config
func (p *Provider) Get(converterType string, config gjson.Result) (Converter, error) {
	factory, ok := p.factories[converterType]
	if !ok {
		return nil, errors.New("No such converter")
	}

	return factory(config)
}
