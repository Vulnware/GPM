packr2 
go build -o ./bin/gpm.exe .
go build -ldflags "-s -w" -o ./bin/gpm_light.exe .
go build -ldflags "-s -w" -o ./bin/gpm_upx.exe .
tinygo build -o ./bin/gpm_tinygo.exe .
upx -9 ./bin/gpm_upx.exe
copy .\bin\gpm.exe .\gpm.exe
