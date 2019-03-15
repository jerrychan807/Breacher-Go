# Breacher-Go

A script to find admin login pages and EAR vulnerabilites.

## Features

- [x] Big path list (482 paths)
- [x] Supports php, asp and html extensions
- [x] Checks for robots.txt
- [x] use goroutine
- [x] Reduce error report
- [ ] Checks for potential EAR vulnerabilites
- [ ] Support for custom patns
- [ ] Remove the same


## Screenshot:

![](https://ws2.sinaimg.cn/large/006tKfTcgy1g139bma45pj30u00uokjl.jpg)


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

## Links:

Blog post:

- [Breacher-Go高并发管理员后台爆破工具](https://jerrychan807.github.io/2019/02/23/Breacher-Go%E5%B9%B6%E5%8F%91%E7%9A%84%E7%AE%A1%E7%90%86%E5%91%98%E5%90%8E%E5%8F%B0%E7%88%86%E7%A0%B4%E5%B7%A5%E5%85%B7/)

## Acknowledgments

This Repository is Go-verion : [s0md3v/Breacher](https://github.com/s0md3v/Breacher)

