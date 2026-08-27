package api

import "github.com/gin-gonic/gin"

// TODO completare e suddividere per gruppi; inserire i middleware per la gestione delle autorizzazioni
func Routes(ph RouteHandler, router *gin.Engine) {

	router.GET("/", HomePage)

	router.POST("/login", ph.HandleLogin)

	authRouter := router.Group("/")
	authRouter.Use(RbacHandler())
	authRouter.GET("/scuola", ph.GetScuola)
	authRouter.GET("/scuola/sala/:id", ph.GetSala)
	authRouter.GET("/scuola/sale", ph.GetSale)
	authRouter.POST("/scuola/sala", ph.PostSala)
	authRouter.PUT("/scuola/sala", ph.PutSala)
	authRouter.DELETE("/scuola/sala/:id", ph.DeleteSala)

}
