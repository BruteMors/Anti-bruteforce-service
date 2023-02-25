package adapters

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/store/postgressql/client"
	"fmt"
)

const (
	isIPExistBl = `SELECT exists(SELECT 1 FROM blacklist WHERE prefix = $1 AND mask = $2)`
	insertIPBl  = `INSERT INTO blacklist (prefix, mask) VALUES ($1, $2)`
	deleteIPBl  = `DELETE FROM blacklist WHERE prefix = $1 AND mask = $2`
	getIPListBl = `SELECT prefix, mask from blacklist`
)

type BlackListRepository struct {
	client *client.PostgresSql
}

func NewBlackListRepository(client *client.PostgresSql) *BlackListRepository {
	return &BlackListRepository{client: client}
}

func (r *BlackListRepository) Add(prefix string, mask string) error {
	var isExist bool
	err := r.client.Db.QueryRow(isIPExistBl, prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("this ip network already exist")
	}
	err = r.client.Db.QueryRow(insertIPBl, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *BlackListRepository) Remove(prefix string, mask string) error {
	err := r.client.Db.QueryRow(deleteIPBl, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *BlackListRepository) Get() ([]entity.IpNetwork, error) {
	ipNetworkList := make([]entity.IpNetwork, 0, 5)
	err := r.client.Db.Select(&ipNetworkList, getIPListBl)
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
