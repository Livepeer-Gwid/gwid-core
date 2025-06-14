package di

import "gwid.io/gwid-core/internal/config"

type Container struct{}

func NewContainer(conf *config.Config) *Container {
	return &Container{}
}
