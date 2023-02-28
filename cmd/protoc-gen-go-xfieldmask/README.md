# protoc-gen-xfieldMask插件

## 说明
- 在protobuf中，`google.protobuf.FieldMask` 字段用于标识消息体中需要被更新（处理）的字段。
- `protoc-gen-xfieldMask`插件基于protobuf自定义扩展字段的方式，以插件->模板的形式在xxx_fm.pb.go文件中生成对应方法，方便客户端和服务端使用fieldMask进行字段选择性处理
- `protoc-gen-xfieldMask`可以根据需要选择性mask 请求和响应中的字段。请求和响应中的mask字段不会有冲突
- `protoc-gen-xfieldMask`可以深层指定message下的嵌套字段

## 使用
### helloworld.proto(proto/helloworld/v1/helloworld.proto)
```
syntax = "proto3";

package helloworld.v1;

import "google/protobuf/field_mask.proto";
import "fieldmask/v1/option.proto";

// The greeting service definition.
service GreeterService {
    //  Sends a goodbye greeting
    rpc SayGoodBye (SayGoodByeRequest) returns (SayGoodByeResponse);
}

// The request message containing the greetings
message SayGoodByeRequest{
    // name of the user
    string name = 1;
    // age of the user
    uint64 age = 2;
    // type filter mode
    Type type = 3;
    // update_mask FieldMask
    google.protobuf.FieldMask update_mask = 4[
        // Whether to mask the request field
        (fieldmask.v1.Option).in = true,
        // Whether to mask the request field
        (fieldmask.v1.Option).out = true
    ];
}

// The response message containing the greetings
message SayGoodByeResponse{
    // Data 响应数据
    message Data {
        // name of the user
        uint64 age = 1;
        // age of the user
        string name = 2;
        // other info of the user
        OtherHelloMessage other = 3;
    }
    // error...
    uint32 error = 1;
    // msg...
    string msg = 2;
    // data...
    Data data = 3;
}

// The response OtherHelloMessage containing the greetings
message OtherHelloMessage{
    // id...
    uint32 id = 1;
    // address...
    string address = 2;
}

// Type
enum Type {
    // TYPE_UNSPECIFIED ...
    TYPE_UNSPECIFIED = 0;
    // TYPE_Filter ... filter模式 表示mask的字段被保留
    TYPE_Filter = 1;
    // TYPE_Prune ...  prune模式 表示mask的字段被剔除
    TYPE_Prune = 2;
}
```

### 客户端调用(proto/helloworld/v1/fieldmask_test.go)
```
protoreq := &SayGoodByeRequest{
	Name: "foo",
	Type: Type_TYPE_Filter, // 表示采用过滤模式
}
// MaskInName:表示需要服务端处理name字段；
// MaskOutDataName：表示需要服务端返回data.name字段；
// MaskOutDataOther：表示需要服务端返回data.other下的所有字段
protoreq.MaskInName().MaskOutDataName().MaskOutDataOther()
```

### 服务端调用(proto/helloworld/v1/helloworld_impl.go)
```
func (s *FooServer) SayGoodBye(ctx context.Context, in *SayGoodByeRequest) (out *SayGoodByeResponse, err error) {
    // 初始化过滤/剔除器 
	var fm = new(SayGoodByeRequest_FieldMask)
	if in.Type == Type_TYPE_Filter {
		fm = in.FieldMaskFilter()
	} else {
		fm = in.FieldMaskPrune()
	}
	out = &SayGoodByeResponse{
		Error: 0,
		Msg:   "请求正常",
		Data: &SayGoodByeResponse_Data{
			Age:  1,
			Name: "",
			Other: &OtherHelloMessage{
				Id:      1,
				Address: "bar",
			},
		},
	}
    // 判断是否需要处理name字段
	if fm.MaskedInName() {
		out.Data.Name = in.GetName()
	}
    // 判断是否需要处理age字段
	if fm.MaskedInAge() {
		out.Data.Age = in.GetAge()
	}
    out1, _ := json.Marshal(out)
	fmt.Println("out1:", string(out1)) // out1: {"error":0,"msg":"请求正常","data":{"age":1,"name":"foo","other":{"id":1,"address":"bar"}}}
    // 过滤响应数据
	_ = fm.Mask(out)
	out2, _ := json.Marshal(out)
	fmt.Println("out2:", string(out2)) // out2: {"error":0,"msg":"请求正常","data":{"age":1,"name":"","other":null}}
	return
}
```


## 注意
- 推荐客户端服务端约定使用filter过滤模式
- message下的嵌套字段如果引入的是外部的消息体，则Mask字段不可再深层指定
- 目前在框架层面过滤了响应体中的error和msg字段，不管是否masked都会保留在响应体中