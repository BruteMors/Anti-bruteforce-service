package handlers

import (
	"Anti-bruteforce-service/internal/domain/entity"
	"Anti-bruteforce-service/internal/domain/service"
	jsoniter "github.com/json-iterator/go"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
	"net/http"
)

type Bucket struct {
	service *service.Authorization
	logger  *zap.SugaredLogger
}

func NewBucket(service *service.Authorization, logger *zap.SugaredLogger) *Bucket {
	return &Bucket{service: service, logger: logger}
}

func (a *Bucket) ResetBucket(rw http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	a.logger.Infoln("Reset Bucket by POST /auth/bucket")
	initHeaders(rw)
	var request entity.Request
	err := jsoniter.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		a.logger.Infof("Invalid json received from client: %v", err)
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	request.Password = "empty"
	isValidate := ValidateRequest(request)
	if !isValidate {
		a.logger.Info("Invalid input request received from client")
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	isLoginReset := a.service.ResetLoginBucket(request.Login)
	if !isLoginReset {
		_, err = rw.Write([]byte("resetLogin=false\n"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err = rw.Write([]byte("resetLogin=true\n"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	isIpReset := a.service.ResetIpBucket(request.Ip)
	if !isIpReset {
		_, err = rw.Write([]byte("resetIp=false"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		_, err = rw.Write([]byte("resetIp=true"))
		if err != nil {
			a.logger.Infof("Troubles with response, err: %v", err)
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
