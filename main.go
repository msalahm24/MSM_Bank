package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/mahmoud24598salah/MSM_Bank/api"
	db "github.com/mahmoud24598salah/MSM_Bank/db/sqlc"
	"github.com/mahmoud24598salah/MSM_Bank/gapi"
	"github.com/mahmoud24598salah/MSM_Bank/pb"
	"github.com/mahmoud24598salah/MSM_Bank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Can not read the config file", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can not connect to database", err)
	}

	store := db.NewStore(conn)
	go runGatewayServer(config,store)
	runGrpcServer(config,store)

}

func runGrpcServer(config util.Config, store db.Store){
	server,err := gapi.Newserver(config,store)
	if err != nil{
		log.Fatal("Cannot create gRPC server", err)
	} 

	grpcServer := grpc.NewServer()
	pb.RegisterMSMBankServer(grpcServer,server)
	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp",config.GRPCSreverAddress)

	if err != nil{
		log.Fatal("can not create listener: ",err)
	}
	log.Printf("start gRPC server at %s",listener.Addr().String())
	err = grpcServer.Serve(listener)
	if  err != nil{
		log.Fatal("can not start gRPC server", err)
	}
}

func runGatewayServer(config util.Config, store db.Store){
	server,err := gapi.Newserver(config,store)
	if err != nil{
		log.Fatal("Cannot create gRPC server", err)
	} 

	grpcmux := runtime.NewServeMux()
	ctx,cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterMSMBankHandlerServer(ctx,grpcmux,server)
	if err != nil{
		log.Fatal("can not register handler server ")
	}

	mux :=http.NewServeMux()
	mux.Handle("/",grpcmux)
	listener, err := net.Listen("tcp",config.HTTPSreverAddress)

	if err != nil{
		log.Fatal("can not create listener: ",err)
	}
	log.Printf("start http gateway server at %s",listener.Addr().String())
	err = http.Serve(listener,mux)
	if  err != nil{
		log.Fatal("can not start http gateway server", err)
	}
}

func runGinServer(config util.Config, store db.Store){
	server, err := api.Newserver(config, store)
	if err != nil {
		log.Fatal("Cannot create server", err)
	}

	err = server.Start(config.GRPCSreverAddress )
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
