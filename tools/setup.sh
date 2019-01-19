#!/bin/sh

set -eux -o pipefail

prefix="pivot"
router="$prefix-router"
net="$prefix-net"
veth="$prefix-veth"
host="$prefix-host"
#root="/tmp/$prefix"
i=0

teardown() {
	echo >&2 -ne "\nErrors detected, tearing down... "
	./teardown.sh
	echo >&2 "Done."
}

trap teardown ERR

while read -r op a b c d; do
	echo -e >&2 "$op $a\\t"
	case "$op" in
		R)
			# create router NS
			ip netns add "$router-$a"
			# enable IPv4 forwarding
			ip netns exec "$router-$a" sysctl net.ipv4.ip_forward=1 >/dev/null
			# setup loopback
			ip netns exec "$router-$a" ip addr add 127.0.0.1/8 dev lo
			ip netns exec "$router-$a" ip link set lo up
			;;
		N)
			# create network NS
			ip netns add "$net-$a"
			# create a bridge in the NS
			ip netns exec "$net-$a" ip link add br0 type bridge
			ip netns exec "$net-$a" ip link set br0 up
			# setup loopback
			ip netns exec "$net-$a" ip addr add 127.0.0.1/8 dev lo
			ip netns exec "$net-$a" ip link set lo up
			;;
		C)
			# create a veth pair connecting router $a with network $b, set both ends up
			ip link add "$veth-$i-0" type veth peer name "$veth-$i-1"
			ip link set "$veth-$i-0" netns "$router-$a"
			ip link set "$veth-$i-1" netns "$net-$b"
			ip netns exec "$router-$a" ip link set "$veth-$i-0" up
			ip netns exec "$net-$b" ip link set "$veth-$i-1" up
			# assign IPv4 address $c to the router end of the veth pair
			ip netns exec "$router-$a" ip addr add "$c" dev "$veth-$i-0"
			# add the other end to the network bridge
			ip netns exec "$net-$b" ip link set "$veth-$i-1" master br0
			i=$(($i + 1))
			;;
		H)
			# create host NS
			ip netns add "$host-$a"
			# create a veth pair connecting host $a with network $b, set both ends up
			ip link add "$veth-$i-0" type veth peer name "$veth-$i-1"
			ip link set "$veth-$i-0" netns "$host-$a"
			ip link set "$veth-$i-1" netns "$net-$b"
			ip netns exec "$host-$a" ip link set "$veth-$i-0" up
			ip netns exec "$net-$b" ip link set "$veth-$i-1" up
			# assign IPv4 address $c to the host end of the veth pair, setup default GW
			ip netns exec "$host-$a" ip addr add "$c" dev "$veth-$i-0"
			ip netns exec "$host-$a" ip route add default via "$d" dev "$veth-$i-0"
			# setup loopback
			ip netns exec "$host-$a" ip addr add 127.0.0.1/8 dev lo
			ip netns exec "$host-$a" ip link set lo up
			# add the other end to the network bridge
			ip netns exec "$net-$b" ip link set "$veth-$i-1" master br0
			i=$(($i + 1))
			;;
		B)
			# create BIRD config file
			cat <<-END >"/tmp/bird.$a.cfg"
				protocol kernel {
					ipv4 {
							export all;
					};
				}

				protocol device {}

				protocol ospf {
					ipv4 {
						import all;
						export all;
					};
					area 0.0.0.0 {
						interface "*";
					};
				}
			END
			# start BIRD in router's NS
			ip netns exec "$router-$a" bird -c "/tmp/bird.$a.cfg" -s "/tmp/bird.$a.ctl"
			;;
		Q)
			# create Quagga ospfd config file
			cat <<-END >"/tmp/ospfd.$a.cfg"
				password 1234
				interface *
					ip ospf area 0.0.0.0

				router ospf
					redistribute connected
					network 0.0.0.0/0 area 0.0.0.0
			END
			# start Quagga core Zebra daemon and OSPF daemon
			ip netns exec "$router-$a" zebra -d -i "/tmp/zebra.$a.pid" -f "/dev/null"
			ip netns exec "$router-$a" ospfd -d -i "/tmp/ospfd.$a.pid" -f "/tmp/ospfd.$a.cfg"
			# expose marvelous Quagga command-line interface (VTY) through UNIX-domain socket
			ip netns exec "$router-$a" socat UNIX-LISTEN:/tmp/ospfd.$a.ctl,fork TCP4-CONNECT:localhost:2604 > /dev/null 2>&1 &

	esac
done
