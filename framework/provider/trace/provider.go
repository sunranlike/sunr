package trace

import (
	"github.com/sunranlike/sunr/framework"
	"github.com/sunranlike/sunr/framework/contract"
)

type HadeTraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a services instance
func (provider *HadeTraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewHadeTraceService
}

// Boot will called when the services instantiate
func (provider *HadeTraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the services instantiate when first make or register
func (provider *HadeTraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *HadeTraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this services
func (provider *HadeTraceProvider) Name() string {
	return contract.TraceKey
}
