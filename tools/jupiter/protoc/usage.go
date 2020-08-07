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

package protoc

//ProtocHelpTemplate ...
const ProtocHelpTemplate = `
jupiter [commands|flags]

The commands & flags are:
  protoc        jupiter protoc tools
  -g,--grpc     whether to generate GRPC code
  -s,--server   whether to generate grpc server code
  -c,--client   generate grpc server code
  -f,--file     path of proto file
  -o,--out      path of code generation
  -p,--prefix   prefix(current project name)
Examples:
   # Generate GRPC code from the Proto file 
   # -f: Proto file address -o: Code generation path -g: Whether to generate GRPC code
   jupiter protoc -f ./pb/hello/hello.proto -o ./pb/hello -g
   # According to the proto file, generate the server implementation
   # -f: Proto file address -o: Code generation path -p:prefix(Current project name) -g: Whether to generate Server code
   jupiter protoc -f ./pb/hello/hello.proto -o ./internal/app/grpc -p jupiter-demo -s
  
`
