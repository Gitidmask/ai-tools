@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

set APP_NAME=ai_tools
set ROOT_DIR=%~dp0..
set BUILD_DIR=%ROOT_DIR%\build
set SIDECAR_DIR=%ROOT_DIR%\sidecar
set FRONTEND_DIR=%ROOT_DIR%\frontend
set SIDECAR_OUT_DIR=%ROOT_DIR%\sidecars
set SIDECAR_EXE=%SIDECAR_OUT_DIR%\claude-sidecar-x86_64-pc-windows-msvc.exe

echo ========================================
echo   AI Toolbox - Build Script
echo ========================================
echo.

REM ---- Step 1: Build Sidecar (Bun --compile) ----
echo [1/3] Building sidecar...

if not exist "%SIDECAR_DIR%\package.json" (
    echo ERROR: sidecar/package.json not found
    goto :error
)

cd /d "%SIDECAR_DIR%"

where bun >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo WARNING: bun not found, skipping sidecar build
    goto :build_frontend
)

bun --version
echo Installing sidecar deps...
call bun install
if %ERRORLEVEL% NEQ 0 goto :error

echo Compiling sidecar...
if not exist "%SIDECAR_OUT_DIR%" mkdir "%SIDECAR_OUT_DIR%"
bun build --compile --target=browser --outfile="%SIDECAR_EXE%" ./src/index.ts
if %ERRORLEVEL% EQU 0 (
    echo sidecar OK: %SIDECAR_EXE%
) else (
    echo sidecar build FAILED, continuing...
)
echo.

REM ---- Step 2: Build Frontend (Vite) ----
:build_frontend
echo [2/3] Building frontend...

cd /d "%FRONTEND_DIR%"

if not exist "node_modules" (
    where npm >nul 2>nul
    if !ERRORLEVEL! NEQ 0 (
        echo WARNING: npm not found, using pre-built dist
        goto :build_app
    )
    echo Installing frontend deps...
    call npm install
    if !ERRORLEVEL! NEQ 0 goto :error
)

if not exist "dist\index.html" (
    echo Building frontend...
    call npm run build
    if !ERRORLEVEL! NEQ 0 goto :error
)
echo frontend OK
echo.

REM ---- Step 3: Build Go App ----
:build_app
echo [3/3] Building Go application...

cd /d "%ROOT_DIR%"

where go >nul 2>nul
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: go not found
    goto :error
)

if not exist "%BUILD_DIR%" mkdir "%BUILD_DIR%"

set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
set HTTPS_PROXY=http://127.0.0.1:7890
set HTTP_PROXY=http://127.0.0.1:7890
set GOPROXY=https://goproxy.cn,direct

go build -o "%BUILD_DIR%\%APP_NAME%.exe" -ldflags="-s -w" .
if %ERRORLEVEL% NEQ 0 (
    echo ERROR: go build failed
    goto :error
)
echo Go app OK: %BUILD_DIR%\%APP_NAME%.exe
echo.

REM ---- Copy sidecar to build output ----
if exist "%SIDECAR_EXE%" (
    copy /Y "%SIDECAR_EXE%" "%BUILD_DIR%\" >nul
    echo sidecar copied to build dir
)

echo ========================================
echo   BUILD COMPLETE
echo ========================================
echo.
echo Output: %BUILD_DIR%\
echo   %BUILD_DIR%\%APP_NAME%.exe
echo.
goto :end

:error
echo.
echo ====== BUILD FAILED ======
echo.

:end
pause
