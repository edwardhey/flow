package models

import (
	"database/sql/driver"

	"mysoft.com/dmp_screen/lib/gorm_json"
)

type KindNodeIsMerge int8

const (
	KindNodeIsMergeFalse KindNodeIsMerge = 0
	KindNodeIsMergeTrue  KindNodeIsMerge = 1
)

type KindOutputEngine int8

const (
	KindOutputEngineFile  KindOutputEngine = 0 //文件
	KindOutputEngineMysql KindOutputEngine = 1 //mysql
)

type NodeOutputConfigFile struct {
	FileName string `json:"filename"`
}

type NodeOutputConfigMysql struct {
	Host     string `json:"h"`
	Port     string `json:"p"`
	Database string `json:"db"`
	Table    string `json:"t"`
}

type KindOutputConfig struct {
	File  *NodeOutputConfigFile  `json:"file,omitempty"`
	Mysql *NodeOutputConfigMysql `json:"mysql,omitempty"`
}

func (oc KindOutputConfig) Value() (driver.Value, error) {
	return gorm_json.Value(oc)
}

func (oc *KindOutputConfig) Scan(input interface{}) error {
	return gorm_json.Scan(input, oc)
}

type Node struct {
	ID              UUID             `gorm:"column:id"`
	Name            string           `gorm:"column:name;type:char(50)"`           //名称
	Description     string           `gorm:"column:description;type:char(255)"`   //描述
	FlowID          UUID             `gorm:"column:flow_id"`                      //流程ID
	IsMerge         KindNodeIsMerge  `gorm:"column:is_merge;type:tinyint(1)"`     //是否合并类型的节点
	IsBegin         uint8            `gorm:"column:is_begin;type:tinyint(1)"`     //是否开始节点
	IsEnd           uint8            `gorm:"column:is_end;type:tinyint(1)"`       //是否结束节点
	OutputEngine    KindOutputEngine `gorm:"column:output_engine"`                //输出引擎
	OutputConfig    KindOutputConfig `gorm:"column:output_config;type:char(255)"` //输出配置
	ImageURL        UUID             `gorm:"column:image_url"`                    //镜像地址
	ImageBootParams string           `gorm:"column:image_boot_params"`            //镜像启动参数
	Base
}

func GetBeginNodeByFlow(f *Flow) (*Node, error) {
	n := &Node{}
	err := db.Model(&Node{}).First(n, "is_begin=1 and flow_id=?", f.ID).Error
	if err != nil {
		return nil, err
	}
	return n, nil
}

func GetNodeByID(id UUID) (*Node, error) {
	var n Node
	err := db.Model(&Node{}).First(&n, "id=?", id).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}
