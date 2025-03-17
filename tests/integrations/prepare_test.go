package integrations

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"
)

type Response struct {
	Status  int                    `json:"status"`
	Success bool                   `json:"success"`
	Meta    map[string]interface{} `json:"meta"`
	Data    interface{}            `json:"data"`
}

var manager *RequestManager

var stopServer = make(chan struct{})

// we need smartPanic to be sure that our server sub-process always will be killed.
func smartPanic(p any) {
	close(stopServer)
	time.Sleep(time.Second * 2)
	panic(p)
}

func TestMain(m *testing.M) {
	if err := os.Chdir(".."); err != nil {
		fmt.Println(err)
	}

	build() // we do not use `go run` command because it left one child process (that process -> go run -> main)

	go run(stopServer)
	time.Sleep(time.Second * 3)

	if err := NewConfig(); err != nil {
		smartPanic(err)
	}

	URL = url.URL{Host: fmt.Sprintf("%v:%v", config.Host, config.Port), Scheme: "http"}

	manager = NewRequestManager()
	rc := m.Run()

	close(stopServer)

	time.Sleep(time.Second * 2)
	os.Exit(rc)
}

func build() {
	cmd := exec.Command("go", "build", "./cmd/main.go")
	res, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())

		return
	}

	fmt.Println(string(res))
}

func run(stop chan struct{}) {
	cmd := exec.Command("./main")

	err := cmd.Start()
	if err != nil {
		fmt.Println(err.Error())

		return
	}

	<-stop

	if err = cmd.Process.Kill(); err != nil {
		fmt.Println("!--", cmd.Process.Pid, "- not killed")
	} else {
		fmt.Println("!--", cmd.Process.Pid, "- killed successfully")
	}
}
