package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwarded                 = "x-forwarded-for"
	userAgentHader             = "user-agent"
)

// grpcgateway-user-agent:[PostmanRuntime/7.44.0] x-forwarded-for:[::1]

type Metadata struct {
	ClientIp string
	UseAgnet string
}

func extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// fmt.Printf("Meta data: %+v", md)
		if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UseAgnet = userAgents[0]
		}
		if userAgents := md.Get(userAgentHader); len(userAgents) > 0 {
			mtdt.UseAgnet = userAgents[0]
		}
		if cliengIps := md.Get(xForwarded); len(cliengIps) > 0 {
			mtdt.ClientIp = cliengIps[0]
		}
		// log.Printf("md: %+v\n", md)
	}

    if p, ok := peer.FromContext(ctx); ok {
        mtdt.ClientIp = p.Addr.String()
        // log.Printf("peer: %+v\n", p)
    }
	return mtdt
}
