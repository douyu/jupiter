package banner

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xdebug"
)

func init() {
	flag.Register(&flag.BoolFlag{Name: "hide-banner", Usage: "--hide-banner=true", Default: false, Action: func(key string, fs *flag.FlagSet) {
		if !fs.Bool(key) {
			printBanner()
		}
	}})
}

//printBanner init
func printBanner() {
	if xdebug.IsTestingMode() {
		return
	}

	const banner = `
   (_)_   _ _ __ (_) |_ ___ _ __
   | | | | | '_ \| | __/ _ \ '__|
   | | |_| | |_) | | ||  __/ |
  _/ |\__,_| .__/|_|\__\___|_|
 |__/      |_|

 Welcome to jupiter, starting application ...
`
	fmt.Println(xcolor.Green(banner))
}
