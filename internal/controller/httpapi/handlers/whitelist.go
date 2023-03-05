package handlers

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

type WhiteList struct {
	service *service.WhiteList
	logger  *zap.SugaredLogger
}

func NewWhiteList(service *service.WhiteList, logger *zap.SugaredLogger) *WhiteList {
	return &WhiteList{service: service, logger: logger}
}

func (a *WhiteList) AddIP(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Add IP in whitelist by POST /auth/whitelist")
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
		if err.Error() == ipAlreadyExist.Error() {
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

func (a *WhiteList) RemoveIP(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Remove IP in whitelist by POST /auth/whitelist")
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

func (a *WhiteList) ShowIPList(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Info("Show IP list in whitelist by GET /auth/whitelist")
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
