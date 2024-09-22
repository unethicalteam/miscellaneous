# Check if the 'hashes' directory exists, and if not, create it
if (-Not (Test-Path "hashes")) {
    Write-Host "Creating 'hashes' directory..."
    try {
        $null = New-Item -Path "hashes" -ItemType Directory -ErrorAction Stop
    } catch {
        Write-Host "Failed to create 'hashes' directory: $($_.Exception.Message)"
        exit 1
    }
} else {
    Write-Host "'hashes' directory already exists."
}

# Get all files in the current directory, excluding 'hasher.ps1'
$files = Get-ChildItem -File | Where-Object { $_.Name -ne 'hasher.ps1' }

# Initialize a hashtable to store file occurrence counts
$fileCount = @{}

# Count occurrences of each unique file basename
Write-Verbose "Counting file occurrences..."
foreach ($file in $files) {
    $fileCount[$file.BaseName] = ($fileCount[$file.BaseName] ?? 0) + 1
}

# Calculate SHA256 hash for each file
Write-Host "Calculating file hashes..."
foreach ($file in $files) {
    $filename = $file.BaseName
    $extension = $file.Extension
    $count = $fileCount[$filename]
    $hashFilename = if ($count -gt 1) { "${filename}${extension}.sha256" } else { "${filename}.sha256" }

    Write-Host "Processing ${filename}${extension}..."
    try {
        $hash = (Get-FileHash -Path $file.FullName -Algorithm SHA256 -ErrorAction Stop).Hash
        $hash | Out-File -FilePath "hashes\$hashFilename" -ErrorAction Stop
        Write-Host "Hashing of ${filename}${extension} completed successfully."
    } catch {
        Write-Host "An error occurred while hashing $($file.FullName): $($_.Exception.Message)"
    }
}
