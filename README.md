 pre-requisites: make sure you have go installed: 
```
go -v
```
or
```
go --v
```



# usage >

for setting headers while testing:
```
go run main.go <URL> [-H key:value]
```

without settings headers:
```
go run main.go <URL>
```

with streaming enabled: (streaming gives live response data back.)
```
go run main.go <URL> [-H key:value] -stream <true/false>
```
