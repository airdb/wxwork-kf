package main

import (
	"log"
	"os"

	"github.com/airdb/sailor/deployutil"
	"github.com/airdb/sailor/faas"
	"github.com/airdb/wxwork-kf/internal/app"
	"github.com/airdb/wxwork-kf/internal/handler"

	// "github.com/airdb/wxwork-kf/pkg/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // docs is generated by Swag CLI, you have to import it.
)

// @title Airdb Serverlesss Example API
// @version 0.0.1
// @description This is a sample server Petstore server.
// @termsOfService https://airdb.github.io/terms.html

// @contact.name API Support
// @contact.url https://work.weixin.qq.com/kfid/kfc5fdb2e0a1f297753
// @contact.email info@airdb.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @info.x-logo.url: http://www.apache.org/licenses/LICENSE-2.0.html

// @host service-iw6drlfr-1251018873.sh.apigw.tencentcs.com
// @BasePath /wxkf
func main() {
	log.Println("start serverless:", deployutil.GetDeployStage(), os.Environ())

	app.InitWxWork()

	mux := chi.NewRouter()
	mux.Use(middleware.Logger)
	mux.Use(render.SetContentType(render.ContentTypeHTML))

	mux.Route("/wxkf", func(r chi.Router) {
		r.Get("/version", faas.HandleVersion)
		r.HandleFunc("/callback", handler.Callback)
		r.Get("/account/list", handler.AccountList)
		r.Put("/invite/image/{usedBy}", handler.InviteImageUpload)
	})

	// http.ListenAndServe(":3333", mux)
	faas.RunTencentChiWithSwagger(mux, "wxkf")
}
