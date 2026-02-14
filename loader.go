package redis

import (
	"github.com/webcore-go/webcore/infra/config"
	"github.com/webcore-go/webcore/port"
)

type MemoryLoader struct {
	Memory *MemoryCache
	name   string
}

func (a *MemoryLoader) SetName(name string) {
	a.name = name
}

func (a *MemoryLoader) Name() string {
	return a.name
}

func (l *MemoryLoader) Init(args ...any) (port.Library, error) {
	config := args[0].(config.MemoryConfig)
	memory, err := NewMemoryCache(config)
	if err != nil {
		return nil, err
	}

	err = memory.Install(args...)
	if err != nil {
		return nil, err
	}

	memory.Connect()

	l.Memory = memory
	return memory, nil
}
