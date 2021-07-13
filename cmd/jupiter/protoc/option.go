package protoc

//Option ...
type Option struct {
	withGRPC      bool
	withServer    bool
	protoFilePath string
	outputDir     string
	prefix        string
}

var (
	option Option
)
