# ddc-parser

server,go


## Config
- default path
```text
~/./ddc-parser/config.toml
```
- configs/config.toml
```text

[database]
addrs= "localhost:27018"
user= "iris"
passwd= "irispassword"
database= "bifrost-sync"


[server]
prometheus_port="9092"
max_operate_tx_count=100
insert_batch_limit=100
incre_height=10



[ddc_client]
gateway_url="http://192.168.150.42:8545"
gateway_api_key=""
gateway_api_value=""
authority_address="0xBcE9AA1924D7197C9C945e43638Bf589f91bcB71"
charge_address="0xF41b6185bFB22E2EFC5fB8395Fa3B952951E2d0b"
ddc_721_address="0x74b6114d011891Ac21FD1d586bc7F3407c63c216"
ddc_1155_address="0x9f7388e114DfDFAbAF8e4b881894E4C7e1b52C17"
log_filepath="./log.log"


```

## Run

```$xslt
make install
ddcparser start
```

## Run with docker
You can run application with docker.

```$xslt
docker build -t ddc-parser .
```

then
```$xslt
docker run --name ddc-parser -p 8080:8080 -v /Users/user/.ddc-parser/config.toml:/root/.ddc-parser/config.toml   ddc-parser
```
