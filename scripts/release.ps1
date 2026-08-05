<#
    release.ps1 - generate sbom.json + release-manifest.json for the current checkout.

    Runs the Go test suites across the 8 Go modules, then generates the SBOM via
    scripts/sbom.go and the release manifest via scripts/release-manifest.go and
    prints the manifest path. This is a convenience wrapper for local development
    / release prep, not a full release pipeline (artifact signing is produced
    separately).

    Usage:
        .\scripts\release.ps1            # run tests + generate SBOM + manifest
        .\scripts\release.ps1 -SkipTests # only generate the artifacts
#>
[CmdletBinding()]
param(
    [switch]$SkipTests
)

$ErrorActionPreference = "Stop"
$root = Split-Path -Parent $PSScriptRoot
Set-Location $root

if (-not $SkipTests) {
    Write-Host "=== Running Go test suites (8-module matrix) ===" -ForegroundColor Cyan
    $modules = Get-ChildItem packages -Directory | Where-Object { Test-Path (Join-Path $_.FullName "go.mod") }
    foreach ($m in $modules) {
        Write-Host ("[TEST] " + $m.Name) -ForegroundColor Cyan
        Push-Location $m.FullName
        try {
            go test ./... -count=1
            if ($LASTEXITCODE -ne 0) { throw "go test failed in $($m.Name)" }
        }
        finally {
            Pop-Location
        }
    }
    Write-Host "[PASS] All module test suites green" -ForegroundColor Green
}

Write-Host "=== Generating SBOM ===" -ForegroundColor Cyan
go run ./scripts/sbom.go
if ($LASTEXITCODE -ne 0) { throw "SBOM generation failed" }

Write-Host "=== Generating release manifest ===" -ForegroundColor Cyan
go run ./scripts/release-manifest.go
if ($LASTEXITCODE -ne 0) { throw "release-manifest generation failed" }

$sbom = Join-Path $root "sbom.json"
Write-Host "[OK] SBOM written to $sbom" -ForegroundColor Green
$manifest = Join-Path $root "release-manifest.json"
Write-Host "[OK] Manifest written to $manifest" -ForegroundColor Green
Write-Host ""
Get-Content $manifest
