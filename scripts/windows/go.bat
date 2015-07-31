if not exist "%TEMP%\go.msi" (
  powershell -Command "(New-Object System.Net.WebClient).DownloadFile('https://storage.googleapis.com/golang/go1.4.2.windows-amd64.msi', '%TEMP%\go.msi')" <NUL
)

msiexec /quiet /i "%TEMP%\go.msi

mkdir "%USERPROFILE%\go"
mkdir "%USERPROFILE%\go\bin"
powershell -Command "$p = [Environment]::GetEnvironmentVariable('PATH', 'Machine'); $pa = [Environment]::GetEnvironmentVariable('USERPROFILE', 'Process'); [Environment]::SetEnvironmentVariable('GOPATH', \"$pa\go\", 'Machine'); [Environment]::SetEnvironmentVariable('PATH', \"$p;$pa\go\bin\", 'Machine');"
