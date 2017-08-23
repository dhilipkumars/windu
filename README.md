# windu
A simple disk utilization program mainly for windows.
This is similar to unix `du` utility


Here is how it works, if you want to narrow down to the heaviest directory

Check the top level directory %GOPATH% and narrow down from there. 
```
D:\>windu %GOPATH%
332681596  317M src
122958886  117M pkg
 32177152   30M bin

D:\>windu %GOPATH%\pkg
122958886  117M windows_amd64

D:\>windu %GOPATH%\pkg\windows_amd64
74140856   70M k8s.io
47910080   45M github.com
  907950  886K golang.org

D:\>windu %GOPATH%\bin
12087296   11M gocode.exe
11481088   10M godep.exe
 6622720    6M golint.exe
 1986048    1M windu.exe
```
