package main

import (
	"bufio"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/edwardhey/flow/models"
	"github.com/edwardhey/flow/modules/nodes"
	"github.com/edwardhey/lib/utils"
	log "github.com/sirupsen/logrus"
)

// Readln returns a single line (without the ending \n)
// from the input buffered reader.
// An error is returned iff there is an error with the
// buffered reader.
func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func main() {
	pid := os.Getpid()
	log.Info(fmt.Sprintf("pid:%d", pid))
	utils.FilePutContents("/tmp/agent.pid", []byte(fmt.Sprintf("%d", pid)), os.O_RDWR|os.O_CREATE|os.O_TRUNC)

	//从环境变量读入节点ID
	sNodeID := os.Getenv("NODE_ID")
	if sNodeID == "" {
		panic("节点ID未配置")
	}

	iNodeID, err := strconv.ParseUint(sNodeID, 10, 64)
	if err != nil {
		panic(err)
	}

	node, err := models.GetNodeByID(models.UUID(iNodeID))
	if err != nil {
		panic(err)
	}
	if node.IsBegin == 0 {
		return
	}

	fifo := "/tmp/fifo"
	wg := sync.WaitGroup{}

	//从上一个节点的输出，写入到FIFO中，然后丢给当前节点处理
	{
		if _, err := os.Stat(fifo); os.IsNotExist(err) {
			err := syscall.Mkfifo(fifo, 0700)
			if err != nil {
				panic(err)
			}
		}
	}
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	fmt.Println("f")
	// 	f, err := os.OpenFile(fifo, os.O_RDONLY, os.ModePerm|os.ModeNamedPipe)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	// defer f.Close()
	// 	r := bufio.NewReader(f)
	// 	for {
	// 		line, err := r.ReadBytes('\n')
	// 		fmt.Println(line, err)
	// 		if err == nil {
	// 			fmt.Print("load string:" + string(line))
	// 		}
	// 	}
	// }()

	wg.Add(1)
	go func() {
		defer wg.Done()

		n := &nodes.OutputEngineFactory{}
		err := n.InitWithNode(node)
		if err != nil {
			log.Error(err)
			return
		}

		f, err := os.OpenFile(fifo, os.O_RDWR|os.O_SYNC, os.ModeNamedPipe)
		if err != nil {
			log.Error(err)
			return
			//return
		}

		//等待信号USR2信号
		{
			sig := make(chan os.Signal, 1)
			signal.Notify(sig, syscall.SIGUSR2)
			log.Info("Wait for USR2 signal")
			<-sig
			log.Info("Receive USR2 signal")
		}
		// w := bufio.NewWriter(f)
		for line := range n.Engine.GetLine() {
			// f.WriteString(line+"\n")
			_, err := f.WriteString(line + "\n")
			// fmt.Println(fifo, line, n, err)
			if err != nil {
				log.Error(err)
				break
			}
		}
	}()

	wg.Wait()
}
