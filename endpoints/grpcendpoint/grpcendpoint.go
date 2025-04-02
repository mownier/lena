package grpcendpoint

import (
	context "context"
	"encoding/json"
	"fmt"
	"lena/auth"
	"lena/errors"
	"log"
	"net"

	"lena/config"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
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
	if err == nil {
		response := RegisterResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresOn:    timestamppb.New(session.AccesTokenExpiry),
		}
		return &response, nil
	}
	domain := fmt.Sprintf("grpcendpoint.Server.Register: in = %v", safeRegisterRequest{in})
	appError := errors.NewAppError(errors.ErrCodeRegistering, domain, err)
	var response errors.UserFriendlyResponse
	if other, contains := appError.ContainsCode(errors.ErrCodeUserAlreadyExists); contains {
		response = other.AsUserFriendlyResponse()
	} else {
		response = appError.AsUserFriendlyResponse()
	}
	jsonData, jsonErr := json.Marshal(response)
	var message string
	if jsonErr != nil {
		message = response.Message
	} else {
		message = string(jsonData)
	}
	return nil, status.Error(codes.Internal, message)
}

func (s *Server) SignIn(ctx context.Context, in *SignInRequest) (*SignInResponse, error) {
	session, err := s.authServer.SignIn(ctx, in.Name, in.Password)
	if err == nil {
		response := SignInResponse{
			AccessToken:  session.AccessToken,
			RefreshToken: session.RefreshToken,
			ExpiresOn:    timestamppb.New(session.AccesTokenExpiry),
		}
		return &response, nil
	}
	domain := fmt.Sprintf("grpcendpoint.Server.SignIn: in = %v", safeSignInRequest{in})
	appError := errors.NewAppError(errors.ErrCodeSigningIn, domain, err)
	var response errors.UserFriendlyResponse
	if other, contains := appError.ContainsCode(errors.ErrCodeUserDoesNotExist); contains {
		response = other.AsUserFriendlyResponse()
	} else if other, contains := appError.ContainsCode(errors.ErrCodeInvalidPassword); contains {
		response = other.AsUserFriendlyResponse()
	} else {
		response = appError.AsUserFriendlyResponse()
	}
	jsonData, jsonErr := json.Marshal(response)
	message := response.Message
	if jsonErr == nil {
		message = string(jsonData)
	}
	return nil, status.Error(codes.Internal, message)
}

func (s *Server) SignOut(ctx context.Context, emp *emptypb.Empty) (*emptypb.Empty, error) {
	domain := "grpcendpoint.Server.SignOut"
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
	}
	err = s.authServer.SignOut(ctx, accessToken)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeSigningOut, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Verify(ctx context.Context, emp *emptypb.Empty) (*emptypb.Empty, error) {
	domain := "grpcendpoint.Server.Verify"
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
	}
	err = s.authServer.Verify(ctx, accessToken)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeVerifyingAccessToken, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Refresh(ctx context.Context, in *RefreshRequest) (*RefreshResponse, error) {
	domain := fmt.Sprintf("grpcendpoint.Server.Refresh: in = %v", in)
	accessToken, err := s.extractAccessToken(ctx)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeGettingAccessToken, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
	}
	session, err := s.authServer.Refresh(ctx, accessToken, in.RefreshToken)
	if err != nil {
		appError := errors.NewAppError(errors.ErrCodeRefreshingAccessToken, domain, err)
		response := appError.AsUserFriendlyResponse()
		jsonData, jsonErr := json.Marshal(response)
		message := response.Message
		if jsonErr == nil {
			message = string(jsonData)
		}
		return nil, status.Error(codes.Internal, message)
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
