@echo off
setlocal enabledelayedexpansion

echo Bundling Tesseract for Windows...
echo.

set BUILD_DIR=build\windows

REM Clean and create directories
if exist "%BUILD_DIR%\bin" rmdir /s /q "%BUILD_DIR%\bin"
if exist "%BUILD_DIR%\lib" rmdir /s /q "%BUILD_DIR%\lib"
if exist "%BUILD_DIR%\tessdata" rmdir /s /q "%BUILD_DIR%\tessdata"

mkdir "%BUILD_DIR%\bin" 2>nul
mkdir "%BUILD_DIR%\lib" 2>nul
mkdir "%BUILD_DIR%\tessdata" 2>nul

REM Try to find Tesseract installation
set TESS_SOURCE=
if exist "C:\Program Files\Tesseract-OCR\tesseract.exe" (
    set TESS_SOURCE=C:\Program Files\Tesseract-OCR
) else if exist "C:\Program Files (x86)\Tesseract-OCR\tesseract.exe" (
    set TESS_SOURCE=C:\Program Files (x86)\Tesseract-OCR
) else (
    echo Error: Tesseract not found!
    echo Please install Tesseract from: https://github.com/UB-Mannheim/tesseract/wiki
    echo Or set TESSERACT_PATH environment variable to your Tesseract installation directory
    if defined TESSERACT_PATH (
        set TESS_SOURCE=%TESSERACT_PATH%
    ) else (
        exit /b 1
    )
)

echo Found Tesseract at: %TESS_SOURCE%
echo.

REM Copy Tesseract binary
echo Copying Tesseract binary...
copy "%TESS_SOURCE%\tesseract.exe" "%BUILD_DIR%\bin\" >nul
if errorlevel 1 (
    echo Error copying tesseract.exe
    exit /b 1
)
echo   tesseract.exe copied

REM Copy all DLL dependencies
echo Copying DLL dependencies...
for %%F in ("%TESS_SOURCE%\*.dll") do (
    copy "%%F" "%BUILD_DIR%\lib\" >nul
    echo   %%~nxF copied
)

REM Copy training data
echo Copying training data...
if exist "%TESS_SOURCE%\tessdata\eng.traineddata" (
    copy "%TESS_SOURCE%\tessdata\eng.traineddata" "%BUILD_DIR%\tessdata\" >nul
    echo   eng.traineddata copied
) else (
    echo Warning: eng.traineddata not found in Tesseract installation
    echo Downloading from GitHub...
    powershell -Command "Invoke-WebRequest -Uri 'https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata' -OutFile '%BUILD_DIR%\tessdata\eng.traineddata'"
    if errorlevel 1 (
        echo Error downloading training data
        exit /b 1
    )
    echo   eng.traineddata downloaded
)

echo.
echo Tesseract bundling complete!
echo Binary: %BUILD_DIR%\bin\tesseract.exe
echo Libraries: %BUILD_DIR%\lib\
echo Training data: %BUILD_DIR%\tessdata\eng.traineddata
