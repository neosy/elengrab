package errorx

import "github.com/neosy/elengrab/internal/pkg/errorx/exceptionx"

type factoryImpl struct{}

func init() {
	exceptionx.RegisterFactory(factoryImpl{})
}

func (f factoryImpl) New(text string, args ...any) error {
	return New(text, args...)
}

func (f factoryImpl) NewFromDomainException(ex exceptionx.DomainException, args ...any) error {
	return NewFromDomainException(ex, args)
}

func (f factoryImpl) NewFromException(ex exceptionx.Exception, args ...any) error {
	return NewFromException(ex, args)
}
