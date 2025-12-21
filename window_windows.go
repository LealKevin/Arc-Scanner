//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	user32               = syscall.NewLazyDLL("user32.dll")
	procFindWindowW      = user32.NewProc("FindWindowW")
	procSetWindowPos     = user32.NewProc("SetWindowPos")
	procGetSystemMetrics = user32.NewProc("GetSystemMetrics")
	procGetDpiForWindow  = user32.NewProc("GetDpiForWindow")
)

const (
	HWND_TOPMOST   = ^uintptr(0) // -1
	SWP_NOSIZE     = 0x0001
	SWP_NOMOVE     = 0x0002
	SWP_NOACTIVATE = 0x0010
	SM_CXSCREEN    = 0
)

func setWindowAboveFullscreen() {
	hwnd := findWindowByTitle("arc-scanner")
	if hwnd == 0 {
		println("Warning: Could not find window handle")
		return
	}

	// Enable layered window for true transparency (per Wails Issue #1296)
	win.SetWindowLong(hwnd, win.GWL_EXSTYLE,
		win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)

	// Get screen width (physical pixels)
	screenWidth, _, _ := procGetSystemMetrics.Call(uintptr(SM_CXSCREEN))

	// Get DPI for this window (96 = 100%, 144 = 150%, etc.)
	dpi, _, _ := procGetDpiForWindow.Call(uintptr(hwnd))
	if dpi == 0 {
		dpi = 96 // Fallback to 100% if API fails (older Windows)
	}

	// Convert to logical coordinates: physical * 96 / dpi
	logicalWidth := int(screenWidth) * 96 / int(dpi)

	windowWidth := 400
	x := logicalWidth - windowWidth
	y := 0

	// Set window to always on top, let frontend handle position
	ret, _, err := procSetWindowPos.Call(
		uintptr(hwnd),
		HWND_TOPMOST,
		uintptr(x), uintptr(y),
		0, 0,
		SWP_NOSIZE|SWP_NOMOVE|SWP_NOACTIVATE,
	)

	if ret == 0 {
		fmt.Printf("Error setting window position: %v\n", err)
	}
}

func findWindowByTitle(title string) win.HWND {
	titlePtr, _ := syscall.UTF16PtrFromString(title)
	hwnd, _, _ := procFindWindowW.Call(
		0,
		uintptr(unsafe.Pointer(titlePtr)),
	)
	return win.HWND(hwnd)
}
