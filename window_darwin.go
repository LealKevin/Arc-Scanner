//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa -framework ApplicationServices
#import <Cocoa/Cocoa.h>
#import <ApplicationServices/ApplicationServices.h>

void setWindowLevelAboveScreensaver() {
    dispatch_async(dispatch_get_main_queue(), ^{
        NSArray *windows = [NSApp windows];
        NSLog(@"Number of windows: %lu", (unsigned long)[windows count]);
        if ([windows count] > 0) {
            NSWindow *window = [windows objectAtIndex:0];
            // Use CGWindowLevelForKey to get maximum window level
            // kCGMaximumWindowLevelKey = 14 (gives us level 2147483630)
            CGWindowLevel maxLevel = CGWindowLevelForKey(kCGMaximumWindowLevelKey);
            NSLog(@"Setting window level to: %d", maxLevel);
            [window setLevel:maxLevel];
            [window setCollectionBehavior:NSWindowCollectionBehaviorCanJoinAllSpaces |
                                          NSWindowCollectionBehaviorStationary |
                                          NSWindowCollectionBehaviorFullScreenAuxiliary];
            NSLog(@"Window level set successfully");
        } else {
            NSLog(@"No windows found!");
        }
    });
}

// Request both permissions upfront so user only needs one restart
void requestPermissions() {
    // Request Screen Recording permission (macOS 10.15+)
    if (@available(macOS 10.15, *)) {
        // CGRequestScreenCaptureAccess prompts user if not already granted
        bool hasScreenAccess = CGPreflightScreenCaptureAccess();
        if (!hasScreenAccess) {
            NSLog(@"Requesting Screen Recording permission...");
            CGRequestScreenCaptureAccess();
        }
    }

    // Check Accessibility permission - this will prompt if not trusted
    bool hasAccessibility = AXIsProcessTrusted();
    if (!hasAccessibility) {
        NSLog(@"Requesting Accessibility permission...");
        // Show the system prompt for accessibility
        NSDictionary *options = @{(__bridge id)kAXTrustedCheckOptionPrompt: @YES};
        AXIsProcessTrustedWithOptions((__bridge CFDictionaryRef)options);
    }
}

int checkPermissions() {
    bool hasScreen = true;
    if (@available(macOS 10.15, *)) {
        hasScreen = CGPreflightScreenCaptureAccess();
    }
    bool hasAccessibility = AXIsProcessTrusted();

    // Return: 0 = both granted, 1 = missing screen, 2 = missing accessibility, 3 = missing both
    if (hasScreen && hasAccessibility) return 0;
    if (!hasScreen && hasAccessibility) return 1;
    if (hasScreen && !hasAccessibility) return 2;
    return 3;
}
*/
import "C"

func setWindowAboveFullscreen() {
	println("Setting window level above fullscreen...")
	C.setWindowLevelAboveScreensaver()
}

// RequestPermissions prompts for both Screen Recording and Accessibility permissions
func requestPermissions() {
	C.requestPermissions()
}

// CheckPermissions returns permission status:
// 0 = both granted, 1 = missing screen, 2 = missing accessibility, 3 = missing both
func checkPermissions() int {
	return int(C.checkPermissions())
}
