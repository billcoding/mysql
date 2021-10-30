# mysql
A mysql client(or terminal) written in Golang

# install
```shell
go install github.com/billcoding/mysql@latest
```

# usage

## on tty
```shell
mysql -H="x.y.z.w" -P="3307" -u="someone" -p="passwd" -d="some_db" 
```

## command only
```shell
mysql -H="x.y.z.w" -P="3307" -u="someone" -p="passwd" -d="some_db" -c="SELECT NOW() AS t" 
```