package test4New

import (
	"github.com/sunranlike/sunr/framework"
)

type Test4NewProvider struct {
	framework.ServiceProvider
	c framework.Container
}

func (sp *Test4NewProvider) Name() string {
	return Test4NewKey
}
func (sp *Test4NewProvider) Register(c framework.Container) framework.NewInstance {
	return NewTest4NewService
}
func (sp *Test4NewProvider) IsDefer() bool {
	return false
}
func (sp *Test4NewProvider) Params(c framework.Container) []interface{} {
	return []interface{}{c}
}
func (sp *Test4NewProvider) Boot(c framework.Container) error {
	return nil
}
