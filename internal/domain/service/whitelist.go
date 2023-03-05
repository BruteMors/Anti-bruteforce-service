package service

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"go.uber.org/zap"
)

type WhiteListStore interface {
	Add(prefix string, mask string) error
	Remove(prefix string, mask string) error
	Get() ([]entity.IpNetwork, error)
}

type WhiteList struct {
	store  WhiteListStore
	logger *zap.SugaredLogger
}

func NewWhiteList(store WhiteListStore, logger *zap.SugaredLogger) *WhiteList {
	return &WhiteList{store: store, logger: logger}
}

func (b *WhiteList) AddIP(network entity.IpNetwork) error {
	b.logger.Infoln("Get prefix")
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = b.store.Add(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (b *WhiteList) RemoveIP(network entity.IpNetwork) error {
	b.logger.Infoln("Get prefix")
	prefix, err := GetPrefix(network.Ip, network.Mask)
	if err != nil {
		return err
	}
	err = b.store.Remove(prefix, network.Mask)
	if err != nil {
		return err
	}
	return nil
}

func (b *WhiteList) GetIPList() ([]entity.IpNetwork, error) {
	ipList, err := b.store.Get()
	if err != nil {
		return nil, err
	}
	return ipList, nil
}
