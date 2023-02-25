package main

import (
	"Anti-bruteforce-service/internal/config"
	"Anti-bruteforce-service/internal/controller/grpcapi/authorizationpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/blacklistpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/bucketpb"
	"Anti-bruteforce-service/internal/controller/grpcapi/whitelistpb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
)

func main() {
	c, err := config.New()
	if err != nil {
		fmt.Println(err)
		return
	}

	insecureTr := grpc.WithTransportCredentials(insecure.NewCredentials())
	dial, err := grpc.Dial(c.Listen.BindIP+":"+c.Listen.Port, insecureTr)
	if err != nil {
		fmt.Println(err)
		return
	}
	clientBL := blacklistpb.NewBlackListServiceClient(dial)
	clientWL := whitelistpb.NewWhiteListServiceClient(dial)
	clientBucket := bucketpb.NewBucketServiceClient(dial)
	clientAuth := authorizationpb.NewAuthorizationClient(dial)
	getIpListInBlackList(clientBL)
	fmt.Println()
	getIpListInWhiteList(clientWL)
	fmt.Println()
	resetBucket(clientBucket)
	fmt.Println()
	tryAuth(clientAuth)

}

func tryAuth(client authorizationpb.AuthorizationClient) {
	response, err := client.TryAuthorization(context.Background(), &authorizationpb.AuthorizationRequest{Request: &authorizationpb.Request{
		Login:    "test",
		Password: "1234",
		Ip:       "192.1.5.1",
	}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.IsAllow)
}

func getIpListInBlackList(client blacklistpb.BlackListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &blacklistpb.GetIpListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		res, err := stream.Recv()

		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}
}

func getIpListInWhiteList(client whitelistpb.WhiteListServiceClient) {
	stream, err := client.GetIpList(context.Background(), &whitelistpb.GetIpListRequest{})
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	}
}

func resetBucket(client bucketpb.BucketServiceClient) {
	response, err := client.ResetBucket(context.Background(), &bucketpb.ResetBucketRequest{Request: &bucketpb.Request{
		Login:    "test",
		Password: "1234",
		Ip:       "192.1.5.1",
	}})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(response.ResetIp, response.ResetLogin)
}
