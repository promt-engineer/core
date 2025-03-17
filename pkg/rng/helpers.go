package rng

import (
	"crypto/tls"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func newClient(host, port string, isSecure bool) (RNGClient, error) {
	addr := host + ":" + port
	var (
		conn *grpc.ClientConn
		err  error
	)

	if isSecure {
		config := &tls.Config{
			InsecureSkipVerify: false,
		}

		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	} else {
		conn, err = grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	if err != nil {
		zap.S().Errorf("can not dial %v: %v", addr, err)

		return nil, err
	}

	return NewRNGClient(conn), nil
}

func sliceOfValues(value uint64, size int) []uint64 {
	buf := make([]uint64, size)
	for i := 0; i < size; i++ {
		buf[i] = value
	}

	return buf
}
