// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

const TableNameCasbinRule = "casbin_rule"

// CasbinRule mapped from table <casbin_rule>
type CasbinRule struct {
	PType string `gorm:"column:p_type;type:varchar(100)" json:"p_type"`
	V0    string `gorm:"column:v0;type:varchar(100)" json:"v0"`
	V1    string `gorm:"column:v1;type:varchar(100)" json:"v1"`
	V2    string `gorm:"column:v2;type:varchar(100)" json:"v2"`
	V3    string `gorm:"column:v3;type:varchar(100)" json:"v3"`
	V4    string `gorm:"column:v4;type:varchar(100)" json:"v4"`
	V5    string `gorm:"column:v5;type:varchar(100)" json:"v5"`
}

// TableName CasbinRule's table name
func (*CasbinRule) TableName() string {
	return TableNameCasbinRule
}
