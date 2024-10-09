package common

type Bootable = func(cfg *Config, deps *Dependencies) error
