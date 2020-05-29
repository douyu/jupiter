// +build !linux

package rotate

import (
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}
