setlocal
@rem set environment variables from `.env` -- does not account for quotes
FOR /F "tokens=*" %%i in ('type .env') do SET %%i

if "%1" == "" (
    echo "Using: Build and run exe mode"
    go build
    knowledgebase.exe
    exit /b 1
)

if "%1" == "--help" (
    echo "Usage: %0 <command>"
    exit /b 1
)

@rem runs the command with env
%*

endlocal