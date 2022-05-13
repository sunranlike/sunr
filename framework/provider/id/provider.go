package id

import (
	"github.com/sunranlike/hade/framework"
	"github.com/sunranlike/hade/framework/contract"
)

type HadeIDProvider struct {
}

// Register registe a new function for make a services instance
func (provider *HadeIDProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeIDService
}

// Boot will called when the services instantiate
func (provider *HadeIDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the services instantiate when first make or register
func (provider *HadeIDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HadeIDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this services
func (provider *HadeIDProvider) Name() string {
	return contract.IDKey
}
