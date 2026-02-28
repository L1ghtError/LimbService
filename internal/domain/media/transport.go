package media

import (
	"bufio"
)

type Transport interface {
	ProcessImage(task ImageTask, w *bufio.Writer) error

	GetWorkersInfo() (AppInfoTask, error)
}
