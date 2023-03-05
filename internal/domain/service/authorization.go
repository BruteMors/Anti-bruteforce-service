package service

import (
	"Anti-bruteforce-service/internal/config"
	"Anti-bruteforce-service/internal/domain/entity"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
	"time"
)

type Authorization struct {
	ipBucketStorage       map[string]*RateLimiterWithLastEventTime
	loginBucketStorage    map[string]*RateLimiterWithLastEventTime
	passwordBucketStorage map[string]*RateLimiterWithLastEventTime
	blackList             *BlackList
	whiteList             *WhiteList
	cfg                   *config.Config
	logger                *zap.SugaredLogger
}

func NewAuthorization(blackList *BlackList, whiteList *WhiteList, cfg *config.Config, logger *zap.SugaredLogger) *Authorization {
	ipBuckets := make(map[string]*RateLimiterWithLastEventTime)
	loginBuckets := make(map[string]*RateLimiterWithLastEventTime)
	passwordBuckets := make(map[string]*RateLimiterWithLastEventTime)
	auth := Authorization{ipBucketStorage: ipBuckets, loginBucketStorage: loginBuckets, passwordBucketStorage: passwordBuckets, blackList: blackList, whiteList: whiteList, cfg: cfg, logger: logger}
	go auth.deleteUnusedBucket()
	return &auth
}

func (a *Authorization) TryAuthorization(request entity.Request) (bool, error) {
	a.logger.Infoln("Check ip in blacklist")
	ipNetworkList, err := a.blackList.GetIPList()
	if err != nil {
		return false, err
	}
	isIpInBlackList, err := a.checkIpByNetworkList(request.Ip, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIpInBlackList {
		return false, nil
	}

	a.logger.Infoln("Check ip in whitelist")
	ipNetworkList, err = a.whiteList.GetIPList()
	if err != nil {
		return false, err
	}
	isIpInWhiteList, err := a.checkIpByNetworkList(request.Ip, ipNetworkList)
	if err != nil {
		return false, err
	}
	if isIpInWhiteList {
		return true, nil
	}

	a.logger.Infoln("Check ip in bucket")
	isAllow := true
	allow := a.tryGetPermissionInBucket(a.ipBucketStorage, request.Ip, a.cfg.Bucket.IpLimit)
	if !allow {
		isAllow = allow
	}

	a.logger.Infoln("Check password in bucket")
	allow = a.tryGetPermissionInBucket(a.passwordBucketStorage, request.Password, a.cfg.Bucket.PasswordLimit)
	if !allow {
		isAllow = allow
	}

	a.logger.Infoln("Check login in bucket")
	allow = a.tryGetPermissionInBucket(a.loginBucketStorage, request.Login, a.cfg.Bucket.LoginLimit)
	if !allow {
		isAllow = allow
	}

	return isAllow, nil
}

func (a *Authorization) checkIpByNetworkList(ip string, ipNetworkList []entity.IpNetwork) (bool, error) {
	for _, network := range ipNetworkList {
		prefix, err := GetPrefix(ip, network.Mask)
		if err != nil {
			return false, err
		}
		if prefix == network.Ip {
			return true, nil
		}
	}
	return false, nil
}

func (a *Authorization) newBucket(limit int) *RateLimiterWithLastEventTime {
	limiter := NewLimiter(rate.Limit(float64(limit)/time.Duration.Seconds(60*time.Second)), limit)
	return limiter
}

func (a *Authorization) tryGetPermissionInBucket(bucketStorage map[string]*RateLimiterWithLastEventTime, requestValue string, limit int) bool {
	limiter, ok := bucketStorage[requestValue]
	if !ok {
		bucketStorage[requestValue] = a.newBucket(limit)
		allow := bucketStorage[requestValue].Allow()
		return allow
	}
	allow := limiter.Allow()
	return allow
}

func (a *Authorization) ResetLoginBucket(login string) bool {
	_, ok := a.loginBucketStorage[login]
	if !ok {
		return false
	}
	delete(a.loginBucketStorage, login)
	return true
}

func (a *Authorization) ResetIpBucket(ip string) bool {
	_, ok := a.ipBucketStorage[ip]
	if !ok {
		return false
	}
	delete(a.ipBucketStorage, ip)
	return true
}

func (a *Authorization) deleteUnusedBucket() {
	ticker := time.NewTicker(60 * time.Second)
	for {
		<-ticker.C
		for ip, limiter := range a.ipBucketStorage {
			if time.Since(limiter.LastEvent) > time.Duration(a.cfg.Bucket.ResetBucketInterval)*time.Second {
				delete(a.ipBucketStorage, ip)
			}
		}

		for login, limiter := range a.loginBucketStorage {
			if time.Since(limiter.LastEvent) > time.Duration(a.cfg.Bucket.ResetBucketInterval)*time.Second {
				delete(a.loginBucketStorage, login)
			}
		}

		for password, limiter := range a.passwordBucketStorage {
			if time.Since(limiter.LastEvent) > time.Duration(a.cfg.Bucket.ResetBucketInterval)*time.Second {
				delete(a.passwordBucketStorage, password)
			}
		}
	}
}
