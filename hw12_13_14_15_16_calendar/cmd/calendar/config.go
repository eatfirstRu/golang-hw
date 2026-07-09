package main

import (
	"github.com/fixme_my_friend/hw12_13_14_15_16_calendar/internal/config"
)

func NewConfig(path string) (*config.Config, error) {
	return config.NewConfig(path)
}
