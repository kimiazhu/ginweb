
:: Global go install batch

:: For all platform, please refer to https://golang.org/doc/install/source

::@echo off
set GOOS=linux
set GOARCH=amd64
go install -ldflags "-s -w" takov2


del bin\linux_amd64\takov2_rc
rename bin\linux_amd64\takov2 takov2_rc

::set GOOS=linux
::set GOARCH=arm
::go install -ldflags "-s -w" github.com/yaosxi/mgox
::go install -ldflags "-s -w" takows
::go install -ldflags "-s -w" takofs
::go install -ldflags "-s -w" takoup

::set GOOS=linux
::set GOARCH=386
::go install -ldflags "-s -w" github.com/yaosxi/mgox
::go install -ldflags "-s -w" takows
::go install -ldflags "-s -w" takofs
::go install -ldflags "-s -w" takoup

::set GOOS=darwin
::set GOARCH=amd64
::go install -ldflags "-s -w" github.com/yaosxi/mgox
::go install -ldflags "-s -w" takows
::go install -ldflags "-s -w" takofs
::go install -ldflags "-s -w" takoup
::go install -ldflags "-s -w" takosign

::set GOOS=darwin
::set GOARCH=386
::go install -ldflags "-s -w" github.com/yaosxi/mgox
::go install -ldflags "-s -w" takows
::go install -ldflags "-s -w" takofs
::go install -ldflags "-s -w" takoup

set GOOS=windows
set GOARCH=amd64
go install -ldflags "-s -w" takov2

::set GOOS=windows
::set GOARCH=386
::go install -ldflags "-s -w" github.com/yaosxi/mgox
::go install -ldflags "-s -w" takows
::go install -ldflags "-s -w" takofs
::go install -ldflags "-s -w" takoup