package config

import (
	"net/http"

	"zhongxuqi/lowtea/app/server/handler"
)

// InitRouter init the router of server
func InitRouter(mainHandler *handler.MainHandler) {

	//---------------------------------
	//
	// init openapi handlers
	//
	//---------------------------------
	openAPIHandler := http.NewServeMux()
	openAPIHandler.HandleFunc("/openapi/login", mainHandler.Login)
	openAPIHandler.HandleFunc("/openapi/register", mainHandler.Register)
	openAPIHandler.HandleFunc("/openapi/logout", mainHandler.Logout)
	mainHandler.Mux.HandleFunc("/openapi/", func(w http.ResponseWriter, r *http.Request) {
		openAPIHandler.ServeHTTP(w, r)
	})

	//---------------------------------
	//
	// init api handlers
	//
	//---------------------------------
	apiHandler := http.NewServeMux()
	userHandler := http.NewServeMux()
	adminHandler := http.NewServeMux()

	// setup /api/ handler
	mainHandler.Mux.HandleFunc("/api/", func(w http.ResponseWriter, r *http.Request) {

		// check cookie
		err := mainHandler.CheckSession(w, r)
		if err != nil {
			http.Error(w, "[Handle /api/user/] "+err.Error(), 400)
			return
		}

		apiHandler.ServeHTTP(w, r)
	})

	// setup /api/user/ handler
	userHandler.HandleFunc("/api/user/userinfo", mainHandler.GetUserInfo)
	userHandler.HandleFunc("/api/user/users", mainHandler.GetUsers)
	apiHandler.HandleFunc("/api/user/", func(w http.ResponseWriter, r *http.Request) {
		userHandler.ServeHTTP(w, r)
	})

	// setup /api/admin/ handler
	adminHandler.HandleFunc("/api/admin/registers", mainHandler.GetRegisters)
	adminHandler.HandleFunc("/api/admin/register", mainHandler.ActionRegister)
	apiHandler.HandleFunc("/api/admin/", func(w http.ResponseWriter, r *http.Request) {

		// check permission
		err := mainHandler.CheckAdmin(r)
		if err != nil {
			http.Error(w, "check permission error: "+err.Error(), 400)
			return
		}

		adminHandler.ServeHTTP(w, r)
	})

	// init web file handler
	mainHandler.Mux.Handle("/", http.FileServer(http.Dir("../front/dist")))
}