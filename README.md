![Logo](https://ibin.co/4R6Hzr2H7l4A.png)

#### A RestAPI Microservice GateWay For Linux

A super light weight remote management tool which uses REST API for real time configuration and performance as well as health monitoring for systems (containers) and applications. It provides fast API based monitoring without affecting the system it's running on.


#### Objectives:
- No client installation required. curl is enough.
- No GUI
- Minimal data transfer using JSON


##### Allows you to configure
- systemd
    - services (start, stop, restart, status)
    - service properties for example CPUShares
    - See service logs.
- networkd config
    - .network
    - .netdev
    - .link
- set and get hostname
   - hostnamed

- configure network (netlink)
   - Link: mtu, up, down
   - Create bridge and enslave links
   - Adddress: Set, Get, Delete
   - Gateway: Default Gateway Add and Delete

- see information from /proc such as netstat, netdev, memory
- See ethtool information

##### Quick Start

First configure your GOPATH. If you have already done this skip this step.
```
# keep in ~/.bashrc
```

```
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
export OS_OUTPUT_GOPATH=1

 ```

clone inside src dir of GOPATH. In my case

```
[sus@Zeus src]$ pwd
/home/sus/go/src
```

##### Install libs
```
$ go get github.com/sirupsen/logrus
$ go get github.com/gorilla/mux
$ go get github.com/vishvananda/netlink
$ go get github.com/shirou/gopsutil
$ go get github.com/coreos/go-systemd/dbus
$ go get github.com/godbus/dbus
$ go get github.com/safchain/ethtool
$ go get github.com/go-ini/ini
```

##### Now build it
```
[sus@Zeus src]$ git clone https://github.com/RestGW/api-routerd
[sus@Zeus src]$ cd api-routerd/cmd
[sus@Zeus cmd]$ pwd
/home/sus/go/src/api-routerd/cmd
[sus@Zeus cmd]$ go build -o api-routerd

```

##### How to configure IP and Port ?

Conf dir: ```/etc/api-routerd/```
Conf File: ```api-routerd.conf```

```
$ cat /etc/api-routerd/api-routerd.conf
[Network]
IPAddress=0.0.0.0
Port=8080
```

##### How to configure users ?
Add user name and authentication string in space separated lines

```
# cat /etc/api-routerd/api-routerd-auth.conf
Susant aaaaa
Max bbbb
Joy ccccc

$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netdev"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
```
##### How to configure TLS ?
Generate private key (.key)
```
# Key considerations for algorithm "RSA" ≥ 2048-bit
$ openssl genrsa -out server.key 2048
Generating RSA private key, 2048 bit long modulus (2 primes)
.......................+++++
.+++++
e is 65537 (0x010001)

openssl genrsa -out server.key 2048
```

Generation of self-signed(x509) public key (PEM-encodings .pem|.crt) based on the private (.key)
```
$ openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
You are about to be asked to enter information that will be incorporated
into your certificate request.
What you are about to enter is what is called a Distinguished Name or a DN.
There are quite a few fields but you can leave some blank
For some fields there will be a default value,
If you enter '.', the field will be left blank.
-----
Country Name (2 letter code) [XX]:

```
Place ```server.crt``` and ```server.key``` in the dir /etc/api-routerd/tls

```
[root@Zeus tls]# ls
server.crt  server.key
[root@Zeus tls]# pwd
/etc/api-routerd/tls

```
Now start curl using ** https ***
```
[sus@Zeus tls]$ curl -k --header "X-Session-Token: aaaaa" --header "Content-Type: application/json" --request GET --data '{"action":"get-link-features", "link":"vmnet8"}' https://localhost:8080/network/ethtool/get --tlsv1.2

{"esp-hw-offload":false,"esp-tx-csum-hw-offload":false,"fcoe-mtu":false,"highdma":false,"hw-tc-offload":false,"l2-fwd-offload":false,"loopback":false,"netns-local":false,"rx-all":false,"rx-checksum":false,"rx-fcs":false,"rx-gro":true,"rx-gro-hw":false,"rx-hashing":false,"rx-lro":false,"rx-ntuple-filter":false,"rx-udp_tunnel-port-offload":false,"rx-vlan-filter":false,"rx-vlan-hw-parse":false,"rx-vlan-stag-filter":false,"rx-vlan-stag-hw-parse":false,"tls-hw-record":false,"tls-hw-rx-offload":false,"tls-hw-tx-offload":false,"tx-checksum-fcoe-crc":false,"tx-checksum-ip-generic":false,"tx-checksum-ipv4":false,"tx-checksum-ipv6":false,"tx-checksum-sctp":false,"tx-esp-segmentation":false,"tx-fcoe-segmentation":false,"tx-generic-segmentation":false,"tx-gre-csum-segmentation":false,"tx-gre-segmentation":false,"tx-gso-partial":false,"tx-gso-robust":false,"tx-ipxip4-segmentation":false,"tx-ipxip6-segmentation":false,"tx-lockless":false,"tx-nocache-copy":false,"tx-scatter-gather":false,"tx-scatter-gather-fraglist":false,"tx-sctp-segmentation":false,"tx-tcp-ecn-segmentation":false,"tx-tcp-mangleid-segmentation":false,"tx-tcp-segmentation":false,"tx-tcp6-segmentation":false,"tx-udp-segmentation":false,"tx-udp_tnl-csum-segmentation":false,"tx-udp_tnl-segmentation":false,"tx-vlan-hw-insert":false,"tx-vlan-stag-hw-insert":false,"vlan-challenged":false}
```

```
[sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/misc --tlsv1.2 -k
{"130":"watchdog","144":"nvram","165":"vmmon","183":"hw_random","184":"microcode","227":"mcelog","228":"hpet","229":"fuse","231":"snapshot","232":"kvm","235":"autofs","236":"device-mapper","53":"vboxnetctl","54":"vsock","55":"vmci","56":"vboxdrvu","57":"vboxdrv","58":"rfkill","59":"memory_bandwidth","60":"network_throughput","61":"network_latency","62":"cpu_dma_latency","63":"vga_arbiter"}

[sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/net/arp --tlsv1.2 -k
[{"ip_address":"192.168.225.1","hw_type":"0x1","flags":"0x2","hw_address":"1a:89:20:38:68:8f","mask":"*","device":"wlp4s0"}]

[sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/modules --tlsv1.2 -k
```
Use case: systemd
```
[sus@Zeus] curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-property", "unit":"sshd.service", "property":"CPUShares", "value":"1024"}' http://localhost:8080/service/systemd/property
{"property":"CPUShares","value":"1024"}

[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"start","unit":"sshd.service"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd

[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"stop","unit":"sshd.service"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd

[sus@Zeus]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"status","unit":"sshd.service"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd
```
Use case:
* command: "GET"
  * netdev
  * version
  * vm
  * netstat
  * interface-stat
  * swap-memory
  * virtual-memory
  * cpuinfo
  * cputimestat
  * avgstat
  * misc
  * arp
  * modules

proc examples:
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netdev"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netstat", "property":"tcp"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/

sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cpuinfo"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cputimestat"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"avgstat"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"virtual-memory"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"swap-memory"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/misc
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/net/arp
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/modules
```

```
[sus@Zeus api-routerd]# curl --header "Content-Type: application/json" --request GET  --header "X-Session-Token: aaaaa" http://localhost:8080/proc/misc
{"130":"watchdog","144":"nvram","165":"vmmon","183":"hw_random","184":"microcode","227":"mcelog","228":"hpet","229":"fuse","231":"snapshot","232":"kvm","235":"autofs","236":"device-mapper","53":"vboxnetctl","54":"vsock","55":"vmci","56":"vboxdrvu","57":"vboxdrv","58":"rfkill","59":"memory_bandwidth","60":"network_throughput","61":"network_latency","62":"cpu_dma_latency","63":"vga_arbiter"}
```
proc vm: property any file name in /proc/sys/vm/
```
[sus@Zeus api-routerd]# curl --header "Content-Type: application/json" --request GET --data '{"property":"swappiness"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/vm
{"property":"swappiness","value":"70"}

```

```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cputimestat"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
[{"cpu":"cpu0","user":231.77,"system":103.54,"idle":11580.69,"nice":1.44,"iowait":6.68,"irq":30.09,"softirq":16.67,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu1","user":233.47,"system":108.38,"idle":11577.43,"nice":2.39,"iowait":5.17,"irq":32.93,"softirq":10.67,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu2","user":233.11,"system":106.65,"idle":11519.95,"nice":1.72,"iowait":5.52,"irq":82.22,"softirq":10.52,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu3","user":235.06,"system":109.29,"idle":11585.23,"nice":1.98,"iowait":6.6,"irq":24.98,"softirq":8.36,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu4","user":233.62,"system":100.02,"idle":11600.14,"nice":2.53,"iowait":6.13,"irq":21.95,"softirq":7.41,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu5","user":225.02,"system":101.52,"idle":11602.33,"nice":7.97,"iowait":7.27,"irq":21.61,"softirq":7.47,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu6","user":238.34,"system":98.43,"idle":11590.73,"nice":1.79,"iowait":6.45,"irq":26.88,"softirq":8.7,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu7","user":238.54,"system":97.5,"idle":11601.67,"nice":1.28,"iowait":6.34,"irq":21.17,"softirq":7.11,"steal":0,"guest":0,"guestNice":0,"stolen":0}]
```

More example
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netdev"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
[{"name":"wlan0","bytesSent":21729026,"bytesRecv":222301420,"packetsSent":127279,"packetsRecv":210421,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"tunl0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"lo","bytesSent":87414,"bytesRecv":87414,"packetsSent":927,"packetsRecv":927,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"eth0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"sit0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"my-macvlan","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"vmnet1","bytesSent":0,"bytesRecv":0,"packetsSent":701,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"vmnet8","bytesSent":0,"bytesRecv":0,"packetsSent":701,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0}]

```

```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
{"hostname":"Zeus","uptime":17747,"bootTime":1545381768,"procs":360,"os":"linux","platform":"fedora","platformFamily":"fedora","platformVersion":"29","kernelVersion":"4.19.2-300.fc29.x86_64","virtualizationSystem":"kvm","virtualizationRole":"host","hostid":"27f7c64c-3148-11b2-a85c-ec64a5733ce1"}

```
set and get any value from ```/proc/sys/net```.
supported: IPv4, IPv6 and core
```
[sus@Zeus api-routerd]# curl --header "Content-Type: application/json" --request GET --data '{"property":"forwarding", "link":"enp0s31f6", "path":"ipv4"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net
{"path":"ipv4","property":"forwarding","value":"0","link":"enp0s31f6"}
[sus@Zeus api-routerd]# curl --header "Content-Type: application/json" --request PUT --data '{"property":"forwarding", "value":"1","link":"enp0s31f6", "path":"ipv4"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net
[sus@Zeus api-routerd]# curl --header "Content-Type: application/json" --request GET --data '{"property":"forwarding", "link":"enp0s31f6", "path":"ipv4"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net
{"path":"ipv4","property":"forwarding","value":"1","link":"enp0s31f6"}
```

##### Use case configure link

Set address
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-address", "address":"192.168.1.131/24", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/address/add
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever

```
Set link up/down
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-up", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/set

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-down", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

Set MTU
 ```
[sus@Zeus api-router]$curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-mtu", "link":"dummy", "mtu":"1280"}' http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1280 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

Set Default GateWay
```
 sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-default-gw", "link":"dummy", "gateway":"192.168.1.1/24", "onlink":"true"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/route/add
[sus@Zeus api-router]$ ip route
default via 192.168.1.1 dev dummy onlink

```

Create a bridge and enslave two links
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-link-bridge", "link":"test-br", "enslave":["dummy","dummy1"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/add

[sus@Zeus log]# ip link
9: dummy: <BROADCAST,NOARP> mtu 12801 qdisc noop master test-br state DOWN mode DEFAULT group default qlen 1000
    link/ether f2:58:ea:f3:83:1e brd ff:ff:ff:ff:ff:ff
10: dummy1: <BROADCAST,NOARP> mtu 1500 qdisc noop master test-br state DOWN mode DEFAULT group default qlen 1000
    link/ether 12:00:9a:65:36:6d brd ff:ff:ff:ff:ff:ff
11: test-br: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default
    link/ether 12:00:9a:65:36:6d brd ff:ff:ff:ff:ff:ff
[sus@Zeus log]#

```

Delete a link
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request DELETE --data '{"action":"delete-link", "link":"test-br"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/delete
```

##### Use Case: networkd
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"eth0", "DHCP":"yes", "LLDP":"yes","Addresses": [{"Address":"192.168.1.2", "Label":"test1"},{"Address":"192.168.1.4", "Label":"test3", "Peer":"192.168.1.5"}], "Routes": [{"Gateway":"192.168.1.1",  "GatewayOnlink":"true"},{"Destination":"192.168.1.10","Table":"10"}]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/network

[sus@Zeus api-router]$ cat /var/run/systemd/network/25-eth0.network
[Match]
Name=eth0

[Network]
DHCP=yes
LLDP=yes


[Address]
Address=192.168.1.2
Label=test1

[Address]
Address=192.168.1.4
Peer=192.168.1.5
Label=test3


[Route]
Gateway=192.168.1.1
GatewayOnlink=yes

[Route]
Destination=192.168.1.10
Table=10

```

networkd NetDev
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"bond-test", "Description":"testing bond", "Kind":"bond", "Mode":"balance-rr"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/netdev

[sus@Zeus log]# cat /var/run/systemd/network/25-bond-test.netdev
[NetDev]
Name=bond-test
Description=testing bond
Kind=bond

[Bond]
Mode=balance-rr

```

Bridge
```

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"bridge-test", "Description":"testing bridge", "Kind":"bridge", "HelloTimeSec":"30"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/netdev
[sus@Zeus api-router]$ cat /var/run/systemd/network/25-bridge-test.netdev
[NetDev]
Name=bridge-test
Description=testing bridge
Kind=bridge

[Bridge]
HelloTimeSec =30

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"eth0", "Description":"etho bridge enslave", "Bridge":"bridge-test"}' http://localhost:8080/network/networkd/network
[sus@Zeus api-router]$ cat /var/run/systemd/network/25-eth0.network
[Match]
Name=eth0

[Network]
Bridge=bridge-test


```

Example: Get and Set Hostname
```
[sus@Zeus api-routerd]$  curl --header "Content-Type: application/json" --request GET --data '{"property":"static"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/
{"method":"","property":"StaticHostname","value":"Zeus"}

[sus@Zeus api-routerd]$ curl --header "Content-Type: application/json" --request PUT --data '{"method":"static", "value":"Zeus1"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/

[sus@Zeus api-routerd]$ curl --header "Content-Type: application/json" --request GET --data '{"property":"static"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/
{"method":"","property":"StaticHostname","value":"Zeus1"}
```


Example: Netlink
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"link":"wlan0"}' --header "X-Session-Token: aaaaa" http://192.168.225.23:8080/network/link/get
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"link":"wlan0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/get
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-hostname"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/get
http://localhost:8080/hostname
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
```

proc: netstat protocol tcp
```
[sus@Zeus api-router]$curl --header "Content-Type: application/json" --request GET --data '{"path":"netstat", "property":"tcp"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/
```

##### ethtool
```
[sus@Zeus src]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-features", "link":"eth0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/ethtool/get
[sus@Zeus src]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-bus", "link":"eth0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/ethtool/get
[sus@Zeus src]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-stat", "link":"eth0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/ethtool/get
```

example:
```
[sus@Zeus ethtool]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-features", "link":"eth0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/ethtool/get
{"esp-hw-offload":false,"esp-tx-csum-hw-offload":false,"fcoe-mtu":false,"highdma":true,"hw-tc-offload":false,"l2-fwd-offload":false,"loopback":false,"netns-local":false,"rx-all":false,"rx-checksum":true,"rx-fcs":false,"rx-gro":true,"rx-gro-hw":false,"rx-hashing":true,"rx-lro":false,"rx-ntuple-filter":false,"rx-udp_tunnel-port-offload":false,"rx-vlan-filter":false,"rx-vlan-hw-parse":true,"rx-vlan-stag-filter":false,"rx-vlan-stag-hw-parse":false,"tls-hw-record":false,"tls-hw-rx-offload":false,"tls-hw-tx-offload":false,"tx-checksum-fcoe-crc":false,"tx-checksum-ip-generic":true,"tx-checksum-ipv4":false,"tx-checksum-ipv6":false,"tx-checksum-sctp":false,"tx-esp-segmentation":false,"tx-fcoe-segmentation":false,"tx-generic-segmentation":true,"tx-gre-csum-segmentation":false,"tx-gre-segmentation":false,"tx-gso-partial":false,"tx-gso-robust":false,"tx-ipxip4-segmentation":false,"tx-ipxip6-segmentation":false,"tx-lockless":false,"tx-nocache-copy":false,"tx-scatter-gather":true,"tx-scatter-gather-fraglist":false,"tx-sctp-segmentation":false,"tx-tcp-ecn-segmentation":false,"tx-tcp-mangleid-segmentation":false,"tx-tcp-segmentation":true,"tx-tcp6-segmentation":true,"tx-udp-segmentation":false,"tx-udp_tnl-csum-segmentation":false,"tx-udp_tnl-segmentation":false,"tx-vlan-hw-insert":true,"tx-vlan-stag-hw-insert":false,"vlan-challenged":false}

[sus@Zeus ethtool]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-driver-name", "link":"eth0"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/ethtool/get
{"action":"get-link-driver-name","link":"eth0","reply":"e1000e"}

```
