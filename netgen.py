#!/usr/bin/env python3

import random
import ipaddress
import sys

def rand_subnet():
    addr = random.randrange(1, pow(2, 32))
    mask = random.randint(9, 24)
    mask_c = 32 - mask
    return ((addr >> mask_c) << mask_c, mask) # subtract mask

def nth_ip(n, subnet, mask):
    assert n <= 254
    return (str(ipaddress.ip_address(subnet + 1 + n)) + "/{}").format(mask)

def connect(router, network, ip):
    print("C", router, network, ip)

def gen(nrouters, nnets):
    # check input parameters
    if nrouters < 1 or nrouters > 254:
        raise ValueError("nrouters must be between 1 and 254")
    if nnets < 1:
        raise ValueError("nnets must be at least 1")

    # setup and print routers
    routers = range(nrouters)
    for i in routers:
        print("R", i + 1)

    # setup and print networks
    nets = range(nnets)
    ncnets = list(range(nnets))
    ipnets = []
    for i in nets:
        print("N", i + 1)
        ipnets.append((rand_subnet(), 0))
    
    # generate connections
    last_net = -1
    for i in range(nrouters):
        k = 1 + random.randrange(nnets)
        rnets = random.sample(nets, k)
        if last_net not in rnets and last_net != -1:
            rnets += [last_net]
        last_net = rnets[k - 1]
        for net in rnets:
            if net in ncnets:
                ncnets.remove(net)
            (netaddr, mask), n = ipnets[net]
            connect(i + 1, net + 1, nth_ip(n, netaddr, mask))
            ipnets[net] = ((netaddr, mask), n + 1)

    # ensure all nets are connected
    for net in ncnets:
        r = random.randrange(nrouters)
        (netaddr, mask), n = ipnets[net]
        connect(r + 1, net + 1, nth_ip(n, netaddr, mask))
        ipnets[net] = ((netaddr, mask), n + 1)

def usage(retval):
    sys.stderr.write("Usage: {} nrouters nnets\n".format(sys.argv[0]))
    exit(retval)

if __name__ == "__main__":
    if len(sys.argv) != 3:
        usage(1)
    try:
        nrouters = int(sys.argv[1])
        nnets = int(sys.argv[2])
        gen(nrouters, nnets)
    except ValueError as e:
        sys.stderr.write("{}: error: {}\n".format(sys.argv[0], e))
        exit(1)
