package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hemanth-1321/rest-api/internal/config"
	"github.com/hemanth-1321/rest-api/internal/http/student"
)


func main(){

 cfg:=config.MustLoad()
	fmt.Printf("Loaded config: %+v\n", cfg)

 router:=http.NewServeMux()
 router.HandleFunc("POST /api/student",student.New())


 server:=http.Server{
	Addr: cfg.HTTPServer.Addr,

	Handler: router,

 }
fmt.Println("server started")


done:=make(chan os.Signal,1)

signal.Notify(done,os.Interrupt,syscall.SIGINT,syscall.SIGTERM)
go func ()  {
	err:=server.ListenAndServe()

if err!=nil{
	log.Fatalf("failed to start server %v",err)
}
}()
<-done


slog.Info("Shutting down the server")

ctx,cancel:=context.WithTimeout(context.Background(),5*time.Second)


defer cancel()

err:=server.Shutdown(ctx)

if err!=nil{
	slog.Error("Failed to shutdown server",slog.String("error",err.Error()))
}


slog.Info("shut down successfully")
}