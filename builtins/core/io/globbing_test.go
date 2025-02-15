package io

import (
	"testing"

	"github.com/lmorg/murex/test"
)

func TestLsG(t *testing.T) {
	tests := []test.MurexTest{
		// g
		{
			Block:  "g: README*",
			Stdout: "README.md",
		},
		{
			Block:   "g: README",
			Stderr:  "Error",
			ExitNum: 1,
		},
		// !g
		{
			Block:   "!g: *",
			Stderr:  "Error",
			ExitNum: 1,
		},
		{
			Block:  "!g: README",
			Stdout: "README.md",
		},
	}
	test.RunMurexTestsRx(tests, t)
}

func TestLsRx(t *testing.T) {
	tests := []test.MurexTest{
		// rx
		{
			Block:  "rx: R*ME",
			Stdout: "README.md",
		},
		{
			Block:   "rx: README$",
			Stderr:  "Error",
			ExitNum: 1,
		},
		// !rx
		{
			Block:   "!rx: .*",
			Stderr:  "Error",
			ExitNum: 1,
		},
		{
			Block:  `!rx: (go|yaml)`,
			Stdout: "README.md",
		},
	}
	test.RunMurexTestsRx(tests, t)
}

func TestLsF(t *testing.T) {
	tests := []test.MurexTest{
		// f
		{
			Block:  "f: +f",
			Stdout: "README.md",
		},
	}
	test.RunMurexTestsRx(tests, t)
}
