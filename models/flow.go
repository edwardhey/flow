package models

import (
	"time"
)

type KindFlowStatus int8

const (
	KindFlowStatusDisabled KindFlowStatus = 0 //停用
	KindFlowStatusEndabled KindFlowStatus = 1 //启用
)

type KindFlowExecStatus int

const (
	KindFlowExecStatusPending KindFlowExecStatus = 0  //等候启动
	KindFlowExecStatusRunning KindFlowExecStatus = 10 //正在运行中
)

type Flow struct {
	ID           UUID               `gorm:"column:id"`
	Name         string             `gorm:"column:name;type:char(50)"`                 //名称
	Description  string             `gorm:"column:description;type:char(255)"`         //描述
	Schedule     string             `gorm:"column:schedule;type:char(10)"`             //调度计划
	Status       KindFlowStatus     `gorm:"column:status;default:0;type:tinyint(1)"`   //状态
	ExecStatus   KindFlowExecStatus `gorm:"column:exec_status;default:0;type:int(11)"` //状态
	LastExecTime *time.Time         `gorm:"column:last_exec_time"`                     //上次执行时间
	NextExecTime *time.Time         `gorm:"column:next_exec_time"`                     //下次执行时间
	Base
}
