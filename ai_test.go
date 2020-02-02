package main

import "testing"

func TestCountFramesTllBelow(t *testing.T) {

	framesBelow := countFramesUntilBelow(20, -5, 5)
	expectedFrames := 3
	if framesBelow != expectedFrames {
		t.Errorf("countFramesUntilBelow expected %v got %v", expectedFrames, framesBelow)
	}

}