package router

import (
	"net/http"

	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/bulletin-board-service/api/handler"
	"github.com/Flo0807/hsfl-master-ai-cloud-engineering/lib/router"
)

type Router struct {
	router http.Handler
}

func NewRouter(
	healthHandler *handler.HealthHandler,
	postHandler *handler.PostHandler,
) *Router {
	r := router.New()
	r.GET("/bulletin-board/health", healthHandler.Health)
	r.GET("/bulletin-board/posts", postHandler.GetPosts)
	r.GET("/bulletin-board/posts/:id", postHandler.GetPost)
	r.POST("/bulletin-board/posts", postHandler.CreatePost, router.JWTAuthMiddleware)
	r.PUT("/bulletin-board/posts/:id", postHandler.UpdatePost)
	r.DELETE("/bulletin-board/posts/:id", postHandler.DeletePost)

	return &Router{r}
}

func (handler *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.router.ServeHTTP(w, r)
}
