# ngrok
ngrok 

配置参照网址 
https://tonybai.com/2015/03/14/selfhost-ngrok-service/
https://blog.csdn.net/sunansheng/article/details/48372149

生成自签名证书
```text
使用ngrok.com官方服务时，我们使用的是官方的SSL证书。
自建ngrokd服务，我们需要生成自己的证书，并提供携带该证书的ngrok客户端。

证书生成过程需要一个NGROK_BASE_DOMAIN。 
以ngrok官方随机生成的地址693c358d.ngrok.com为例，
其NGROK_BASE_DOMAIN就是"ngrok.com"，
如果你要 提供服务的地址为"example.chan.yanghy.cn"，
那NGROK_BASE_DOMAIN就应该 是"chan.yanghy.cn"。

以NGROK_BASE_DOMAIN="chan.yanghy.cn"为例
```
```text
cd ~/goproj/src/github.com/inconshreveable/ngrok
openssl genrsa -out rootCA.key 2048
openssl req -x509 -new -nodes -key rootCA.key -subj "/CN=chan.yanghy.cn" -days 5000 -out rootCA.pem
openssl genrsa -out device.key 2048
openssl req -new -key device.key -subj "/CN=chan.yanghy.cn" -out device.csr
openssl x509 -req -in device.csr -CA rootCA.pem -CAkey rootCA.key -CAcreateserial -out device.crt -days 5000

```
```text
执行完以上命令，在ngrok目录下就会新生成6个文件：
device.crt
device.csr
device.key
rootCA.key
rootCA.pem
rootCA.srl

ngrok通过bindata将ngrok源码目录下的assets目录（资源文件）
打包到可执行文件(ngrokd和ngrok)中去，
assets/client/tls和assets/server/tls下分别存放着用于ngrok和ngrokd的默认证书文件，
我们需要将它们替换成我们自己生成的：(因此这一步务必放在编译可执行文件之前)

cp rootCA.pem assets/client/tls/ngrokroot.crt
cp device.crt assets/server/tls/snakeoil.crt
cp device.key assets/server/tls/snakeoil.key
```
服务端和客户端配置
```text
1、启动ngrokd
ngrokd -domain="chan.yanghy.cn" -httpAddr=":80" -httpsAddr=":443" -tunnelAddr=":4443"

2、启动ngrok方式1 还是使用方式2吧
配置文件ngrok.cfg
server_addr: "chan.yanghy.cn:4443"
trust_host_root_certs: false
执行 ngrok -subdomain example -config=ngrok.cfg 8080
或 ngrok -config=ngrok.cfg -subdomain demo 80

3、启动ngrok方式2
配置文件debug.yml
server_addr: yanghy.cn:4443
trust_host_root_certs: false
tunnels:
  demo:
    proto:
      http: 8080
执行 ngrok -config=debug.yml -log=ngrok.log -subdomain=demo 8080
```
```text
注意事项

客户端ngrok.cfg中server_addr后的值必须严格与-domain以及证书中的NGROK_BASE_DOMAIN相同，否则Server端就会出现如下错误日志：

[03/13/15 09:55:46] [INFO] [tun:15dd7522] New connection from 54.149.100.42:38252
[03/13/15 09:55:46] [DEBG] [tun:15dd7522] Waiting to read message
[03/13/15 09:55:46] [WARN] [tun:15dd7522] Failed to read message: remote error: bad certificate
[03/13/15 09:55:46] [DEBG] [tun:15dd7522] Closing
```