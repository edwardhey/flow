package nodes

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/edwardhey/flow/models"
)

type INode interface {
	Init() error
	Run() error
	Complete() error
}

type IOutputEngine interface {
	GetLine() <-chan string
	SetConfig(*models.KindOutputConfig)
}

type OutputEngineFactory struct {
	Engine IOutputEngine
}

func (o *OutputEngineFactory) InitWithNode(n *models.Node) error {
	// var engine IOutputEngine
	switch n.OutputEngine {
	case models.KindOutputEngineFile:
		o.Engine = &OutputEngineFile{}
	case models.KindOutputEngineMysql:
		o.Engine = &OutputEngineMysql{}
	default:
		return fmt.Errorf("can't init output engine with node:%s", n.Name)
	}
	o.Engine.SetConfig(&n.OutputConfig)
	return nil
}

type OutputEngineBase struct {
	Config *models.KindOutputConfig
}

func (o *OutputEngineBase) SetConfig(c *models.KindOutputConfig) {
	o.Config = c
}

type OutputEngineFile struct {
	OutputEngineBase
}

func (of *OutputEngineFile) GetLine() <-chan string {
	c := make(chan string)
	go func() {
		defer close(c)
		f, err := os.OpenFile(of.Config.File.FileName, os.O_RDONLY, 0600)
		if err != nil {
			//TODO: 写日志
			return
		}
		defer f.Close()
		scanner := bufio.NewScanner(f)
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			c <- strings.TrimRight(scanner.Text(), "\n")
		}
	}()
	return c
}

type OutputEngineMysql struct {
	OutputEngineBase
}

func (om *OutputEngineMysql) GetLine() <-chan string {
	c := make(chan string)
	go func() {
	}()
	return c
}
