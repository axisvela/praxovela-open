<#
    verify-release.ps1 - verify a staged PRAXOVELA open-distribution release set.

    The open repository is distribution-only: release assets are produced by the
    core repository and merely confirmed here. This script verifies a staged
    release directory against the core-produced release manifest, SBOM, and
    SHA256SUMS. It NEVER generates release evidence - it only reads.

    Checks performed:
      1. canonical artifact set for the version (exact names, no extras/missing)
      2. SHA256SUMS agreement: exactly one matching entry per artifact, no
         duplicates, no unexpected entries, and no self-entry for SHA256SUMS
      3. manifest / SBOM / conformance-report version agreement
      4. manifest provenance: source_commit and distribution_commit both present
         and not equal

    Usage:
        .\scripts\verify-release.ps1 -ReleaseDir dist\v2.0.0
        .\scripts\verify-release.ps1 -ReleaseDir dist\v2.0.0 -Version 2.0.0

    Exit code is 0 when every check passes and 1 (with a [FAIL] list) otherwise.
#>
[CmdletBinding()]
param(
    [string]$ReleaseDir = "dist",
    [string]$Version
)

$ErrorActionPreference = "Stop"

# The script lives in <repo>/scripts; the repo root is its parent.
$repoRoot = Split-Path -Parent $PSScriptRoot

function Add-Failure {
    param(
        [System.Collections.Generic.List[string]]$List,
        [string]$Message
    )
    $List.Add($Message) | Out-Null
}

# Resolve a relative -ReleaseDir against the repo root so the default is stable
# regardless of the caller's working directory.
if (-not [System.IO.Path]::IsPathRooted($ReleaseDir)) {
    $ReleaseDir = Join-Path $repoRoot $ReleaseDir
}

# Resolve the version from version.json when not supplied explicitly.
if (-not $Version) {
    $versionFile = Join-Path $repoRoot "version.json"
    if (-not (Test-Path -LiteralPath $versionFile)) {
        Write-Host "[FAIL] version.json not found at $versionFile; pass -Version explicitly"
        exit 1
    }
    try {
        $Version = (Get-Content -LiteralPath $versionFile -Raw | ConvertFrom-Json).version
    }
    catch {
        Write-Host "[FAIL] cannot parse version.json: $($_.Exception.Message)"
        exit 1
    }
}

if ($Version -notmatch '^[0-9]+\.[0-9]+\.[0-9]+$') {
    Write-Host "[FAIL] version '$Version' is not clean major.minor.patch"
    exit 1
}

if (-not (Test-Path -LiteralPath $ReleaseDir -PathType Container)) {
    Write-Host "[FAIL] release directory not found: $ReleaseDir"
    exit 1
}

$sumsFile = "SHA256SUMS"
$manifestName = "praxovela-${Version}-release-manifest.json"
$sbomName = "praxovela-${Version}-sbom.cdx.json"
$conformanceName = "praxovela-${Version}-syndovela-conformance.txt"

$expected = @(
    "axond-windows-amd64.exe",
    "axond-linux-amd64",
    "PRAXOVELA_${Version}_x64-setup.exe",
    "PRAXOVELA_${Version}_x64_en-US.msi",
    $sbomName,
    $manifestName,
    $conformanceName,
    $sumsFile
)

$failures = New-Object 'System.Collections.Generic.List[string]'

# --- Check 1: canonical artifact set (exact names, no extras) -----------------
$expectedSet = @{}
foreach ($name in $expected) { $expectedSet[$name] = $true }

$have = @{}
foreach ($f in Get-ChildItem -LiteralPath $ReleaseDir -File) { $have[$f.Name] = $true }

foreach ($name in $expected) {
    if (-not $have.ContainsKey($name)) {
        Add-Failure $failures "artifact '$name' missing from $ReleaseDir"
    }
}
foreach ($name in ($have.Keys | Sort-Object)) {
    if (-not $expectedSet.ContainsKey($name)) {
        Add-Failure $failures "unexpected artifact '$name' (stale embedded version?)"
    }
}

# --- Check 2: SHA256SUMS agreement -------------------------------------------
$sumsPath = Join-Path $ReleaseDir $sumsFile
if (-not (Test-Path -LiteralPath $sumsPath)) {
    Add-Failure $failures "$sumsFile missing from $ReleaseDir"
}
else {
    $recorded = @{}
    $lineNo = 0
    foreach ($raw in (Get-Content -LiteralPath $sumsPath)) {
        $lineNo++
        $line = $raw.Trim()
        if ($line -eq "") { continue }
        $m = [regex]::Match($line, '^([0-9a-fA-F]{64})\s+\*?(.+)$')
        if (-not $m.Success) {
            Add-Failure $failures "$($sumsFile):$($lineNo): malformed checksum line '$line'"
            continue
        }
        $digest = $m.Groups[1].Value.ToLower()
        $name = $m.Groups[2].Value
        if ($name -eq $sumsFile) {
            Add-Failure $failures "$sumsFile must not list itself"
        }
        elseif ($recorded.ContainsKey($name)) {
            Add-Failure $failures "$sumsFile has duplicate entry for '$name'"
        }
        else {
            $recorded[$name] = $digest
        }
    }

    foreach ($name in $expected) {
        if ($name -eq $sumsFile) { continue }
        if (-not $recorded.ContainsKey($name)) {
            Add-Failure $failures "$sumsFile missing entry for '$name'"
            continue
        }
        $path = Join-Path $ReleaseDir $name
        if (-not (Test-Path -LiteralPath $path)) { continue }
        $actual = (Get-FileHash -LiteralPath $path -Algorithm SHA256).Hash.ToLower()
        if ($actual -ne $recorded[$name]) {
            Add-Failure $failures "$sumsFile entry for '$name' = $($recorded[$name]), actual sha256 $actual"
        }
    }

    foreach ($name in ($recorded.Keys | Sort-Object)) {
        if (-not $expectedSet.ContainsKey($name)) {
            Add-Failure $failures "$sumsFile unexpected entry for '$name'"
        }
    }
}

# --- Check 3: manifest / SBOM / conformance version agreement -----------------
$manifestPath = Join-Path $ReleaseDir $manifestName
if (-not (Test-Path -LiteralPath $manifestPath)) {
    Add-Failure $failures "$manifestName missing from $ReleaseDir"
}
else {
    try {
        $manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json

        $manifestVersion = $null
        foreach ($key in @("product_version", "version")) {
            $prop = $manifest.PSObject.Properties[$key]
            if ($prop -and $prop.Value -is [string] -and $prop.Value -ne "") {
                $manifestVersion = $prop.Value
                break
            }
        }
        if (-not $manifestVersion) {
            Add-Failure $failures "$manifestName has no product_version or version field"
        }
        elseif ($manifestVersion -ne $Version) {
            Add-Failure $failures "$manifestName version '$manifestVersion' != '$Version'"
        }
    }
    catch {
        Add-Failure $failures "$($manifestName): $($_.Exception.Message)"
    }
}

$sbomPath = Join-Path $ReleaseDir $sbomName
if (-not (Test-Path -LiteralPath $sbomPath)) {
    Add-Failure $failures "$sbomName missing from $ReleaseDir"
}
else {
    try {
        $sbom = Get-Content -LiteralPath $sbomPath -Raw | ConvertFrom-Json
        $sbomVersion = $sbom.metadata.component.version
        if ($sbomVersion -ne $Version) {
            Add-Failure $failures "$sbomName component version '$sbomVersion' != '$Version'"
        }
    }
    catch {
        Add-Failure $failures "$($sbomName): $($_.Exception.Message)"
    }
}

$conformancePath = Join-Path $ReleaseDir $conformanceName
if (-not (Test-Path -LiteralPath $conformancePath)) {
    Add-Failure $failures "$conformanceName missing from $ReleaseDir"
}
else {
    $report = Get-Content -LiteralPath $conformancePath -Raw
    $m = [regex]::Match($report, '(?im)^\s*["'']?product_version["'']?\s*[:=]\s*["'']?([0-9]+\.[0-9]+\.[0-9]+)["'']?\s*$')
    if (-not $m.Success) {
        Add-Failure $failures "$conformanceName has no product_version line"
    }
    elseif ($m.Groups[1].Value -ne $Version) {
        Add-Failure $failures "$conformanceName version '$($m.Groups[1].Value)' != '$Version'"
    }
}

# --- Check 4: manifest provenance (source != distribution, both present) ------
if (Test-Path -LiteralPath $manifestPath) {
    try {
        $manifest = Get-Content -LiteralPath $manifestPath -Raw | ConvertFrom-Json
        $sourceProp = $manifest.PSObject.Properties["source_commit"]
        $distProp = $manifest.PSObject.Properties["distribution_commit"]
        $sourceCommit = if ($sourceProp) { [string]$sourceProp.Value } else { "" }
        $distributionCommit = if ($distProp) { [string]$distProp.Value } else { "" }

        if (-not $sourceCommit) {
            Add-Failure $failures "$manifestName is missing source_commit"
        }
        if (-not $distributionCommit) {
            Add-Failure $failures "$manifestName is missing distribution_commit"
        }
        if ($sourceCommit -and $sourceCommit -eq $distributionCommit) {
            Add-Failure $failures "$manifestName conflates source_commit and distribution_commit ($sourceCommit)"
        }
    }
    catch {
        Add-Failure $failures "$manifestName provenance: $($_.Exception.Message)"
    }
}

# --- Report -------------------------------------------------------------------
Write-Host "verify-release: dir=$ReleaseDir version=$Version"
if ($failures.Count -gt 0) {
    Write-Host "[FAIL] $($failures.Count) problem(s):"
    foreach ($f in $failures) {
        Write-Host "  - $f"
    }
    exit 1
}

Write-Host "[OK] all release verification checks passed"
exit 0
