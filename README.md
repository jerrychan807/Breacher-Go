# Breacher-Go

A script to find admin login pages and EAR vulnerabilites.

## Acknowledgments

This Repository is Go-verion : [s0md3v/Breacher](https://github.com/s0md3v/Breacher)



## Features

- [x] Big path list (482 paths)
- [x] Supports php, asp and html extensions
- [x] Checks for robots.txt
- [ ] Checks for potential EAR vulnerabilites
- [x] use goroutine
- [ ] Support for custom patns
- [x] Reduce error report
- [ ] Remove the same

## Usages

- Check all paths with php extension
```
go run breacher.go -u "example.com" -t "php"
```

- Check all paths with php extension with goroutine
```
go run breacher.go -u example.com --type php --fast
```
- Check all paths without goroutine
```
go run breacher.go -u example.com
```


