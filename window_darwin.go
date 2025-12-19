//go:build darwin

package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

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
*/
import "C"

func setWindowAboveFullscreen() {
	println("Setting window level above fullscreen...")
	C.setWindowLevelAboveScreensaver()
}
