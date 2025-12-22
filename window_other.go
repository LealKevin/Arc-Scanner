//go:build !darwin && !windows

package main

func setWindowAboveFullscreen() {}

func requestPermissions() {}

func checkPermissions() int { return 0 }
