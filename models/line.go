package models

type Line struct {
	ID         UUID `gorm:"column:id"`
	PrevNodeID UUID `gorm:"column:prev_node_id"` //上一个节点ID
	NextNodeID UUID `gorm:"column:next_node_id"` //下一个节点ID
	Base
}
