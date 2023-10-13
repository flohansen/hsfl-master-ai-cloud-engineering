$ROOT_DIR = $PWD

Get-ChildItem -Path ./ -Filter *.mod -Recurse -ErrorAction SilentlyContinue -Force | ForEach-Object {
    $currentDir = $_.DirectoryName
    Write-Output ""
    Write-Output "### $currentDir ###"
    Write-Output ""
    
    Set-Location $currentDir

    $VOID = (go mod tidy)
    $VOID = (go clean -testcache)

    $VOID = (go test ./... -coverprofile=cover)

    if (-Not $?) {
      Remove-Item -Force cover
      exit 1
    }

    $COVERAGE = (go tool cover -func cover | Select-String -Pattern 'total')

    Remove-Item -Force cover


    $COVERAGE_NUMBER = Write-Output $COVERAGE # | sed -e 's/\.[0-9]*//' -e 's/%//'

    Write-Output $COVERAGE_NUMBER

    Set-Location $ROOT_DIR
    
}