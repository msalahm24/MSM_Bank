package gapi

import (
	"context"
	"net"

	"google.golang.org/grpc/metadata"
)

func getMetaData(ctx context.Context)(string,string){
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok{
		return "",""
	}
	userAgents := md.Get("user-agent")
	if len(userAgents)==0{
		return "",""
	}

	userAgent := userAgents[0]
	clientIP := ""
	xRealIPs,ok := md["x-real-ip"]

	if ok && len(xRealIPs)>0{
		clientIP = xRealIPs[0]
	}

	xForwardedFor := md["x-forwarded-for"]

	if ok && len(xForwardedFor) > 0{
		clientIP = xForwardedFor[0]
	} 
	clientAddr := net.ParseIP(clientIP)

	if clientAddr.To4() == nil{
		clientIP = clientAddr.To16().To4().String()
	}
	return userAgent,clientIP
}