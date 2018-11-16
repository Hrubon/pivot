# The PIVOT Interactive Visualizer of OSPF Topology

## TOOLS

* `netgen.py` generates a file describing a random network in a format understood
by other tools.
* `setup.sh` reads network description from standard input (in the same format that
`netgen.py` generates) and sets up a virtual network with the same topology.
* `teardown.sh` cleans up whatever network objects `setup.sh` had generated before.

## REQUIREMENTS

These packages are necessary in order to run smoothly all of the scripts and binaries:

* Python 3.X
* BIRD 2.X
* Quagga 1.X
* `socat` utility
