@echo off
setlocal enabledelayedexpansion

set BUILD_DIR=build\windows
set BIN_DIR=build\bin

if not exist "%BIN_DIR%\arc-scanner.exe" (
    echo Error: arc-scanner.exe not found at %BIN_DIR%\arc-scanner.exe
    echo Run 'wails build -platform windows/amd64' first
    exit /b 1
)

echo Copying bundled resources alongside executable...
echo.

REM Create directory structure next to the .exe
echo Creating directory structure...
if not exist "%BIN_DIR%\windows" mkdir "%BIN_DIR%\windows"
if not exist "%BIN_DIR%\windows\bin" mkdir "%BIN_DIR%\windows\bin"
if not exist "%BIN_DIR%\windows\lib" mkdir "%BIN_DIR%\windows\lib"
if not exist "%BIN_DIR%\windows\tessdata" mkdir "%BIN_DIR%\windows\tessdata"

REM Copy bundled resources
echo Copying bin...
copy "%BUILD_DIR%\bin\*" "%BIN_DIR%\windows\bin\" >nul
echo   Copied binaries

echo Copying lib...
copy "%BUILD_DIR%\lib\*" "%BIN_DIR%\windows\lib\" >nul
echo   Copied libraries

echo Copying tessdata...
copy "%BUILD_DIR%\tessdata\*" "%BIN_DIR%\windows\tessdata\" >nul
echo   Copied training data

echo.
echo Done! Bundled resources copied to %BIN_DIR%\windows\
echo.
echo To distribute:
echo   - Copy the entire 'build\bin' folder
echo   - Users can run arc-scanner.exe directly (no installation needed)
echo   - The app will automatically use the bundled Tesseract
