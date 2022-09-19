package xfreecache

type Size int

const (
	Byte Size = 1
	KB        = 1024 * Byte
	MB        = 1024 * KB
	GB        = 1024 * MB
	TB        = 1024 * GB
)
