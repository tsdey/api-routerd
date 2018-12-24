![Logo](https://ibin.co/4R6Hzr2H7l4A.png)

#### A RestAPI GateWay For Linux

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

##### Basic Setup and Use of api-routerd

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

Install libs
```
$ go get github.com/sirupsen/logrus
$ go get github.com/gorilla/mux
$ go get github.com/vishvananda/netlink
$ go get github.com/shirou/gopsutil
$ go get github.com/coreos/go-systemd/dbus
$ go get github.com/godbus/dbus
$ go get github.com/safchain/ethtool
```

Now build it
```
[sus@Zeus src]$ git clone https://github.com/RestGW/api-routerd
[sus@Zeus src]$ cd api-routerd/cmd
[sus@Zeus cmd]$ pwd
/home/sus/go/src/api-routerd/cmd
[sus@Zeus cmd]$ go build -o api-routerd

```

Use case: systemd
```
[sus@Zeus] curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-property", "unit":"sshd.service", "property":"CPUShares", "value":"1024"}' http://localhost:8080/service/systemd/property
{"property":"CPUShares","value":"1024"}

[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"start","unit":"sshd.service"}' http://localhost:8080/service/systemd

[sus@Zeus]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"stop","unit":"sshd.service"}' http://localhost:8080/service/systemd

[sus@Zeus]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"status","unit":"sshd.service"}' http://localhost:8080/service/systemd
```
Use Case:
         command: "GET"
                       "netdev, version", "vm", "netstat", "interface-stat":
                       "swap-memory", "virtual-memory", "cpuinfo","cputimestat","avgstat":

```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netdev"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"vm"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netstat", "property":"udp"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cpuinfo"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cputimestat"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"avgstat"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"virtual-memory"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"swap-memory"}' http://localhost:8080/proc/
```


proc vm: property any file name in /proc/sys/vm/
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"vm", "property":"swappiness"}' http://localhost:8080/proc/
{"path":"","property":"swappiness","value":"60"}
```

```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"cputimestat"}' http://localhost:8080/proc/
[{"cpu":"cpu0","user":231.77,"system":103.54,"idle":11580.69,"nice":1.44,"iowait":6.68,"irq":30.09,"softirq":16.67,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu1","user":233.47,"system":108.38,"idle":11577.43,"nice":2.39,"iowait":5.17,"irq":32.93,"softirq":10.67,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu2","user":233.11,"system":106.65,"idle":11519.95,"nice":1.72,"iowait":5.52,"irq":82.22,"softirq":10.52,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu3","user":235.06,"system":109.29,"idle":11585.23,"nice":1.98,"iowait":6.6,"irq":24.98,"softirq":8.36,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu4","user":233.62,"system":100.02,"idle":11600.14,"nice":2.53,"iowait":6.13,"irq":21.95,"softirq":7.41,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu5","user":225.02,"system":101.52,"idle":11602.33,"nice":7.97,"iowait":7.27,"irq":21.61,"softirq":7.47,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu6","user":238.34,"system":98.43,"idle":11590.73,"nice":1.79,"iowait":6.45,"irq":26.88,"softirq":8.7,"steal":0,"guest":0,"guestNice":0,"stolen":0},{"cpu":"cpu7","user":238.54,"system":97.5,"idle":11601.67,"nice":1.28,"iowait":6.34,"irq":21.17,"softirq":7.11,"steal":0,"guest":0,"guestNice":0,"stolen":0}]
```

More example
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netdev"}' http://localhost:8080/proc/
[{"name":"wlan0","bytesSent":21729026,"bytesRecv":222301420,"packetsSent":127279,"packetsRecv":210421,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"tunl0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"lo","bytesSent":87414,"bytesRecv":87414,"packetsSent":927,"packetsRecv":927,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"eth0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"sit0","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"my-macvlan","bytesSent":0,"bytesRecv":0,"packetsSent":0,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"vmnet1","bytesSent":0,"bytesRecv":0,"packetsSent":701,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0},{"name":"vmnet8","bytesSent":0,"bytesRecv":0,"packetsSent":701,"packetsRecv":0,"errin":0,"errout":0,"dropin":0,"dropout":0,"fifoin":0,"fifoout":0}]

```

```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' http://localhost:8080/proc/
{"hostname":"Zeus","uptime":17747,"bootTime":1545381768,"procs":360,"os":"linux","platform":"fedora","platformFamily":"fedora","platformVersion":"29","kernelVersion":"4.19.2-300.fc29.x86_64","virtualizationSystem":"kvm","virtualizationRole":"host","hostid":"27f7c64c-3148-11b2-a85c-ec64a5733ce1"}

```

use case configure link

set address
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-address", "address":"192.168.1.131/24", "link":"dummy"}' http://localhost:8080/network/address/add
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noop state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever


```
set link up/down
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-up", "link":"dummy"}' http://localhost:8080/network/link/set

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-down", "link":"dummy"}' http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1500 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

set MTU
 ```
[sus@Zeus api-router]$curl --header "Content-Type: application/json" --request POST --data '{"action":"set-link-mtu", "link":"dummy", "mtu":"1280"}' http://localhost:8080/network/link/set
[sus@Zeus api-router]$ ip addr show dummy
9: dummy: <BROADCAST,NOARP> mtu 1280 qdisc noqueue state DOWN group default qlen 1000
    link/ether ea:d0:e3:be:ea:25 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.131/24 brd 192.168.1.255 scope global dummy
       valid_lft forever preferred_lft forever
```

Set GateWay
```
 sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-default-gw", "link":"dummy", "gateway":"192.168.1.1/24", "onlink":"true"}' http://localhost:8080/network/route/add
[sus@Zeus api-router]$ ip route
default via 192.168.1.1 dev dummy onlink

```

Create a bridge and enslave two links
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"add-link-bridge", "link":"test-br", "enslave":["dummy","dummy1"]}' http://localhost:8080/network/link/add

[root@Zeus log]# ip link
9: dummy: <BROADCAST,NOARP> mtu 12801 qdisc noop master test-br state DOWN mode DEFAULT group default qlen 1000
    link/ether f2:58:ea:f3:83:1e brd ff:ff:ff:ff:ff:ff
10: dummy1: <BROADCAST,NOARP> mtu 1500 qdisc noop master test-br state DOWN mode DEFAULT group default qlen 1000
    link/ether 12:00:9a:65:36:6d brd ff:ff:ff:ff:ff:ff
11: test-br: <BROADCAST,MULTICAST> mtu 1500 qdisc noop state DOWN mode DEFAULT group default
    link/ether 12:00:9a:65:36:6d brd ff:ff:ff:ff:ff:ff
[root@Zeus log]#

```

delete a link
```
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request DELETE --data '{"action":"delete-link", "link":"test-br"}' http://localhost:8080/network/link/delete
```

Use Case: networkd
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"eth0", "DHCP":"yes", "LLDP":"yes","Addresses": [{"Address":"192.168.1.2", "Label":"test1"},{"Address":"192.168.1.4", "Label":"test3", "Peer":"192.168.1.5"}], "Routes": [{"Gateway":"192.168.1.1",  "GatewayOnlink":"true"},{"Destination":"192.168.1.10","Table":"10"}]}' http://localhost:8080/network/networkd/network

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
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"bond-test", "Description":"testing bond", "Kind":"bond", "Mode":"balance-rr"}' http://localhost:8080/network/networkd/netdev

[root@Zeus log]# cat /var/run/systemd/network/25-bond-test.netdev
[NetDev]
Name=bond-test
Description=testing bond
Kind=bond

[Bond]
Mode=balance-rr

```

Bridge
```

[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request PUT --data '{"Name":"bridge-test", "Description":"testing bridge", "Kind":"bridge", "HelloTimeSec":"30"}' http://localhost:8080/network/networkd/netdev
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

Examples:
```
[sus@Zeus hostname]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-hostname", "property":"static"}' http://localhost:8080/hostname/
{"action":"get-hostname","method":"","property":"StaticHostname","value":"Zeus"}

[sus@Zeus hostname]$ curl --header "Content-Type: application/json" --request PUT --data '{"action":"set-hostname", "method":"static", "value":"zeus1"}' http://localhost:8080/hostname/
[sus@Zeus hostname]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-hostname", "property":"static"}' http://localhost:8080/hostname/
{"action":"get-hostname","method":"","property":"StaticHostname","value":"zeus1"}

sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' http://localhost:8080/proc/
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"link":"wlan0"}' http://192.168.225.23:8080/network/link/get
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"link":"wlan0"}' http://localhost:8080/network/link/get
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-hostname"}' http://localhost:8080/hostname/get
http://localhost:8080/hostname
sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"version"}' http://localhost:8080/proc/
```

proc: netstat protocol tcp
```
[sus@Zeus api-router]$ curl --header "Content-Type: application/json" --request GET --data '{"path":"netstat", "property":"tcp"}' http://localhost:8080/proc/

[{"fd":11,"family":2,"type":1,"localaddr":{"ip":"0.0.0.0","port":902},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1577},{"fd":21,"family":2,"type":1,"localaddr":{"ip":"127.0.0.1","port":8307},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1639},{"fd":4,"family":2,"type":1,"localaddr":{"ip":"0.0.0.0","port":22},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1030},{"fd":7,"family":2,"type":1,"localaddr":{"ip":"127.0.0.1","port":631},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1027},{"fd":15,"family":2,"type":1,"localaddr":{"ip":"0.0.0.0","port":443},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1639},{"fd":14,"family":2,"type":1,"localaddr":{"ip":"127.0.0.1","port":8125},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[970,970,970,970],"pid":1257},{"fd":3,"family":2,"type":1,"localaddr":{"ip":"0.0.0.0","port":19999},"remoteaddr":{"ip":"0.0.0.0","port":0},"status":"LISTEN","uids":[970,970,970,970],"pid":1257},{"fd":351,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49868},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":144,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":60314},"remoteaddr":{"ip":"172.217.166.67","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":135,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":57972},"remoteaddr":{"ip":"192.30.253.125","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":294,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":55414},"remoteaddr":{"ip":"172.217.167.170","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":249,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":34938},"remoteaddr":{"ip":"216.58.203.174","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":204,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49846},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":131,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":39134},"remoteaddr":{"ip":"192.30.253.124","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":218,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49848},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":210,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":36848},"remoteaddr":{"ip":"172.217.194.188","port":5228},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":176,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":41490},"remoteaddr":{"ip":"192.30.253.124","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":170,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49834},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":363,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49878},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":223,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49488},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":362,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49876},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":302,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":33452},"remoteaddr":{"ip":"172.217.160.163","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":186,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":48990},"remoteaddr":{"ip":"172.217.166.46","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":361,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49874},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":229,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":37740},"remoteaddr":{"ip":"185.199.108.154","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":356,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49864},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":352,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49870},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":360,"family":2,"type":1,"localaddr":{"ip":"192.168.225.101","port":49872},"remoteaddr":{"ip":"151.101.152.133","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":10,"family":10,"type":1,"localaddr":{"ip":"::","port":902},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1577},{"fd":5,"family":10,"type":1,"localaddr":{"ip":"::","port":8080},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":21553},{"fd":20,"family":10,"type":1,"localaddr":{"ip":"::1","port":8307},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1639},{"fd":5,"family":10,"type":1,"localaddr":{"ip":"::","port":22},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1030},{"fd":6,"family":10,"type":1,"localaddr":{"ip":"::1","port":631},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1027},{"fd":17,"family":10,"type":1,"localaddr":{"ip":"::","port":443},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[0,0,0,0],"pid":1639},{"fd":11,"family":10,"type":1,"localaddr":{"ip":"::1","port":8125},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[970,970,970,970],"pid":1257},{"fd":4,"family":10,"type":1,"localaddr":{"ip":"::","port":19999},"remoteaddr":{"ip":"::","port":0},"status":"LISTEN","uids":[970,970,970,970],"pid":1257},{"fd":6,"family":10,"type":1,"localaddr":{"ip":"::1","port":8080},"remoteaddr":{"ip":"::1","port":45618},"status":"ESTABLISHED","uids":[0,0,0,0],"pid":21553},{"fd":3,"family":10,"type":1,"localaddr":{"ip":"::1","port":45618},"remoteaddr":{"ip":"::1","port":8080},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":25067},{"fd":177,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":58948},"remoteaddr":{"ip":"2a03:2880:f02f:f:face:b00c:0:2","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":212,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":48310},"remoteaddr":{"ip":"2a02:26f0:7b:89b::25ea","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":197,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":58956},"remoteaddr":{"ip":"2a03:2880:f02f:f:face:b00c:0:2","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":169,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":45970},"remoteaddr":{"ip":"2a00:1450:400c:c0b::5e","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":262,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":34300},"remoteaddr":{"ip":"2600:1417:75:192::25ea","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783},{"fd":171,"family":10,"type":1,"localaddr":{"ip":"2409:4042:200b:c0ce:8556:56c3:7d23:6211","port":36256},"remoteaddr":{"ip":"2a03:2880:f12f:83:face:b00c:0:25de","port":443},"status":"ESTABLISHED","uids":[1000,1000,1000,1000],"pid":2783}]
```


ethtool
```
curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-features", "link":"eth0"}' http://localhost:8080/network/ethtool/get
curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-bus", "link":"eth0"}' http://localhost:8080/network/ethtool/get
curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-stat", "link":"eth0"}' http://localhost:8080/network/ethtool/get
```

example:
```
[sus@Zeus ethtool]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-features", "link":"eth0"}' http://localhost:8080/network/ethtool/get
{"esp-hw-offload":false,"esp-tx-csum-hw-offload":false,"fcoe-mtu":false,"highdma":true,"hw-tc-offload":false,"l2-fwd-offload":false,"loopback":false,"netns-local":false,"rx-all":false,"rx-checksum":true,"rx-fcs":false,"rx-gro":true,"rx-gro-hw":false,"rx-hashing":true,"rx-lro":false,"rx-ntuple-filter":false,"rx-udp_tunnel-port-offload":false,"rx-vlan-filter":false,"rx-vlan-hw-parse":true,"rx-vlan-stag-filter":false,"rx-vlan-stag-hw-parse":false,"tls-hw-record":false,"tls-hw-rx-offload":false,"tls-hw-tx-offload":false,"tx-checksum-fcoe-crc":false,"tx-checksum-ip-generic":true,"tx-checksum-ipv4":false,"tx-checksum-ipv6":false,"tx-checksum-sctp":false,"tx-esp-segmentation":false,"tx-fcoe-segmentation":false,"tx-generic-segmentation":true,"tx-gre-csum-segmentation":false,"tx-gre-segmentation":false,"tx-gso-partial":false,"tx-gso-robust":false,"tx-ipxip4-segmentation":false,"tx-ipxip6-segmentation":false,"tx-lockless":false,"tx-nocache-copy":false,"tx-scatter-gather":true,"tx-scatter-gather-fraglist":false,"tx-sctp-segmentation":false,"tx-tcp-ecn-segmentation":false,"tx-tcp-mangleid-segmentation":false,"tx-tcp-segmentation":true,"tx-tcp6-segmentation":true,"tx-udp-segmentation":false,"tx-udp_tnl-csum-segmentation":false,"tx-udp_tnl-segmentation":false,"tx-vlan-hw-insert":true,"tx-vlan-stag-hw-insert":false,"vlan-challenged":false}

[sus@Zeus ethtool]$ curl --header "Content-Type: application/json" --request GET --data '{"action":"get-link-driver-name", "link":"eth0"}' http://localhost:8080/network/ethtool/get
{"action":"get-link-driver-name","link":"eth0","reply":"e1000e"}

```
