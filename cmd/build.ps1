# Define the file prefix for the output files
$filePrefix = "youdao"

# Define the main output directory
$mainOutputDir = "build"

# Define the output directories for different platforms and architectures with "youdao" prefix
$outputDirs = @{
    "windows_amd64" = "$mainOutputDir\youdao_windows_amd64"
    "windows_arm64" = "$mainOutputDir\youdao_windows_arm64"
    "linux_amd64"   = "$mainOutputDir\youdao_linux_amd64"
    "linux_arm64"   = "$mainOutputDir\youdao_linux_arm64"
    "mac_amd64"     = "$mainOutputDir\youdao_mac_amd64"
    "mac_arm64"     = "$mainOutputDir\youdao_mac_arm64"
}

# Create the main output directory if it doesn't exist
if (-not (Test-Path $mainOutputDir)) {
    New-Item -ItemType Directory -Path $mainOutputDir | Out-Null
}

# Create the output directories if they don't exist
foreach ($dir in $outputDirs.Values) {
    if (-not (Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir | Out-Null
    }
}

# Function to build for a specific platform and architecture
function Build-ForPlatform {
    param (
        [string]$goos,
        [string]$goarch,
        [string]$outputDir
    )

    $env:GOOS = $goos
    $env:GOARCH = $goarch

    $outputFile = "$outputDir\$filePrefix"
    if ($goos -eq "windows") {
        $outputFile += ".exe"
    }

    Write-Output "Building for $goos/$goarch..."
    go build -o $outputFile .

    # Clean up the environment variables
    Remove-Item Env:GOOS
    Remove-Item Env:GOARCH
}

# Function to create a zip file from a directory
function Create-Zip {
    param (
        [string]$dir
    )

    $zipFile = "$dir.zip"

    Write-Output "Creating zip file $zipFile..."
    Compress-Archive -Path $dir -DestinationPath $zipFile -Update
}

# Function to remove a directory
function Remove-Directory {
    param (
        [string]$dir
    )

    if (Test-Path $dir) {
        Write-Output "Removing directory $dir..."
        Remove-Item -Path $dir -Recurse -Force
    }
}

# Build for Windows AMD64
Build-ForPlatform -goos "windows" -goarch "amd64" -outputDir $outputDirs["windows_amd64"]
Create-Zip -dir $outputDirs["windows_amd64"]
Remove-Directory -dir $outputDirs["windows_amd64"]

# Build for Windows ARM64
Build-ForPlatform -goos "windows" -goarch "arm64" -outputDir $outputDirs["windows_arm64"]
Create-Zip -dir $outputDirs["windows_arm64"]
Remove-Directory -dir $outputDirs["windows_arm64"]

# Build for Linux AMD64
Build-ForPlatform -goos "linux" -goarch "amd64" -outputDir $outputDirs["linux_amd64"]
Create-Zip -dir $outputDirs["linux_amd64"]
Remove-Directory -dir $outputDirs["linux_amd64"]

# Build for Linux ARM64
Build-ForPlatform -goos "linux" -goarch "arm64" -outputDir $outputDirs["linux_arm64"]
Create-Zip -dir $outputDirs["linux_arm64"]
Remove-Directory -dir $outputDirs["linux_arm64"]

# Build for macOS AMD64
Build-ForPlatform -goos "darwin" -goarch "amd64" -outputDir $outputDirs["mac_amd64"]
Create-Zip -dir $outputDirs["mac_amd64"]
Remove-Directory -dir $outputDirs["mac_amd64"]

# Build for macOS ARM64
Build-ForPlatform -goos "darwin" -goarch "arm64" -outputDir $outputDirs["mac_arm64"]
Create-Zip -dir $outputDirs["mac_arm64"]
Remove-Directory -dir $outputDirs["mac_arm64"]

Write-Output "Builds complete. Check the zip files in the build directory."
