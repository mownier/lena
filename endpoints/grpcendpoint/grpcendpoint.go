package grpcendpoint

import (
	context "context"
	"fmt"
	"lena/auth"
	"lena/errors"
	"log"
	"net"

	"lena/config"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Server struct {
	authServer *auth.Server
	UnimplementedLenaServiceServer
}

func NewServer(authServer *auth.Server) *Server {
	return &Server{authServer: authServer}
}

func (s *Server) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	session, err := s.authServer.Register(ctx, in.Name, in.Password)
	if err != nil {
		domain := fmt.Sprintf("grpcendpoint.Server.Register: in = %v", safeRegisterRequest{in})
		return &RegisterResponse{}, errors.NewAppError(errors.ErrCodeRegistering, domain, err)
	}
	response := RegisterResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresOn:    timestamppb.New(session.AccesTokenExpiry),
	}
	return &response, nil
}

func (s *Server) SignIn(ctx context.Context, in *SignInRequest) (*SignInResponse, error) {
	session, err := s.authServer.SignIn(ctx, in.Name, in.Password)
	if err != nil {
		domain := fmt.Sprintf("grpcendpoint.Server.SignIn: in = %v", safeSignInRequest{in})
		return &SignInResponse{}, errors.NewAppError(errors.ErrCodeSigningIn, domain, err)
	}
	response := SignInResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresOn:    timestamppb.New(session.AccesTokenExpiry),
	}
	return &response, nil
}

func (s *Server) SignOut(ctx context.Context, emp *emptypb.Empty) (*emptypb.Empty, error) {
	domain := "grpcendpoint.Server.SignOut"
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
	}
	err = s.authServer.SignOut(ctx, accessToken)
	if err != nil {
		return &emptypb.Empty{}, errors.NewAppError(errors.ErrCodeSigningOut, domain, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Verify(ctx context.Context, emp *emptypb.Empty) (*emptypb.Empty, error) {
	domain := "grpcendpoint.Server.Verify"
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		return &emptypb.Empty{}, errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
	}
	err = s.authServer.Verify(ctx, accessToken)
	if err != nil {
		return &emptypb.Empty{}, errors.NewAppError(errors.ErrCodeVerifyingAccessToken, domain, err)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Refresh(ctx context.Context, in *RefreshRequest) (*RefreshResponse, error) {
	domain := fmt.Sprintf("grpcendpoint.Server.Refresh: in = %v", in)
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		return &RefreshResponse{}, errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
	}
	session, err := s.authServer.Refresh(ctx, accessToken, in.RefreshToken)
	if err != nil {
		return &RefreshResponse{}, errors.NewAppError(errors.ErrCodeRefreshingAccessToken, domain, err)
	}
	response := RefreshResponse{
		AccessToken:  session.AccessToken,
		RefreshToken: session.RefreshToken,
		ExpiresOn:    timestamppb.New(session.AccesTokenExpiry),
	}
	return &response, nil
}

func (s *Server) extractAccessToken(ctx context.Context) (string, error) {
	domain := "grpcendpoint.Server.extractAccessToken"
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.NewAppError(errors.ErrCodeMetadataNotOkay, domain, nil)
	}
	values := md.Get("Authorization")
	if len(values) == 0 {
		return "", errors.NewAppError(errors.ErrCodeAuthorizationNotSet, domain, nil)
	}
	accessToken := values[0]
	if accessToken == "" {
		return "", errors.NewAppError(errors.ErrCodeEmptyAuthorization, domain, nil)
	}
	return accessToken, nil
}

func ListenAndServe(config config.Config, authServer *auth.Server) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalln("failed to listen:", err)
	}
	grpcServer := grpc.NewServer()
	if config.Reflection {
		reflection.Register(grpcServer)
	}
	RegisterLenaServiceServer(grpcServer, NewServer(authServer))
	fmt.Printf("lena GRPC server listening on: tcp://%s:%d\n", config.LocalIP, config.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln("failed to server:", err)
	}
}

type safeRegisterRequest struct {
	*RegisterRequest
}

func (r safeRegisterRequest) String() string {
	return fmt.Sprintf("RegisterRequest{Name: %v, Password: ****}", r.Name)
}

type safeSignInRequest struct {
	*SignInRequest
}

func (r safeSignInRequest) String() string {
	return fmt.Sprintf("SignInRequest{Name: %v, Password: ****}", r.Name)
}
