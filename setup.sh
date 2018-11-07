#!/bin/sh -eu

prefix="pivot"
router="$prefix-router"
net="$prefix-net"
veth="$prefix-veth"
host="$prefix-host"
i=0

teardown() {
	echo >&2 -ne "\nErrors detected, tearing down... "
	./teardown.sh
	echo >&2 "Done."
}

trap teardown ERR

while read -r op a b c d; do
	stdbuf -e0 echo -ne >&2 "$op $a\\t"
	case "$op" in
		R)
			# create router NS
			ip netns add "$router-$a"
			# enable IPv4 forwarding
			ip netns exec "$router-$a" sysctl net.ipv4.ip_forward=1 >/dev/null
			;;
		N)
			# create network NS
			ip netns add "$net-$a"
			# create a bridge in the NS
			ip netns exec "$net-$a" ip link add br0 type bridge
			ip netns exec "$net-$a" ip link set br0 up
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
			# add the other end to the network bridge
			ip netns exec "$net-$b" ip link set "$veth-$i-1" master br0
			i=$(($i + 1))
	esac
	echo >&2 -e "\u2713"
done
