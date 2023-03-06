# Kitty domain

> a simple fast and efficient subdomain miner  


# Quick start

**build** 
```golang
$ go build .
```

**command parameter** 
>    __ __     __  __  _        ___                  _    
>   / //_/__ _/ /_/ /_(_)_ __  / _ \___  __ _  ___ _(_)__ 
>  / ,< / _ ./ __/ __/ / // / / // / _ \/  ' \/ _ ./ / _ \
> /_/|_|\_,_/\__/\__/_/\_, / /____/\___/_/_/_/\_,_/_/_//_/
>                     /___/                               
> 		kitty domain version: 0.3
> 
>   -d string
>     	target domain
>   -dns string
>     	dns server (default "8.8.8.8")
>   -f string
>     	set subdomain dict file (default "subdomains.txt")
>   -fd string
>     	fuzz domain the placeholder is ? example: "w?.baidu.com" -> "www.baidu.com"
>   -o string
>     	set out put file
>   -rd string
>     	use regexp fuzz domain must be used -reg example: -rd "w?.baidu.com" -reg "\?"
>   -reg string
>     	set regexp for -rd  



**example**
```shell
$ ./kdomain -d "baidu.com"
$ ./kdomain -fd "ww?.baidu.com"
```
