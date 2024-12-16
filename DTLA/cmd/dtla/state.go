package main

import (
	"context"
	"database/sql"
	"dtla/internal/util"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"time"
)

type ServerState struct {
	srv       http.Server
	mux       *http.ServeMux
	HttpIP    *string
	HttpPort  *string
	TLS       *bool
	TLSCert   *string
	TLSPKey   *string
	DBName    *string
	DB        *sql.DB
	PublicDir *string
	TmplDir   *string
	Tmpl      *template.Template
	Verbose   *bool
}

func (s *ServerState) Init() error {
	var err error

	// filepath.Clean() twice, first for help messages and second for flag value changes
	*s = ServerState{
		HttpIP:    flag.String("host", "127.0.0.1", "IP adrese uz kuras klausīties HTTP vaicājumus"),
		HttpPort:  flag.String("port", "30000", "Ports uz kura klausīties HTTP vaicājumus"),
		TLS:       flag.Bool("tls", true, "Vai klausīties izmantojot TLS"),
		TLSCert:   flag.String("cert", "cert", "TLS sertifikāts"),
		TLSPKey:   flag.String("key", "pkey", "TLS privātā atslēga"),
		DBName:    flag.String("db", filepath.Clean("db"), "Datubāzes fails"),
		PublicDir: flag.String("public", filepath.Clean("public"), "Publisko failu direktorija/folderis ar HTML, CSS, JavaScript, utt."),
		TmplDir:   flag.String("tmpl", filepath.Clean("public/tmpl"), "Veidņu direktorija/folderis ar veidnēm, ko izmanto lai ģenerētu HTML saturu"),
		Verbose:   flag.Bool("v", false, "Vairāk info"),
	}
	flag.Parse()

	// Absolute because of os.Chdir() later
	// Abs() calls Clean()
	*s.TLSCert, err = filepath.Abs(*s.TLSCert)
	if err != nil {
		util.LogFatal(err.Error())
	}

	*s.TLSPKey, err = filepath.Abs(*s.TLSPKey)
	if err != nil {
		util.LogFatal(err.Error())
	}

	*s.DBName, err = filepath.Abs(*s.DBName)
	if err != nil {
		util.LogFatal(err.Error())
	}

	*s.PublicDir, err = filepath.Abs(*s.PublicDir)
	if err != nil {
		util.LogFatal(err.Error())
	}

	*s.TmplDir, err = filepath.Abs(*s.TmplDir)
	if err != nil {
		util.LogFatal(err.Error())
	}

	err = os.Chdir(*s.PublicDir)
	if err != nil {
		util.LogFatal(err.Error())
	}

	s.mux = http.NewServeMux()
	s.srv = http.Server{
		Addr:    *s.HttpIP + ":" + *s.HttpPort,
		Handler: s.mux,
	}

	s.Tmpl, err = template.ParseGlob(*s.TmplDir + "/*.tmpl.html")
	if err != nil {
		util.LogError(err.Error())
	}
	if *s.Verbose {
		fmt.Println((*s).Tmpl.DefinedTemplates())
	}

	s.DB, err = sql.Open("sqlite", *s.DBName)
	if err != nil {
		return err
	}

	return nil
}

func ListenShutdown(sstate *ServerState) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch

	signal.Stop(ch)
	Shutdown(sstate, true, true)
}

func Shutdown(sstate *ServerState, shutdownHTTP bool, shutdownDB bool) {
	var err error

	if shutdownDB {
		err = sstate.DB.Close()
		if err != nil {
			log.Printf("Database closed with error: %s\n", err.Error())
		}
	}

	if shutdownHTTP {
		ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
		defer cancel()
		err = sstate.srv.Shutdown(ctx)
		if err != nil {
			log.Printf("HTTP Server closed with error: %s\n", err.Error())
		}
	}
}
