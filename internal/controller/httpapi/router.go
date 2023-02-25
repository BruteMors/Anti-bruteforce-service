package httpapi

import (
	"Anti-bruteforce-service/internal/controller/httpapi/handlers"
	"github.com/julienschmidt/httprouter"
	"go.uber.org/zap"
)

type ApiRouter struct {
	router    *httprouter.Router
	auth      *handlers.Authorization
	blackList *handlers.BlackList
	whiteList *handlers.WhiteList
	bucket    *handlers.Bucket
	logger    *zap.SugaredLogger
}

func NewRouter(auth *handlers.Authorization, blackList *handlers.BlackList, whiteList *handlers.WhiteList, bucket *handlers.Bucket, logger *zap.SugaredLogger) *ApiRouter {
	router := httprouter.New()
	return &ApiRouter{
		router:    router,
		auth:      auth,
		blackList: blackList,
		whiteList: whiteList,
		bucket:    bucket,
		logger:    logger,
	}
}

func (r *ApiRouter) RegisterRoutes() {
	r.router.POST("/auth/check", r.auth.TryAuthorization)
	r.router.DELETE("/auth/bucket", r.bucket.ResetBucket)
	r.router.POST("/auth/blacklist", r.blackList.AddIP)
	r.router.DELETE("/auth/blacklist", r.blackList.RemoveIP)
	r.router.GET("/auth/blacklist", r.blackList.ShowIPList)
	r.router.POST("/auth/whitelist", r.whiteList.AddIP)
	r.router.DELETE("/auth/whitelist", r.whiteList.RemoveIP)
	r.router.GET("/auth/whitelist", r.whiteList.ShowIPList)
}

func (r *ApiRouter) GetRouter() *httprouter.Router {
	return r.router
}
