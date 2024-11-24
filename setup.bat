@echo off
echo Installing development tools...

REM Install solc
echo Installing solc...
go run tools/solc/install.go
if errorlevel 1 goto error

REM Install abigen
echo Installing abigen...
go run tools/abigen/install.go
if errorlevel 1 goto error

echo Setup complete! Tools installed in .\bin directory
goto :eof

:error
echo Failed to install tools
exit /b 1