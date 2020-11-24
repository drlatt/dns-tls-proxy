**DNS to DNS-over-TLS proxy**

This application implements a DNS proxy that listens to conventional DNS and sends it over TLS.

Technology used
1. Golang
    - I decided to use golang because the language is easy to use and understand. Golang is easily understood by anyone going through code written in the language.
    - The language has a robust standard library which offers lots of functionality inbuilt without the engineer needing to import lots of external libraries in order to implement application functionality.
    - Majority of golang's libraries provide error handling functionality by default which is a nice to have feature.
2. Cloudflare-dns
    - Cloudflare dns was used as the provider which handles DNS queries over a TLS connection.

**Implementation**\
This is a simple dns-over-tls-proxy application that listens on TCP port 53 for client DNS queries. The proxy then dials the upstream server(server which handles DNS requests) using a TLS network connection.

Once the TLS connection is successful, the proxy then forwards the client's DNS query to the upstream server.

The upstream server responds with a DNS answer which the proxy then forwards to the client.

- Upstream server: server which handles DNS queries e.g. Cloudflare, Rackspace, ClouDNS etc.
- Client: computer/device which sends DNS requests to be resolved by an external DNS provider

**How to use**
- Clone repo or download zip file then unzip
- Navigate to application folder containing Dockerfile, .env file and main.go file
- Create a .env file in the format below
    ```
    CONN_PORT=53
    UPSTREAM_SERVER=1.1.1.1
    UPSTREAM_SERVER_PORT=853
    ```
- Run `docker build -t dns_tls_proxy .` to build a docker image
- Run `docker run -p 53:53 --env-file .env --name dns dns_tls_proxy` to run a docker container running the proxy solution. This command also maps TCP port 53 on the local machine to port 53 on the running container.
- Once the docker container starts successfully, DNS queries can be sent to the proxy application

*Testing the proxy*\
`dig @localhost +tcp primark.de`

response
```
lateefkoshemani@DFSs-MBP  ~/Documents/repos/myprojects/dns-over-tls/misc2  dig @localhost +tcp primark.de 
; <<>> DiG 9.10.6 <<>> @localhost +tcp primark.de
; (2 servers found)
;; global options: +cmd
;; Got answer:
;; ->>HEADER<<- opcode: QUERY, status: NOERROR, id: 39016
;; flags: qr rd ra; QUERY: 1, ANSWER: 1, AUTHORITY: 0, ADDITIONAL: 1

;; OPT PSEUDOSECTION:
; EDNS: version: 0, flags:; udp: 1452
; PAD: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 (".....................................................................")
;; QUESTION SECTION:
;primark.de.                    IN      A

;; ANSWER SECTION:
primark.de.             3599    IN      A       217.19.248.132

;; Query time: 156 msec
;; SERVER: 127.0.0.1#53(127.0.0.1)
;; WHEN: Mon Oct 14 14:51:12 CEST 2019
;; MSG SIZE  rcvd: 128
```
