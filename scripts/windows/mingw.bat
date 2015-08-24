if not exist "%TEMP%\7z920-x64.msi" (
    powershell -Command "(New-Object System.Net.WebClient).DownloadFile('http://downloads.sourceforge.net/sevenzip/7z920-x64.msi', '%TEMP%\7z920-x64.msi')" <NUL
)

msiexec /qb /i "%TEMP%\7z920-x64.msi"

pushd "%TEMP%"
rem Make sure the '%' nd '&' in URL are escaped when copy-and-pasting in
if not exist "%TEMP%\mingw.7z" (
  rem Use curl because it can follow redirects
  "%PROGRAMFILES(X86)%\Git\bin\curl.exe" -L "http://downloads.sourceforge.net/project/mingw-w64/Toolchains%%20targetting%%20Win64/Personal%%20Builds/mingw-builds/5.1.0/threads-posix/seh/x86_64-5.1.0-release-posix-seh-rt_v4-rev0.7z?r=http%%3A%%2F%%2Fsourceforge.net%%2Fprojects%%2Fmingw-w64%%2Ffiles%%2FToolchains%%2520targetting%%2520Win64%%2FPersonal%%2520Builds%%2Fmingw-builds%%2F5.1.0%%2Fthreads-posix%%2Fseh%%2F" -o mingw.7z
)
"%PROGRAMFILES%\7-Zip\7z.exe" x mingw.7z -oC:\
rename C:\mingw64 MinGW
rem del mingw.7z
popd

rem msiexec /qb /x "%TEMP%\7z920-x64.msi"

powershell -Command "$p = [Environment]::SetEnvironmentVariable('C_INCLUDE_PATH', 'C:\MinGW\include', 'Machine');"
powershell -Command "$p = [Environment]::SetEnvironmentVariable('LIBRARY_PATH', 'C:\MinGW\lib', 'Machine');"
powershell -Command "$p = [Environment]::GetEnvironmentVariable('PATH', 'Machine'); [Environment]::SetEnvironmentVariable('PATH', \"$p;C:\MinGW\bin\", 'Machine');"
