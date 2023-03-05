package service

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"go.uber.org/zap"
)

type BlackListStore interface {
	Add(prefix string, mask string) error
	Remove(prefix string, mask string) error
	Get() ([]entity.IpNetwork, error)
}

type BlackList struct {
	store  BlackListStore
	logger *zap.SugaredLogger
}

func NewBlackList(store BlackListStore, logger *zap.SugaredLogger) *BlackList {
	return &BlackList{store: store, logger: logger}
}

func (b *BlackList) AddIP(network entity.IpNetwork) error {
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

func (b *BlackList) RemoveIP(network entity.IpNetwork) error {
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

func (b *BlackList) GetIPList() ([]entity.IpNetwork, error) {
	ipList, err := b.store.Get()
	if err != nil {
		return nil, err
	}
	return ipList, nil
}
