package handlers

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

type BlackList struct {
	service *service.BlackList
	logger  *zap.SugaredLogger
}

func NewBlackList(service *service.BlackList, logger *zap.SugaredLogger) *BlackList {
	return &BlackList{service: service, logger: logger}
}

func (a *BlackList) AddIP(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Add IP in blacklist by POST /auth/blacklist")
	initHeaders(rw)
	var inputIp entity.IpNetwork
	err := jsoniter.NewDecoder(r.Body).Decode(&inputIp)
	if err != nil {
		a.logger.Infof("Invalid json received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateIP(inputIp)
	if !isValidate {
		a.logger.Info("Invalid input IP received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = a.service.AddIP(inputIp)
	if err != nil {
		if errors.Is(err, ipAlreadyExist) {
			a.logger.Info(err)
			rw.WriteHeader(http.StatusBadRequest)
			_, err = rw.Write([]byte(err.Error()))
			if err != nil {
				a.logger.Info(err)
				return
			}
			return
		}
		a.logger.Infof("Troubles with add ip: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func (a *BlackList) RemoveIP(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Remove IP in blacklist by POST /auth/blacklist")
	initHeaders(rw)
	var inputIp entity.IpNetwork
	err := jsoniter.NewDecoder(r.Body).Decode(&inputIp)
	if err != nil {
		a.logger.Infof("Invalid json received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	isValidate := ValidateIP(inputIp)
	if !isValidate {
		a.logger.Info("Invalid input IP received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	err = a.service.RemoveIP(inputIp)
	if err != nil {
		a.logger.Infof("Troubles with remove ip: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusNoContent)
}

func (a *BlackList) ShowIPList(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Show IP list in blacklist by GET /auth/blacklist")
	initHeaders(rw)
	ipList, err := a.service.GetIPList()
	if err != nil {
		a.logger.Infof("Troubles with show ip list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = jsoniter.NewEncoder(rw).Encode(ipList)
	if err != nil {
		a.logger.Infof("Troubles with encode ip list: %v", err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}
