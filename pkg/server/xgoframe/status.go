// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package xgoframe

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	rstatus "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	jsonpb "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// EmptyMessage ...
type EmptyMessage struct{}

// Reset ...
func (m *EmptyMessage) Reset() { *m = EmptyMessage{} }

// String ...
func (m *EmptyMessage) String() string { return "{}" }

// ProtoMessage ...
func (*EmptyMessage) ProtoMessage() {}

// GRPCProxyMessage ...
type GRPCProxyMessage struct {
	Error   int           `protobuf:"varint,1,opt,name=error" json:"error"`
	Message string        `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Data    proto.Message `protobuf:"bytes,3,opt,name=data" json:"data"`
}

// Reset ...
func (m *GRPCProxyMessage) Reset() { *m = GRPCProxyMessage{} }

// String ...
func (m *GRPCProxyMessage) String() string { return jsonpb.Format(m.Data) }

// ProtoMessage ...
func (*GRPCProxyMessage) ProtoMessage() {}

// MarshalJSONPB ...
func (m *GRPCProxyMessage) MarshalJSONPB() ([]byte, error) {

	ss, err := jsonpb.Marshal(m.Data)
	if err != nil {
		return []byte{}, err
	}

	msg := struct {
		Error   int             `json:"error"`
		Message string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
	}{
		Error:   m.Error,
		Message: m.Message,
		Data:    json.RawMessage([]byte(ss)),
	}

	return json.Marshal(msg)
}

type statusErr struct {
	s *rstatus.Status
}

// Proto ...
func (e *statusErr) Proto() *rstatus.Status {
	if e.s == nil {
		return nil
	}
	return proto.Clone(e.s).(*rstatus.Status)
}
func init() {
	s, _ := status.FromError(errMicroDefault)
	de, _ := statusFromString(s.Message())
	_ = de.Proto()
}
func statusFromString(s string) (*statusErr, bool) {
	i := strings.Index(s, ":")
	if i == -1 {
		return nil, false
	}
	u64, err := strconv.ParseInt(s[:i], 10, 32)
	if err != nil {
		return nil, false
	}

	return &statusErr{
		&rstatus.Status{
			Code:    int32(u64),
			Message: s[i:],
			Details: []*anypb.Any{},
		},
	}, true
}
func createStatusErr(code uint32, msg string) string {
	return fmt.Sprintf("%d:%s", code, msg)
}
