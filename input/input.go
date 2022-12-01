package ginput

import (
	"bufio"
	"fmt"
	"github.com/basicfu/gf/text/gstr"
	"os"
)

func WaitInput(tip string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(tip)
	fmt.Print("-> ")
	text, err := reader.ReadString('\n')
	if err != nil {
		panic(fmt.Errorf("error: %w \n", err))
	}
	return gstr.Replace(gstr.Replace(text, "\r\n", ""), "\n", "")
}
