@echo off
setlocal
REM Change to project root (one level up from scripts/)
cd /d "%~dp0.."

echo ===== AxisRobo Agent - Full Build =====

echo [1/4] Running Go tests...
go test ./packages/archon-guard/... ./packages/janus-gateway/... ./packages/argus-trace/... ./packages/axon-core/... ./packages/vulcan-forge/... -count=1 -timeout 60s
if %ERRORLEVEL% neq 0 (
    echo [FAIL] Tests failed
    exit /b 1
)
echo [OK] All tests passed

echo [2/4] Building axond.exe...
go build -o apps\desktop\src-tauri\binaries\axond.exe .\packages\axon-core\cmd\axond\
if %ERRORLEVEL% neq 0 (
    echo [ERROR] axond build failed
    exit /b 1
)
echo [OK] axond.exe built

echo [3/4] Copying sidecar...
copy /Y apps\desktop\src-tauri\binaries\axond.exe apps\desktop\src-tauri\binaries\axond-x86_64-pc-windows-gnu.exe >nul
echo [OK] Sidecar copied

echo [4/4] Building Tauri release...
if exist apps\desktop\src-tauri\dist rmdir /S /Q apps\desktop\src-tauri\dist
xcopy /E /I /Y apps\desktop\dist apps\desktop\src-tauri\dist >nul
set PATH=D:\dev\msys64\mingw64\bin;D:\data\Rust\.cargo\bin;%PATH%
cargo build --manifest-path apps\desktop\src-tauri\Cargo.toml --release --target-dir "D:\data\Rust\global_target"
if %ERRORLEVEL% neq 0 exit /b 1

REM Copy binaries
mkdir D:\data\Rust\global_target\release\binaries 2>nul
copy /Y apps\desktop\src-tauri\binaries\axond.exe D:\data\Rust\global_target\release\binaries\axond.exe >nul

echo ===== Build complete! =====
echo Release: D:\data\Rust\global_target\release\axisrobo-desktop.exe
echo Run: D:\data\Rust\global_target\release\axisrobo-desktop.exe
exit /b 0
