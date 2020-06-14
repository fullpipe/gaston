package converter

import (
	"errors"

	"github.com/tidwall/gjson"
)

type Provider struct {
	factories map[string]ConverterFactory
}

func NewProvider() *Provider {
	p := Provider{
		factories: make(map[string]ConverterFactory),
	}

	p.Add("rename", NewRename)
	p.Add("overwite", NewOverwrite)

	return &p
}

func (p *Provider) Add(converterType string, factory ConverterFactory) {
	p.factories[converterType] = factory
}

func (p *Provider) Get(converterType string, json gjson.Result) (Converter, error) {
	factory, ok := p.factories[converterType]
	if !ok {
		return nil, errors.New("No such converter")
	}

	return factory(json)
}
