package main

import (
	"context"
	"flag"
	"log"
	"time"

	userv1 "github.com/chyiyaqing/gmicro-proto/golang/user/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type AuthClient struct {
	service  userv1.UserClient
	username string
	password string
}

func NewAuthClient(cc *grpc.ClientConn, username, password string) *AuthClient {
	service := userv1.NewUserClient(cc)
	return &AuthClient{service: service, username: username, password: password}
}

func (client *AuthClient) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &userv1.LoginRequest{
		Username: client.username,
		Password: client.password,
	}

	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}
	return res.GetAccessToken(), nil
}

type AuthInterceptor struct {
	authClient  *AuthClient
	authMethods map[string]bool
	accessToken string
}

func NewAuthInterceptor(authClient *AuthClient, authMethods map[string]bool, refreshDuration time.Duration) (*AuthInterceptor, error) {
	interceptor := &AuthInterceptor{
		authClient:  authClient,
		authMethods: authMethods,
	}

	// 刷新访问令牌
	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

func (interceptor *AuthInterceptor) refreshToken() error {
	accessToken, err := interceptor.authClient.Login()
	if err != nil {
		return err
	}

	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)
	return nil
}

func (interceptor *AuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration
		for {
			time.Sleep(wait)
			err := interceptor.refreshToken()
			if err != nil {
				wait = time.Second
			} else {
				wait = refreshDuration
			}
		}
	}()

	return nil
}

func (interceptor *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		log.Printf("--> unary interceptor: %s", method)
		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func (interceptor *AuthInterceptor) Stream() grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		log.Printf("--> stream interceptor: %s", method)

		if interceptor.authMethods[method] {
			return streamer(interceptor.attachToken(ctx), desc, cc, method, opts...)
		}

		return streamer(ctx, desc, cc, method, opts...)
	}
}

func (interceptor *AuthInterceptor) attachToken(ctx context.Context) context.Context {
	// 将输入上下文与授权密钥哈访问令牌值一起传入
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

// userClient is a client to call user service RPCs
type UserClient struct {
	service userv1.UserClient
}

// NewuserClient creates a new userClient
func NewUserClient(cc *grpc.ClientConn) *UserClient {
	service := userv1.NewUserClient(cc)
	return &UserClient{service: service}
}

func (UserClient *UserClient) GetUser(ctx context.Context, id int64) (*userv1.GetUserResponse, error) {
	req := &userv1.GetUserRequest{UserId: id}
	return UserClient.service.Get(ctx, req)
}

func authMethods() map[string]bool {
	const userServicePath = "/user.v1.User/"

	return map[string]bool{
		userServicePath + "Get": true,
	}
}

const (
	username        = "admin"
	password        = "admin123"
	refreshDuration = 5 * time.Second
)

func main() {
	serverAddress := flag.String("address", "", "the server address")
	flag.Parse()
	log.Printf("dial server %s", *serverAddress)

	cc1, err := grpc.NewClient(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	authClient := NewAuthClient(cc1, username, password)

	interceptor, err := NewAuthInterceptor(authClient, authMethods(), refreshDuration)
	if err != nil {
		log.Fatal("cannot create auth interceptor: ", err)
	}

	cc2, err := grpc.NewClient(
		*serverAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor.Unary()),
		grpc.WithStreamInterceptor(interceptor.Stream()),
	)
	if err != nil {
		log.Fatal("cannot dial server: ", err)
	}

	userClient := NewUserClient(cc2)
	for {
		userClient.service.Get(context.Background(), &userv1.GetUserRequest{UserId: 1})
		time.Sleep(10 * time.Second)
	}
}
