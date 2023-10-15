# Check if the 'hashes' directory exists, and if not, create it
if (-Not (Test-Path "hashes")) {
    Write-Host "Creating 'hashes' directory..."
    try {
        New-Item -Path "hashes" -ItemType Directory -ErrorAction Stop
    } catch {
        Write-Host "Failed to create 'hashes' directory: $($_.Exception.Message)"
        exit 1
    }
} else {
    Write-Host "'hashes' directory already exists."
}

# Initialize a hashtable to store file occurrence counts
Write-Verbose "Initializing hashtable for file occurrences..."
$fileCount = @{}

# Count occurrences of each unique file basename
Write-Verbose "Counting file occurrences..."
Get-ChildItem -File | ForEach-Object {
    if ($_.Name -eq 'hasher.ps1') {
        Write-Verbose "Skipping $($_.Name)"
        return
    }

    if ($fileCount.ContainsKey($_.BaseName)) {
        $fileCount[$_.BaseName] += 1
    } else {
        $fileCount[$_.BaseName] = 1
    }
}

# Calculate SHA256 hash for each file
Write-Host "Calculating file hashes..."
Get-ChildItem -File | ForEach-Object {
    if ($_.Name -eq 'hasher.ps1') {
        Write-Verbose "Skipping $($_.Name)"
        return
    }

    $filename = $_.BaseName
    $extension = $_.Extension
    $count = $fileCount[$filename]
    $hashFilename = if ($count -gt 1) { "${filename}${extension}.sha256" } else { "${filename}.sha256" }

    Write-Host "Processing ${filename}${extension}..."
    try {
        $hash = (Get-FileHash -Path $_.FullName -Algorithm SHA256 -ErrorAction Stop).Hash
        $hash | Out-File -FilePath "hashes\$hashFilename" -ErrorAction Stop
        Write-Host "Hashing of ${filename}${extension} completed successfully."
    } catch {
        Write-Host "An error occurred while hashing $($_.FullName): $($_.Exception.Message)"
    }
}
