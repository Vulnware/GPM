packr2 
go build -o ./bin/gpm_cli.exe .
go build -ldflags "-s -w" -o ./bin/gpm_cli_light.exe .
go build -ldflags "-s -w" -o ./bin/gpm_cli_upx.exe .
tinygo build -o ./bin/gpm_cli_tinygo.exe .
upx -9 ./bin/gpm_cli_upx.exe
copy .\bin\gpm_cli.exe .\gpm_cli.exe
