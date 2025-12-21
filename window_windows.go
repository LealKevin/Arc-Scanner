//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	user32           = syscall.NewLazyDLL("user32.dll")
	procFindWindowW  = user32.NewProc("FindWindowW")
	procSetWindowPos = user32.NewProc("SetWindowPos")
)

const (
	HWND_TOPMOST   = ^uintptr(0) // -1
	SWP_NOMOVE     = 0x0002
	SWP_NOSIZE     = 0x0001
	SWP_NOACTIVATE = 0x0010
)

func setWindowAboveFullscreen() {
	println("Setting window level to always on top (Windows)...")

	// Find the window by enumerating all windows and looking for our title
	// Since we don't have direct access to the HWND from Wails runtime,
	// we'll use the window title to find it
	hwnd := findWindowByTitle("arc-scanner")
	if hwnd == 0 {
		println("Warning: Could not find window handle")
		return
	}

	// Set window to be always on top
	ret, _, err := procSetWindowPos.Call(
		uintptr(hwnd),
		HWND_TOPMOST,
		0, 0, 0, 0,
		SWP_NOMOVE|SWP_NOSIZE|SWP_NOACTIVATE,
	)

	if ret == 0 {
		fmt.Printf("Error setting window to topmost: %v\n", err)
		return
	}

	println("Window set to always on top successfully")
}

func findWindowByTitle(title string) win.HWND {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(titlePtr)),
	)
	return win.HWND(hwnd)
}
