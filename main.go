package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
  "flag"
  "strings"


	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
	"github.com/robfig/cron/v3"
	"github.com/lm317379829/livetv/global"
	"github.com/lm317379829/livetv/route"
	"github.com/lm317379829/livetv/service"
)

func main() {
  Dir := flag.String("DIR","/etc","路径")
  Url := flag.String("URL","0.0.0.0","IP")
  Port := flag.String("PORT","9000","端口号")
  flag.Parse()
  var LIVETV_LISTEN string = *Url + ":" + *Port
  var LIVETV_DATADIR string = *Dir
    if strings.HasSuffix(LIVETV_DATADIR, "/") {
    LIVETV_DATADIR = LIVETV_DATADIR[:len(LIVETV_DATADIR)-1]
  }
	rand.Seed(time.Now().UnixNano())
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Server listen", LIVETV_LISTEN)
	log.Println("Server datadir", LIVETV_DATADIR)
	logFile, err := os.OpenFile(LIVETV_DATADIR+"/livetv.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Panicln(err)
	}
	log.SetOutput(io.MultiWriter(os.Stderr, logFile))
	err = global.InitDB(LIVETV_DATADIR + "/livetv.db")
	if err != nil {
		log.Panicf("init: %s\n", err)
	}
	log.Println("LiveTV starting...")
	go service.LoadChannelCache()
	c := cron.New()
	_, err = c.AddFunc("0 */4 * * *", service.UpdateURLCache)
	if err != nil {
		log.Panicf("preloadCron: %s\n", err)
	}
	c.Start()
	sessionSecert, err := service.GetConfig("password")
	if err != nil {
		sessionSecert = "sessionSecert"
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	store := cookie.NewStore([]byte(sessionSecert))
	router.Use(sessions.Sessions("mysession", store))
	router.Static("/assert", "./assert")
	route.Register(router)
	srv := &http.Server{
		Addr:    LIVETV_LISTEN,
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Panicf("listen: %s\n", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shuting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Panicf("Server forced to shutdown: %s\n", err)
	}
	log.Println("Server exiting")
	logFile.Close()
}
