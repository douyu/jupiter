package xgoframe

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var errMicroDefault = status.Errorf(codes.Internal, createStatusErr(codeMS, "micro default"))

