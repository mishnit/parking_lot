Install GO on Ubuntu:
```
$ sudo apt update
$ wget https://dl.google.com/go/go1.13.1.linux-amd64.tar.gz
$ sudo tar -xvf go1.13.1.linux-amd64.tar.gz
$ sudo mv go /usr/local
$ sudo vi ~/.profile
(paste the below line at end, save and exit)
>> export GOROOT=/usr/local/go
>> export GOPATH=$HOME/go
>> export PATH=$PATH:$GOROOT/bin:$GOPATH/bin
$ source ~/.profile
$ go version [to check golang version]
$ go env [to check golang internal env variables]
```

Install Docker and Docker Compose on Ubuntu:
```
$ sudo apt install apt-transport-https ca-certificates curl software-properties-common
$ curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
$ sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu bionic stable"
$ sudo apt update
$ apt-cache policy docker-ce
$ sudo apt install docker-ce
$ docker --version
$ sudo curl -L https://github.com/docker/compose/releases/download/1.21.2/docker-compose-`uname -s`-`uname -m` -o /usr/local/bin/docker-compose
$ sudo chmod +x /usr/local/bin/docker-compose
$ docker-compose --version
