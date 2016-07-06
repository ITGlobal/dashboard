@echo off
setlocal enableextensions

if "%1" == "" (
    call :_MAKE_VERSION
    
    call :_BUILD_WIN32
    call :_BUILD_WIN64

    call :_BUILD_LINUX32
    call :_BUILD_LINUX64
    call :_BUILD_LINUX_ARM

    goto _EXIT
)

if "%1" == "win32" (
    call :_MAKE_VERSION
    call :_BUILD_WIN32
    goto _EXIT 
)

if "%1" == "win64" (
    call :_MAKE_VERSION
    call :_BUILD_WIN64
    goto _EXIT 
)

if "%1" == "linux32" (
    call :_MAKE_VERSION
    call :_BUILD_LINUX32
    goto _EXIT 
)


if "%1" == "linux64" (
 call :_MAKE_VERSION
    call :_BUILD_LINUX64
    goto _EXIT 
)

if "%1" == "linuxarm" (
    call :_MAKE_VERSION
    call :_BUILD_LINUX_ARM
    goto _EXIT 
)

echo Invalid argument! 
echo Use any of: 
echo .\build 
echo .\build win32 
echo .\build win64 
echo .\build linux32 
echo .\build linux64 
echo .\build linuxarm
goto _EXIT 

rem =========================================================================== 
rem _MAKE_VERSION
rem ===========================================================================
:_MAKE_VERSION
if not exist .\out mkdir .\out
git describe --tags > .\out\version.txt
set /p BUILD_VERSION= < .\out\version.txt
echo Building IT Global Dashboard v%BUILD_VERSION%
exit /B 0
rem ===========================================================================  


rem =========================================================================== 
rem _BUILD_WIN32
rem ===========================================================================  
:_BUILD_WIN32

set GOARCH=386
set GOOS=windows

if not exist .\out mkdir .\out
if not exist .\out\win32 mkdir .\out\win32

go build -o .\out\win32\dashboard.exe -ldflags "-X main.Version=v%BUILD_VERSION% -X main.BuildConfiguration=[win32]"
7z a .\out\dashboard-win32.zip .\out\win32\dashboard.exe > nul

echo Compiled Dashboard v%BUILD_VERSION% for Win32:     See .\out\dashboard-win32.zip

exit /B %ERRORLEVEL%
rem ===========================================================================


rem =========================================================================== 
rem _BUILD_WIN64
rem ===========================================================================  
:_BUILD_WIN64

set GOARCH=amd64
set GOOS=windows

if not exist .\out mkdir .\out
if not exist .\out\win64 mkdir .\out\win64

go build -o .\out\win64\dashboard.exe -ldflags "-X main.Version=v%BUILD_VERSION% -X main.BuildConfiguration=[win64]"
7z a .\out\dashboard-win64.zip .\out\win64\dashboard.exe > nul

echo Compiled Dashboard v%BUILD_VERSION% for Win64:     See .\out\dashboard-win64.zip

exit /B %ERRORLEVEL%
rem ===========================================================================


rem =========================================================================== 
rem _BUILD_LINUX32
rem ===========================================================================  
:_BUILD_LINUX32

set GOARCH=386
set GOOS=linux

if not exist .\out mkdir .\out
if not exist .\out\linux32 mkdir .\out\linux32

go build -o .\out\linux32\dashboard -ldflags "-X main.Version=v%BUILD_VERSION% -X main.BuildConfiguration=[linux32]"
7z a .\out\dashboard-linux32.zip .\out\linux32\dashboard > nul

echo Compiled Dashboard v%BUILD_VERSION% for Linux32:   See .\out\dashboard-linux32.zip

exit /B %ERRORLEVEL%
rem ===========================================================================


rem =========================================================================== 
rem _BUILD_LINUX64
rem ===========================================================================  
:_BUILD_LINUX64

set GOARCH=amd64
set GOOS=linux

if not exist .\out mkdir .\out
if not exist .\out\linux64 mkdir .\out\linux64

go build -o .\out\linux64\dashboard -ldflags "-X main.Version=v%BUILD_VERSION% -X main.BuildConfiguration=[linux64]"
7z a .\out\dashboard-linux64.zip .\out\linux64\dashboard > nul

echo Compiled Dashboard v%BUILD_VERSION% for Linux64:   See .\out\dashboard-linux64.zip

exit /B %ERRORLEVEL%
rem ===========================================================================


rem =========================================================================== 
rem _BUILD_LINUX_ARM
rem ===========================================================================  
:_BUILD_LINUX_ARM

set GOARCH=amd64
set GOOS=linux

if not exist .\out mkdir .\out
if not exist .\out\arm mkdir .\out\arm

go build -o .\out\arm\dashboard -ldflags "-X main.Version=v%BUILD_VERSION% -X main.BuildConfiguration=[arm]"
7z a .\out\dashboard-arm.zip .\out\arm\dashboard > nul

echo Compiled Dashboard v%BUILD_VERSION% for Linux ARM: See .\out\dashboard-arm.zip

exit /B %ERRORLEVEL%
rem ==========================================================================


:_EXIT
endlocal