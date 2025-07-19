package main

import (
	"testing"
)

func TestCLIStructure(t *testing.T) {
	// Test that the CLI struct can be created
	cli := CLI{}
	if cli.SplitVideo.InputFile != "" {
		t.Error("Expected empty InputFile initially")
	}
}

func TestSplitVideoCmdDefaults(t *testing.T) {
	cmd := SplitVideoCmd{}
	if cmd.ChunkDuration != 0 { // Default is set by Kong, not struct
		t.Error("Expected ChunkDuration to be 0 initially")
	}
	if cmd.OutputDir != "" { // Default is set by Kong, not struct
		t.Error("Expected OutputDir to be empty initially")
	}
}