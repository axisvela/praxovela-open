@echo off
setlocal
cd /d "%~dp0.."

echo ===== AxisRobo Agent - Test Suite =====
echo.

echo [TEST] Running all Go tests...
go test ./packages/archon-guard/... ./packages/janus-gateway/... ./packages/argus-trace/... ./packages/axon-core/... ./packages/vulcan-forge/... -count=1 -timeout 60s
if %ERRORLEVEL% neq 0 (
    echo [FAIL] Tests failed
    exit /b 1
)
echo [PASS] All tests passed

echo [BUILD] Building axond...
go build -o apps\desktop\src-tauri\binaries\axond.exe .\packages\axon-core\cmd\axond\
echo [OK] axond built

echo ===== All checks passed! =====
pause
