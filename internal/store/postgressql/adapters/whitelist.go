package adapters

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/store/postgressql/client"
	"fmt"
)

const (
	isIPExistWl = `SELECT exists(SELECT 1 FROM whitelist WHERE prefix = $1 AND mask = $2)`
	insertIPWl  = `INSERT INTO whitelist (prefix, mask) VALUES ($1, $2)`
	deleteIPWl  = `DELETE FROM whitelist WHERE prefix = $1 AND mask = $2`
	getIPListWl = `SELECT prefix, mask from whitelist`
)

type WhiteListRepository struct {
	client *client.PostgresSql
}

func NewWhiteListRepository(client *client.PostgresSql) *WhiteListRepository {
	return &WhiteListRepository{client: client}
}

func (r *WhiteListRepository) Add(prefix string, mask string) error {
	var isExist bool
	err := r.client.Db.QueryRow(isIPExistWl, prefix, mask).Scan(&isExist)
	if err != nil {
		return err
	}
	if isExist {
		return fmt.Errorf("this ip network already exist")
	}
	err = r.client.Db.QueryRow(insertIPWl, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *WhiteListRepository) Remove(prefix string, mask string) error {
	err := r.client.Db.QueryRow(deleteIPWl, prefix, mask).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *WhiteListRepository) Get() ([]entity.IpNetwork, error) {
	ipNetworkList := make([]entity.IpNetwork, 0, 5)
	err := r.client.Db.Select(&ipNetworkList, getIPListWl)
	if err != nil {
		return nil, err
	}
	return ipNetworkList, nil
}
