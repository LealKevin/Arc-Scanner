//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/lxn/win"
)

var (
	user32                         = syscall.NewLazyDLL("user32.dll")
	dwmapi                         = syscall.NewLazyDLL("dwmapi.dll")
	procFindWindowW                = user32.NewProc("FindWindowW")
	procSetWindowPos               = user32.NewProc("SetWindowPos")
	procGetSystemMetrics           = user32.NewProc("GetSystemMetrics")
	procGetDpiForWindow            = user32.NewProc("GetDpiForWindow")
	procDwmExtendFrameIntoClientArea = dwmapi.NewProc("DwmExtendFrameIntoClientArea")
)

const (
	HWND_TOPMOST   = ^uintptr(0) // -1
	SWP_NOSIZE     = 0x0001
	SWP_NOMOVE     = 0x0002
	SWP_NOACTIVATE = 0x0010
	SM_CXSCREEN    = 0
)

// MARGINS structure for DwmExtendFrameIntoClientArea
type MARGINS struct {
	cxLeftWidth    int32
	cxRightWidth   int32
	cyTopHeight    int32
	cyBottomHeight int32
}

func setWindowAboveFullscreen() {
	hwnd := findWindowByTitle("arc-scanner")
	if hwnd == 0 {
		println("Warning: Could not find window handle")
		return
	}

	// Use DWM composition for transparency (fixes Windows 11 black background issue)
	// Set margins to -1 to extend the frame into the entire client area
	margins := MARGINS{-1, -1, -1, -1}
	dwmRet, _, _ := procDwmExtendFrameIntoClientArea.Call(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&margins)),
	)
	if dwmRet != 0 {
		// If DWM fails (older Windows or DWM disabled), fall back to layered window
		fmt.Printf("DwmExtendFrameIntoClientArea failed with HRESULT 0x%X, falling back to WS_EX_LAYERED\n", dwmRet)
		win.SetWindowLong(hwnd, win.GWL_EXSTYLE,
			win.GetWindowLong(hwnd, win.GWL_EXSTYLE)|win.WS_EX_LAYERED)
	}

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

// Windows doesn't require special permissions like macOS
func requestPermissions() {}

func checkPermissions() int { return 0 }
