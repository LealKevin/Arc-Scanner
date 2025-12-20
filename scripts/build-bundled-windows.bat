@echo off
setlocal enabledelayedexpansion

echo ===============================================
echo Building self-contained Arc Scanner for Windows
echo ===============================================
echo.

REM Step 1: Bundle Tesseract
echo Step 1/3: Bundling Tesseract...
echo -----------------------------------------------
call scripts\bundle-tesseract-windows.bat
if errorlevel 1 (
    echo.
    echo Error: Tesseract bundling failed!
    exit /b 1
)
echo.

REM Step 2: Build with Wails
echo Step 2/3: Building app with Wails...
echo -----------------------------------------------
wails build -platform windows/amd64
if errorlevel 1 (
    echo.
    echo Error: Wails build failed!
    exit /b 1
)
echo.

REM Step 3: Copy bundled resources
echo Step 3/3: Copying bundled resources to app...
echo -----------------------------------------------
call scripts\post-build-windows.bat
if errorlevel 1 (
    echo.
    echo Error: Post-build script failed!
    exit /b 1
)
echo.

echo ===============================================
echo Build complete!
echo ===============================================
echo   App location: build\bin\arc-scanner.exe
echo   App is self-contained with bundled Tesseract
echo   Total package size: ~15-20MB
echo.
echo To run: build\bin\arc-scanner.exe
echo To distribute: Copy the entire 'build\bin' folder
