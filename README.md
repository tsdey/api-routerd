![Logo](https://ibin.co/4R6Hzr2H7l4A.png)

#### A RestAPI Microservice GateWay For Linux

A super light weight remote management tool which uses REST API for real time configuration and performance as well as health monitoring for systems (containers) and applications. It provides fast API based monitoring without affecting the system it's running on.


#### Objectives:
- No client installation required. curl is enough.
- No GUI
- Minimal data transfer using JSON


##### Allows you to configure
- systemd
    - systemd informations
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
- configure nameserver /etc/resolv.conf

 ### api-routerd json API
 Refer spreadsheet [API](https://docs.google.com/spreadsheets/d/e/2PACX-1vTl2Vmp-BdTE5Vgi_PiW-qKPJnbLxdSso9kT2GAkAxCu_iWrw3_PZLlEuyXz0lbFgd7DoofXlmmb3dP/pubhtml
)


#### Quick Start

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

Now start using https
```
[sus@Zeus tls]$ curl --header "X-Session-Token: aaaaa" --request GET https://localhost:8080/network/ethtool/vmnet8/get-link-features -k --tlsv1.2

{"esp-hw-offload":false,"esp-tx-csum-hw-offload":false,"fcoe-mtu":false,"highdma":false,"hw-tc-offload":false,"l2-fwd-offload":false,"loopback":false,"netns-local":false,"rx-all":false,"rx-checksum":false,"rx-fcs":false,"rx-gro":true,"rx-gro-hw":false,"rx-hashing":false,"rx-lro":false,"rx-ntuple-filter":false,"rx-udp_tunnel-port-offload":false,"rx-vlan-filter":false,"rx-vlan-hw-parse":false,"rx-vlan-stag-filter":false,"rx-vlan-stag-hw-parse":false,"tls-hw-record":false,"tls-hw-rx-offload":false,"tls-hw-tx-offload":false,"tx-checksum-fcoe-crc":false,"tx-checksum-ip-generic":false,"tx-checksum-ipv4":false,"tx-checksum-ipv6":false,"tx-checksum-sctp":false,"tx-esp-segmentation":false,"tx-fcoe-segmentation":false,"tx-generic-segmentation":false,"tx-gre-csum-segmentation":false,"tx-gre-segmentation":false,"tx-gso-partial":false,"tx-gso-robust":false,"tx-ipxip4-segmentation":false,"tx-ipxip6-segmentation":false,"tx-lockless":false,"tx-nocache-copy":false,"tx-scatter-gather":false,"tx-scatter-gather-fraglist":false,"tx-sctp-segmentation":false,"tx-tcp-ecn-segmentation":false,"tx-tcp-mangleid-segmentation":false,"tx-tcp-segmentation":false,"tx-tcp6-segmentation":false,"tx-udp-segmentation":false,"tx-udp_tnl-csum-segmentation":false,"tx-udp_tnl-segmentation":false,"tx-vlan-hw-insert":false,"tx-vlan-stag-hw-insert":false,"vlan-challenged":false}
```

```
[sus@Zeus cmd]$ curl --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/misc --tlsv1.2 -k
{"130":"watchdog","144":"nvram","165":"vmmon","183":"hw_random","184":"microcode","227":"mcelog","228":"hpet","229":"fuse","231":"snapshot","232":"kvm","235":"autofs","236":"device-mapper","53":"vboxnetctl","54":"vsock","55":"vmci","56":"vboxdrvu","57":"vboxdrv","58":"rfkill","59":"memory_bandwidth","60":"network_throughput","61":"network_latency","62":"cpu_dma_latency","63":"vga_arbiter"}

[sus@Zeus cmd]$ curl --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/net/arp --tlsv1.2 -k
[{"ip_address":"192.168.225.1","hw_type":"0x1","flags":"0x2","hw_address":"1a:89:20:38:68:8f","mask":"*","device":"wlp4s0"}]

[sus@Zeus cmd]$ curl --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/modules --tlsv1.2 -k
```
Use case: systemd

Know systemd information
```
http://localhost:8080/service/systemd/version
http://localhost:8080/service/systemd/state
http://localhost:8080/service/systemd/features
http://localhost:8080/service/systemd/virtualization
http://localhost:8080/service/systemd/architecture
````
Get all units
```
http://localhost:8080/service/systemd/units
```
Get all properties of a unit
```
http://localhost:8080/service/systemd/sshd.service/get
```
Status of a unit
```
http://localhost:8080/service/systemd/sshd.service/status
```
Set and get propetties
```
[sus@Zeus] curl --request GET --header "X-Session-Token: aaaaa"http://localhost:8080/service/systemd/sshd.service/get/LimitNOFILESoft
[sus@Zeus]  curl --request GET --header "X-Session-Token: aaaaa"  http://localhost:8080/service/systemd/sshd.service/get/LimitNOFILE
[sus@Zeus]  curl --request GET --header "X-Session-Token: aaaaa"  http://localhost:8080/service/systemd/sshd.service/get
[sus@Zeus] curl --header "X-Session-Token: aaaaa" --header "Content-Type: application/json" --request PUT --data '{"value":"1100"}' http://localhost:8080/service/systemd/sshd.service/set/CPUShares
[sus@Zeus]  curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd/sshd.service/get/CPUShares
[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"start","unit":"sshd.service"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd
[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"stop","unit":"sshd.service"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd
[sus@Zeus proc]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd/sshd.service/status
{"property":"active","unit":"sshd.service"}
```

Send a signal to the service
```
$curl --header "Content-Type: application/json" --request POST --data '{"action":"kill","unit":"sshd.service", "value":"9"}' --header "X-Session-Token: aaaaa" http://localhost:8080/service/systemd
```
Use case:
* command: "GET"
  * netdev
  * version
  * vm
  * netstat
  * proto-counter-stat
  * proto-pid-stat
  * interface-stat
  * swap-memory
  * virtual-memory
  * cpuinfo
  * cputimestat
  * avgstat
  * misc
  * arp
  * modules
  * userstat
  * temperaturestat

/proc examples:
```
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/netdev
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/version
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/netstat/tcp
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/proto-pid-stat/2881/tcp

sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/cpuinfo
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/cputimestat
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/virtual-memory
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/swap-memory
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/userstat
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/temperaturestat
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/proto-counter-stat
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/misc
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" https://localhost:8080/proc/net/arp
sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/modules
```

information by pid request = "GET"
```
[sus@Zeus cmd]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/process/2769/pid-connections/
[{"fd":270,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":45354},"remoteaddr":{"ip":"74.125.24.102","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2769},{"fd":196,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49138},"remoteaddr":{"ip":"172.217.194.94","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":27

http://localhost:8080/proc/process/2769/pid-rlimit/
http://localhost:8080/proc/process/2769/pid-rlimit-usage/
http://localhost:8080/proc/process/2769/pid-status/
http://localhost:8080/proc/process/2769/pid-username/
http://localhost:8080/proc/process/2769/pid-open-files/
http://localhost:8080/proc/process/2769/pid-fds/
http://localhost:8080/proc/process/2769/pid-name/
http://localhost:8080/proc/process/2769/pid-memory-percent/
http://localhost:8080/proc/process/2769/pid-memory-maps/
http://localhost:8080/proc/process/2769/pid-memory-info/
http://localhost:8080/proc/process/2769/pid-io-counters/
```

```
[sus@Zeus api-routerd]# curl --header --request GET  --header "X-Session-Token: aaaaa" http://localhost:8080/proc/misc
{"130":"watchdog","144":"nvram","165":"vmmon","183":"hw_random","184":"microcode","227":"mcelog","228":"hpet","229":"fuse","231":"snapshot","232":"kvm","235":"autofs","236":"device-mapper","53":"vboxnetctl","54":"vsock","55":"vmci","56":"vboxdrvu","57":"vboxdrv","58":"rfkill","59":"memory_bandwidth","60":"network_throughput","61":"network_latency","62":"cpu_dma_latency","63":"vga_arbiter"}
```
proc vm: property any file name in /proc/sys/vm/
```
[sus@Zeus api-routerd]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/vm/swappiness
{"property":"swappiness","value":"60"}
[sus@Zeus api-routerd]$ curl --header "Content-Type: application/json" --request PUT --data '{"value":"70"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/vm/swappiness
{"property":"swappiness","value":"70"}

```

```
[sus@Zeus api-router]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/version
{"hostname":"Zeus","uptime":17747,"bootTime":1545381768,"procs":360,"os":"linux","platform":"fedora","platformFamily":"fedora","platformVersion":"29","kernelVersion":"4.19.2-300.fc29.x86_64","virtualizationSystem":"kvm","virtualizationRole":"host","hostid":"27f7c64c-3148-11b2-a85c-ec64a5733ce1"}

```
set and get any value from ```/proc/sys/net```.
supported: IPv4, IPv6 and core
```
[sus@Zeus proc]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net/ipv4/enp0s31f6/forwarding
{"path":"ipv4","property":"forwarding","value":"0","link":"enp0s31f6"}
[sus@Zeus proc]$  curl --header "Content-Type: application/json" --request PUT --data '{"value":"1"}' --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net/ipv4/enp0s31f6/forwarding
[sus@Zeus proc]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/proc/sys/net/ipv4/enp0s31f6/forwarding
{"path":"ipv4","property":"forwarding","value":"1","link":"enp0s31f6"}
[sus@Zeus proc]$

```

##### Use case configure link

Add address
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"add-address", "address":"192.168.1.131/24", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/address/add
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever

```
Set link up/down
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-link-up", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/set

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-link-down", "link":"dummy"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

Set MTU
 ```
[sus@Zeus api-router]$curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-link-mtu", "link":"dummy", "mtu":"1280"}' http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1280 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

Set Default GateWay
```
 sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"add-default-gw", "link":"dummy", "gateway":"192.168.1.1/24", "onlink":"true"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/route/add
[sus@Zeus api-router]$ ip route
default via 192.168.1.1 dev dummy onlink

```

Create a bridge and enslave two links
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"add-link-bridge", "link":"test-br", "enslave":["dummy","dummy1"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/link/add

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
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"Name":"eth0", "DHCP":"yes", "LLDP":"yes","Addresses": [{"Address":"192.168.1.2", "Label":"test1"},{"Address":"192.168.1.4", "Label":"test3", "Peer":"192.168.1.5"}], "Routes": [{"Gateway":"192.168.1.1",  "GatewayOnlink":"true"},{"Destination":"192.168.1.10","Table":"10"}]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/network

[sus@Zeus api-router]$ cat /etc/systemd/network/25-eth0.network
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
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"Name":"bond-test", "Description":"testing bond", "Kind":"bond", "Mode":"balance-rr"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/netdev

[sus@Zeus log]# cat /etc/systemd/network/25-bond-test.netdev
[NetDev]
Name=bond-test
Description=testing bond
Kind=bond

[Bond]
Mode=balance-rr

```

networkd Link
```
sus@Zeus api-routerd]$  curl --header "Content-Type: application/json" --request POST --data '{"MAC":"00:50:56:c0:00:08", "Name":"test","Description":"testing link", "WakeOnLan":"phy", "TCPSegmentationOffload":"yes"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/link

[sus@Zeus network]# pwd
/etc/systemd/network
[sus@Zeus network]# cat 00-test.link
[Match]
MACAddress=00:50:56:c0:00:08

[Link]
Description=testing link
Name=test
WakeOnLan=phy
TCPSegmentationOffload=yes
```

Bridge
```

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"Name":"bridge-test", "Description":"testing bridge", "Kind":"bridge", "HelloTimeSec":"30"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/netdev
[sus@Zeus api-router]$ cat /etc/systemd/network/25-bridge-test.netdev
[NetDev]
Name=bridge-test
Description=testing bridge
Kind=bridge

[Bridge]
HelloTimeSec =30

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"Name":"eth0", "Description":"etho bridge enslave", "Bridge":"bridge-test"}' http://localhost:8080/network/networkd/network
[sus@Zeus api-router]$ cat /etc/systemd/network/25-eth0.network
[Match]
Name=eth0

[Network]
Bridge=bridge-test


Tunnel
[sus@Zeus api-routerd]$ curl --header "Content-Type: application/json" --request POST --data '{"Name":"tunnel-test", "Description":"testing tunnel", "Kind":"tunnel", "Local":"192.168.1.2", "Remote":"192.168.1.2", "Independent":"true"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/networkd/netdev
[sus@Zeus api-routerd]$ cat /etc/systemd/network/
00-.link               00-test.link           25-bond-test.netdev    25-tunnel-test.netdev
[sus@Zeus api-routerd]$ cat /etc/systemd/network/25-tunnel-test.netdev
[NetDev]
Name=tunnel-test
Description=testing tunnel
Kind=tunnel

[Tunnel]
Local=192.168.1.2
Remote=192.168.1.2
Independent=true
```

##### Hostname

Example: Get and Set Hostname
```
[sus@Zeus ~]$ curl --request GET --header "X-Session-Token: aaaaa" http://localhost:8080/hostname
{"Chassis":"laptop","Deployment":"","HomeURL":"https://fedoraproject.org/","Hostname":"Zeus","IconName":"computer-laptop","KernelName":"Linux","KernelRelease":"4.19.2-300.fc29.x86_64","KernelVersion":"#1 SMP Wed Nov 14 19:05:24 UTC 2018","Location":"","OperatingSystemCPEName":"cpe:/o:fedoraproject:fedora:29","OperatingSystemPrettyName":"Fedora 29 (Twenty Nine)","PrettyHostname":"","StaticHostname":"Zeus"}

[sus@Zeus ~]$ curl --header "Content-Type: application/json" --request PUT --data '{"property":"SetStaticHostname","value":"Zeus1"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/set
[sus@Zeus ~]$ hostname
Zeus1
[sus@Zeus ~]$ curl --header "Content-Type: application/json" --request PUT --data '{"property":"SetStaticHostname","value":"Zeus"}' --header "X-Session-Token: aaaaa" http://localhost:8080/hostname/set
[sus@Zeus ~]$ hostname
Zeus
```

Supported Property for hostname
```
        "Hostname"
        "StaticHostname"
        "PrettyHostname"
        "IconName"
        "Chassis"
        "Deployment"
        "Location"
        "KernelName"
        "KernelRelease"
        "KernelVersion",
        "OperatingSystemPrettyName"
        "OperatingSystemCPEName"
        "HomeURL"
}
```

Supported Property (Methods) for setting hostname. For example: ```'{"property":"SetStaticHostname","value":"Zeus"}'```
```
        "SetHostname"
        "SetStaticHostname"
        "SetPrettyHostname"
        "SetIconName"
        "SetChassis"
        "SetDeployment"
        "SetLocation"
```

##### ethtool

```
[sus@Zeus src]$ curl --header "X-Session-Token: aaaaa" --request GET  http://localhost:8080/network/ethtool/vmnet8/get-link-features

{"esp-hw-offload":false,"esp-tx-csum-hw-offload":false,"fcoe-mtu":false,"highdma":false,"hw-tc-offload":false,"l2-fwd-offload":false,"loopback":false,"netns-local":false,"rx-all":false,"rx-checksum":false,"rx-fcs":false,"rx-gro":true,"rx-gro-hw":false,"rx-hashing":false,"rx-lro":false,"rx-ntuple-filter":false,"rx-udp_tunnel-port-offload":false,"rx-vlan-filter":false,"rx-vlan-hw-parse":false,"rx-vlan-stag-filter":false,"rx-vlan-stag-hw-parse":false,"tls-hw-record":false,"tls-hw-rx-offload":false,"tls-hw-tx-offload":false,"tx-checksum-fcoe-crc":false,"tx-checksum-ip-generic":false,"tx-checksum-ipv4":false,"tx-checksum-ipv6":false,"tx-checksum-sctp":false,"tx-esp-segmentation":false,"tx-fcoe-segmentation":false,"tx-generic-segmentation":false,"tx-gre-csum-segmentation":false,"tx-gre-segmentation":false,"tx-gso-partial":false,"tx-gso-robust":false,"tx-ipxip4-segmentation":false,"tx-ipxip6-segmentation":false,"tx-lockless":false,"tx-nocache-copy":false,"tx-scatter-gather":false,"tx-scatter-gather-fraglist":false,"tx-sctp-segmentation":false,"tx-tcp-ecn-segmentation":false,"tx-tcp-mangleid-segmentation":false,"tx-tcp-segmentation":false,"tx-tcp6-segmentation":false,"tx-udp-segmentation":false,"tx-udp_tnl-csum-segmentation":false,"tx-udp_tnl-segmentation":false,"tx-vlan-hw-insert":false,"tx-vlan-stag-hw-insert":false,"vlan-challenged":false}
```

```
[sus@Zeus src]$ curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/ethtool/wlp4s0/get-link-stat
{"ch_time":18446744073709551615,"ch_time_busy":18446744073709551615,"ch_time_ext_busy":18446744073709551615,"ch_time_rx":18446744073709551615,"ch_time_tx":18446744073709551615,"channel":0,"noise":18446744073709551615,"rx_bytes":1387313,"rx_dropped":45,"rx_duplicates":0,"rx_fragments":3255,"rx_packets":3275,"rxrate":117000000,"signal":227,"sta_state":4,"tx_bytes":584843,"tx_filtered":0,"tx_packets":2949,"tx_retries":321,"tx_retry_failed":0,"txrate":144400000}
```

```
[sus@Zeus cmd]$  curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/ethtool/wlp4s0/get-link-stat
{"ch_time":18446744073709551615,"ch_time_busy":18446744073709551615,"ch_time_ext_busy":18446744073709551615,"ch_time_rx":18446744073709551615,"ch_time_tx":18446744073709551615,"channel":0,"noise":18446744073709551615,"rx_bytes":1387313,"rx_dropped":45,"rx_duplicates":0,"rx_fragments":3255,"rx_packets":3275,"rxrate":117000000,"signal":227,"sta_state":4,"tx_bytes":584843,"tx_filtered":0,"tx_packets":2949,"tx_retries":321,"tx_retry_failed":0,"txrate":144400000}
```

```
[sus@Zeus cmd]$  curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/ethtool/wlp4s0/get-link-driver-name
{"action":"get-link-driver-name","link":"wlp4s0","reply":"iwlwifi"}
```

Get link netlink
```
[sus@Zeus ethtool]$ curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/link/get/wlp4s0
{"index":5,"MTU":1500,"TxQLen":0,"Name":"wlp4s0","HardwareAdd":"7c:76:35:ea:89:90","LinkOperState":""}

[sus@Zeus cmd]$  curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/address/get/wlp4s0
[{"action":"","link":"wlp4s0","address":"192.168.43.105/24","label":""},{"action":"","link":"wlp4s0","address":"2409:4042:239c:7f9d:e45f:27a9:c6de:c39e/64","label":""},{"action":"","link":"wlp4s0","address":"fe80::c912:39ce:e9a3:aaca/64","label":""}]
```

get all links
```
curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/link/get
```

get all routes
```
curl --header "X-Session-Token: aaaaa" --request GET http://localhost:8080/network/route/get
```
Replace route
```
[sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"replace-default-gw", "link":"dummy", "gateway":"192.168.1.3/24", "onlink":"true"}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/route/add
[sus@Zeus cmd]$ ip route
default via 192.168.1.3 dev dummy onlink
default via 192.168.225.1 dev wlp4s0 proto dhcp metric 600
172.16.130.0/24 dev vmnet8 proto kernel scope link src 172.16.130.1
192.168.1.0/24 dev dummy proto kernel scope link src 192.168.1.131
192.168.137.0/24 dev vmnet1 proto kernel scope link src 192.168.137.1
192.168.225.0/24 dev wlp4s0 proto kernel scope link src 192.168.225.101 metric 600
```

Get all addresses
```
http://localhost:8080/network/address/get
```

##### Add/Read/Delete confs in resolv.conf

Get
```
sus@Zeus cmd]$ http://localhost:8080/network/resolv
```

Add
```
sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request POST --data '{"servers":["192.168.1.131","192.168.1.132"], "search":["hello","hello2"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/resolv/add
```
Delete
```
sus@Zeus cmd]$ curl --header "Content-Type: application/json" --request DELETE --data '{"servers":["192.168.225.3","192.168.225.2"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/resolv/delete

```


Configure systemd-resolved
```
To get
http://localhost:8080/network/systemdresolved

Add
[sus@Zeus ~]$  curl --header "Content-Type: application/json" --request POST --data '{"dns":["192.168.1.131","192.168.1.132"], "fallback_dns":["hello","hello2"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/systemdresolved/add
{"dns":["10.68.5.26 10.64.63.6 192.168.225.1","192.168.1.131","192.168.1.132"],"fallback_dns":["8.8.8.8 8.8.4.4 2001:4860:4860::8888 2001:4860:4860::8844","hello","hello2"]}

Delete
[sus@Zeus ~]$  curl --header "Content-Type: application/json" --request DELETE --data '{"dns":["192.168.1.131"]}' --header "X-Session-Token: aaaaa" http://localhost:8080/network/systemdresolved/delete
{"dns":["10.68.5.26","10.64.63.6","192.168.225.1"],"fallback_dns":["8.8.8.8","8.8.4.4","2001:4860:4860::8888","2001:4860:4860::8844"]}

```
