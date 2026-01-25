
$version = Get-Content .\version.txt
Write-Host $version
git tag $version
git push origin $version

