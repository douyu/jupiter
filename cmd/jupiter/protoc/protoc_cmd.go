package protoc

import "github.com/urfave/cli"

var Cmd = cli.Command{
	Name:            "protoc",
	Aliases:         []string{"p"},
	Usage:           "jupiter protoc tools",
	Action:          Run,
	SkipFlagParsing: false,
	UsageText:       ProtocHelpTemplate,
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:        "grpc,g",
			Usage:       "whether to generate GRPC code",
			Destination: &option.withGRPC,
		},
		&cli.BoolFlag{
			Name:        "server,s",
			Usage:       "whether to generate grpc server code",
			Destination: &option.withServer,
		},
		&cli.StringFlag{
			Name:        "file,f",
			Usage:       "Path of proto file",
			Required:    true,
			Destination: &option.protoFilePath,
		},
		&cli.StringFlag{
			Name:        "out,o",
			Usage:       "Path of code generation",
			Required:    true,
			Destination: &option.outputDir,
		},
		&cli.StringFlag{
			Name:        "prefix,p",
			Usage:       "prefix(current project name)",
			Required:    false,
			Destination: &option.prefix,
		},
	},
}
