# FTP

simple client-server ftp implementation in golang

## ftpClient
currently implements the following features:
- DIR
- CD
- PWD

## ftpServer
communicates with ftpClient over TCP socket to reveal directory info.

hardcoded to use localhost port 1202

currently implements:
- DIR
- CD
- PWD

