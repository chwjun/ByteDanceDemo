package service

import (
	"testing"
)

func TestNewVSIInstance(t *testing.T) {
	vsi1 := NewVSIInstance()
	vsi2 := NewVSIInstance()
	if vsi1 != vsi2 {
		t.Error("单例测试出错")
	}
}
