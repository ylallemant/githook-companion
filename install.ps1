#!/usr/bin/env pwsh
# inherit from https://deno.land/x/install@v0.1.4/install.ps1
# Copyright 2018 the Deno authors. All rights reserved. MIT license.

# required:
# 1. $repo or $r
# 2. $version or $v
# 2. $exe or $e

$ErrorActionPreference = 'Stop'
$inputVersion = if ($version) {
  "${version}"
} else {
  "${v}"
}

$arr = $inputRepo.Split('/')

$owner = "ylallemant"
$repoName = "githooks-butler"
$exeName = "githooks-butler"
$separator = "-"
$githubUrl = "https://github.com"

if ($exeName -eq "") {
  $exeName = "${repoName}"
}

if ($inputVersion) {
  $version = "${inputVersion}"
}

if ([Environment]::Is64BitProcess) {
  $arch = "amd64"
} else {
  $arch = "386"
}

$BinDir = "$Home\bin"
$downloadedTagGz = "$BinDir\${exeName}.tar.gz"
$downloadedExe = "$BinDir\${exeName}.exe"
$fileName = "${exeName}${separator}${version}${separator}windows${separator}${arch}.tar.gz"

# GitHub requires TLS 1.2
[Net.ServicePointManager]::SecurityProtocol = [Net.SecurityProtocolType]::Tls12

$ResourceUri = "${githubUrl}/${owner}/${repoName}/releases/download/${Version}/${fileName}"

if (!(Test-Path $BinDir)) {
  New-Item $BinDir -ItemType Directory | Out-Null
}

Invoke-WebRequest $ResourceUri -OutFile $downloadedTagGz -UseBasicParsing -ErrorAction Stop

function Check-Command {
  param($Command)
  $found = $false
  try
  {
      $Command | Out-Null
      $found = $true
  }
  catch [System.Management.Automation.CommandNotFoundException]
  {
      $found = $false
  }

  $found
}

if (Check-Command -Command tar) {
  Invoke-Expression "tar -xvzf $downloadedTagGz -C $BinDir"
} else {
  function Expand-Tar($tarFile, $dest) {

      if (-not (Get-Command Expand-7Zip -ErrorAction Ignore)) {
          Install-Package -Scope CurrentUser -Force 7Zip4PowerShell > $null
      }

      Expand-7Zip $tarFile $dest
  }

  Expand-Tar $downloadedTagGz $BinDir
}

Remove-Item $downloadedTagGz

$User = [EnvironmentVariableTarget]::User
$Path = [Environment]::GetEnvironmentVariable('Path', $User)
if (!(";$Path;".ToLower() -like "*;$BinDir;*".ToLower())) {
  [Environment]::SetEnvironmentVariable('Path', "$Path;$BinDir", $User)
  $Env:Path += ";$BinDir"
}

Write-Output "${exeName} was installed successfully to $downloadedExe"
Write-Output "Run '${exeName} --help' to get started"