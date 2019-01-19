package sources

import (
	"bufio"
	"fmt"
	"github.com/Hrubon/pivot/model"
	"github.com/pkg/errors"
	"net"
	"strings"
	"strconv"
)

const (
	dumpCmd = "show ospf state"
)

type BIRDSource struct {
	Source
	ctlPath string
}

func (s *BIRDSource) GetRoutes() (*model.RouteList, error) {
	c, err := net.Dial("unix", s.ctlPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to BIRD control socket")
	}
	defer c.Close()
	rlist := model.NewRouteList()
	fmt.Fprintf(c, "%s\n", dumpCmd)
	scanner := bufio.NewScanner(c)
	scanner.Scan() // skip banner
	var routerID string
	for scanner.Scan() {
		l := scanner.Text()
		if len(l) == 0 {
			return nil, errors.Errorf("BIRD CLI sent empty line (w/o status code)")
		}
		if l[0] != ' ' {
			if len(l) < 5 {
				return nil, errors.Errorf("BIRD CLI sent too short a line: '%s'", l)
			}
			status, err := strconv.Atoi(l[0:4])
			if err != nil {
				return nil, errors.Wrap(err, "BIRD CLI sent malformed status code")
			}
			if status < 1000 {
				return rlist, nil // success
			}
			if status >= 8000 {
				return nil, errors.Errorf("BIRD CLI replied: %s", l)
			}
			l = l[5:]
		} else {
			l = l[1:]
		}
		t := strings.TrimSpace(l)
		if len(t) == 0 {
			routerID = ""
			continue // regular empty line
		}
		p := strings.Split(t, " ")
		switch (p[0]) {
		case "router":
			if len(p) < 2 {
				return nil, errors.New("router section missing router ID")
			}
			routerID = p[1]
			continue
		case "stubnet", "network":
			if routerID == "" {
				continue
			}
			if len(p) < 4 || p[2] != "metric" {
				fmt.Println("L:", t)
				return nil, errors.New("malformed 'network' line in router section")
			}
			_, ipNet, err := net.ParseCIDR(p[1])
			if err != nil {
				return nil, errors.Wrap(err, "malformed network address")
			}
			metric, err := strconv.Atoi(p[3])
			if err != nil {
				return nil, errors.Wrap(err, "malformed metric")
			}
			rlist.Routes = append(rlist.Routes, &model.Route{
				RouterID: routerID,
				Network: ipNet,
				Metric: metric,
			})
		default:
			fmt.Println("Skipped:", t)
		}
	}
	return nil, scanner.Err()
}

func NewBIRDSource(ctlPath string) *BIRDSource {
	return &BIRDSource{
		ctlPath: ctlPath,
	}
}
