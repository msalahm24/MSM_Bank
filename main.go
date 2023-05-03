package main

import (
	"database/sql"
	"log"
	"net"

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
