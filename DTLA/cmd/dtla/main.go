package main

import (
	"dtla/internal/sockets"
	"dtla/internal/util"
	"errors"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
	_ "modernc.org/sqlite"
)

func main() {
	var err error

	log.SetFlags(log.Ltime | log.Llongfile)

	var sstate ServerState
	err = (&sstate).Init()
	if err != nil {
		log.Fatal(err.Error())
	}

	sstate.mux.HandleFunc("GET /{$}", makeHandler(rootHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /view/{$}", makeHandler(viewAllHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /view/", makeHandler(viewHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /edit/", makeHandler(editHandler, &sstate, true))
	sstate.mux.HandleFunc("POST /save/", makeHandler(saveHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /tools/", makeHandler(toolsHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /new/", makeHandler(newHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /delete/", makeHandler(deleteHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /login", makeHandler(loginHandler, &sstate, true))
	sstate.mux.HandleFunc("POST /login", makeHandler(loginPostHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /logout", makeHandler(logoutHandler, &sstate, true))
	sstate.mux.HandleFunc("GET /LICENSE", makeHandler(licenseHandler, nil, false))
	sstate.mux.HandleFunc("GET /", makeHandler(getHandler, nil, false))
	sstate.mux.Handle("GET /api/sockets", websocket.Handler(sockets.Handler))

	go ListenShutdown(&sstate)

	var httpProtocol string
	if *sstate.TLS {
		httpProtocol = "HTTPS"
	} else {
		httpProtocol = "HTTP"
	}

	fmt.Printf("%s serveris palaists uz %s:%s\n", httpProtocol, *sstate.HttpIP, *sstate.HttpPort)
	fmt.Printf("Servē failus no '%s' un datubāzes '%s'\n", *sstate.PublicDir, *sstate.DBName)
	if *sstate.TLS {
		err = sstate.srv.ListenAndServeTLS(*sstate.TLSCert, *sstate.TLSPKey)
	} else {
		err = sstate.srv.ListenAndServe()
	}

	if !errors.Is(err, http.ErrServerClosed) {
		util.LogError(err.Error())
		return
	}

}
