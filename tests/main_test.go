package main

import "testing"

func TestOneCheck(t *testing.T) {
  var b int = 1

  if b != 1 {
    t.Error("wrong", b)

  }
}
