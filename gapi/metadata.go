package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	// HTTP
	httpContentTypeHeader      = "grpcgateway-content-type"
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	xForwardedForHeader        = "x-forwarded-for"

	// gRPC
	grpcContentTypeHeader = "content-type"
	userAGentHeader       = "user-agent"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (s *Server) extractMetadata(ctx context.Context) Metadata {
	metadataFromContext, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return Metadata{}
	}

	var data Metadata
	switch {
	case len(metadataFromContext.Get(httpContentTypeHeader)) > 0:
		if userAgent := metadataFromContext.Get(grpcGatewayUserAgentHeader); len(userAgent) > 0 {
			data.UserAgent = userAgent[0]
		}

		if clientIPs := metadataFromContext.Get(xForwardedForHeader); len(clientIPs) > 0 {
			data.ClientIP = clientIPs[0]
		}

	case len(metadataFromContext.Get(grpcContentTypeHeader)) > 0:
		if userAgent := metadataFromContext.Get(userAGentHeader); len(userAgent) > 0 {
			data.UserAgent = userAgent[0]
		}

		if p, ok := peer.FromContext(ctx); ok {
			data.ClientIP = p.Addr.String()
		}
	}

	return data
}
