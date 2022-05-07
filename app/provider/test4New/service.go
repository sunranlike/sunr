package test4New

import "github.com/sunranlike/hade/framework"

type Test4NewService struct {
	container framework.Container
}

func NewTest4NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.Container)
	return &Test4NewService{container: container}, nil
}
func (s *Test4NewService) Foo() string {
	return ""
}
