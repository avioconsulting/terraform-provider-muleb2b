param(
  [Parameter(Mandatory=$True)]
  [string]$version
)

# https://stackoverflow.com/a/46876070/847953
function Expand-Tar($tarFile, $dest) {

    if (-not (Get-Command Expand-7Zip -ErrorAction Ignore)) {
        Install-Package -Scope CurrentUser -Force 7Zip4PowerShell > $null
    }

    Expand-7Zip -ArchiveFileName $tarFile -TargetPath $dest
}

# https://www.codeproject.com/Tips/638039/GZipStream-length-when-uncompressed
function Original-GzipFileSize {
    param(
        [Parameter(Mandatory=$true)]
        [string] $gzipFile
    )

    $fs = New-Object System.IO.FileStream $gzipFile, ([IO.FileMode]::Open), ([IO.FileAccess]::Read), ([IO.FileShare]::Read)

    try
    {
        $fh = New-Object byte[](3)
        $fs.Read($fh, 0, 3) | Out-Null
        # If magic numbers are 31 and 139 and the deflation id is 8 then this is a file to process
        if ($fh[0] -eq 31 -and $fh[1] -eq 139 -and $fh[2] -eq 8)
        {
            $ba = New-Object byte[](4)
            $fs.Seek(-4, [System.IO.SeekOrigin]::End) | Out-Null
            $fs.Read($ba, 0, 4) | Out-Null

            return [int32][System.BitConverter]::ToInt32($ba, 0)
        }
        else
        {
            throw "File '$gzipFile' does not have the correct gzip header"
        }
    }
    finally
    {
        $fs.Close()
    }
}

# https://stackoverflow.com/a/42165686/847953
function Expand-GZip {
    [cmdletbinding(SupportsShouldProcess=$True)]
    param(
        [Parameter(Mandatory=$true)]
        [string]$infile,
        [Parameter(Mandatory=$true)]
        [string]$outFile,
        [int]$bufferSize = 1024
    )
    $fileSize = Original-GzipFileSize $inFile
    $processed = 0

    if ($PSCmdlet.ShouldProcess($infile,"Expand gzip stream")) {
        $input = New-Object System.IO.FileStream $inFile, ([IO.FileMode]::Open), ([IO.FileAccess]::Read), ([IO.FileShare]::Read)
        $output = New-Object System.IO.FileStream $outFile, ([IO.FileMode]::Create), ([IO.FileAccess]::Write), ([IO.FileShare]::None)
        $gzipStream = New-Object System.IO.Compression.GzipStream $input, ([IO.Compression.CompressionMode]::Decompress)

        $buffer = New-Object byte[]($bufferSize)
        while($true){

            $pc = (($processed / $fileSize) * 100) % 100
            Write-Progress "Extracting tar from gzip" -PercentComplete $pc

            $read = $gzipstream.Read($buffer, 0, $bufferSize)

            $processed = $processed + $read

            if ($read -le 0)
            {
                Write-Progress "Extracting tar from gzip" -Completed
                break
            }
            $output.Write($buffer, 0, $read)
        }

        $gzipStream.Close()
        $output.Close()
        $input.Close()
    }
}

$targzFile = "terraform-provider-muleb2b_windows_amd64.tar.gz"
$tarFile = "terraform-provider-muleb2b_windows_amd64.tar"
$url = "https://github.com/avioconsulting/terraform-provider-muleb2b/releases/download/$version/$targzFile"
$output = ".\terraform-provider-muleb2b_windows_amd64.tar.gz"
Invoke-WebRequest -Uri $url -OutFile $output

if (-NOT (Test-Path $env:APPDATA\terraform.d\plugins)) {
  if (-NOT (Test-Path $env:APPDATA\terraform.d)) {
    New-Item -Path "$env:APPDATA" -Name "terraform.d" -ItemType "directory"
  }
  New-Item -Path "$env:APPDATA\terraform.d" -Name "plugins" -ItemType "directory"
}

Expand-Gzip $targzFile $tarFile 16384
Expand-Tar $tarFile "$env:APPDATA\terraform.d\plugins\"

Remove-Item -Path .\terraform-provider-muleb2b_windows_amd64.tar.gz
Remove-Item $tarFile