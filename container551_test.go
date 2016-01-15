package container551_test

import (
	"github.com/go51/container551"
	"testing"
)

func TestNew(t *testing.T) {
	c1 := container551.New()
	c2 := container551.New()

	if c1 == nil {
		t.Error("インスタンスの生成に失敗しました。")
	}

	if c2 == nil {
		t.Error("インスタンスの生成に失敗しました。")
	}

	if &c1 == &c2 {
		t.Error("インスタンスの生成に失敗しました。")
	}
}

func BenchmarkNew(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = container551.New()
	}
}
