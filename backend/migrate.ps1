param(
    [string]$Direction = "up"
)

$ErrorActionPreference = "Stop"

$container = "mysql-task"
$user      = "root"
$password  = "root"
$database  = "taskdb"

$migrationsDir = Join-Path $PSScriptRoot "migrations"

if (-not (Test-Path $migrationsDir)) {
    Write-Error "Folder migrations tidak ditemukan: $migrationsDir"
}

if ($Direction -eq "up") {
    $files = Get-ChildItem $migrationsDir -Filter "*.up.sql" | Sort-Object Name
} elseif ($Direction -eq "down") {
    $files = Get-ChildItem $migrationsDir -Filter "*.down.sql" | Sort-Object Name -Descending
} else {
    Write-Error "Direction harus 'up' atau 'down'"
}

if ($files.Count -eq 0) {
    Write-Host "Tidak ada migration $Direction."
    exit 0
}

foreach ($file in $files) {
    Write-Host "→ Applying $($file.Name) ($Direction)"
    Get-Content $file.FullName | docker exec -i $container mysql -u$user -p$password $database
}

Write-Host "✅ Migration $Direction selesai."