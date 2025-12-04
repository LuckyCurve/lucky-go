package valuation

import (
	"testing"
)

func TestCAPECommand(t *testing.T) {
	t.Run("CommandStructure", func(t *testing.T) {
		cmd := NewCommand()

		if cmd.Use != "cape" {
			t.Errorf("expected command use 'cape', got '%s'", cmd.Use)
		}

		if cmd.Short != "查询标普500席勒CAPE估值" {
			t.Errorf("unexpected short description: '%s'", cmd.Short)
		}

		// 验证 push flag 存在
		pushFlag := cmd.Flags().Lookup("push")
		if pushFlag == nil {
			t.Fatal("expected push flag to exist")
			return
		}
		if pushFlag.Shorthand != "p" {
			t.Errorf("expected push shorthand 'p', got '%s'", pushFlag.Shorthand)
		}
	})
}
