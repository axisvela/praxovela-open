@echo off
setlocal
cd /d "%~dp0.."

echo ==========================================
echo   AxisRobo Agent - Development Launcher
echo ==========================================
echo.

where go >nul 2>&1
if %ERRORLEVEL% neq 0 (
    echo [ERROR] Go not found. Install from https://go.dev/dl/
    exit /b 1
)
echo [OK] Go found

echo [BUILD] Building axond.exe...
go build -o apps\desktop\src-tauri\binaries\axond.exe .\packages\axon-core\cmd\axond\
if %ERRORLEVEL% neq 0 (
    echo [ERROR] axond build failed
    exit /b 1
)
echo [OK] axond.exe built

echo [START] Starting AXON Core on port 8420...
start "AXON Core" /B apps\desktop\src-tauri\binaries\axond.exe
echo [OK] AXON Core started

echo [WAIT] Waiting for AXON Core...
:wait
timeout /t 2 /nobreak >nul
curl -s http://localhost:8420/health >nul 2>&1
if %ERRORLEVEL% neq 0 goto wait
echo [OK] AXON Core healthy

echo ==========================================
echo   AXON Core running on http://localhost:8420
echo   Open apps\desktop\dist\index.html in browser
echo ==========================================
pause
