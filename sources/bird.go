package sources

import (
	"bufio"
	"fmt"
	"github.com/Hrubon/pivot/model"
	"github.com/pkg/errors"
	"net"
	"regexp"
	"strconv"
)

const (
	dumpCmd = "show route"
)

type BirdSource struct {
	Source
	ctlPath string
}

func parseRoute(l string) (*model.Route, error) {
	re := regexp.MustCompile(`[\s\[\]]+`)
	parts := re.Split(l, -1)
	_, ipNet, err := net.ParseCIDR(parts[0])
	if err != nil {
		return nil, err
	}
	rid := parts[len(parts)-2]
	return &model.Route{
		Network: ipNet,
		RouterID: rid,
	}, nil
}

func (s *BirdSource) GetRoutes() (*model.RouteList, error) {
	c, err := net.Dial("unix", s.ctlPath)
	if err != nil {
		return nil, errors.Wrap(err, "cannot connect to BIRD control socket")
	}
	defer c.Close()
	rlist := model.NewRouteList()
	fmt.Fprintf(c, "%s\n", dumpCmd)
	scanner := bufio.NewScanner(c)
	scanner.Scan() // skip banner
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
		if len(l) == 0 {
			continue // regular empty line
		}
		if l[0] >= byte('0') && l[0] <= byte('9') {
			r, err := parseRoute(l)
			if err != nil {
				return nil, errors.Wrap(err, "error while parsing route")
			}
			rlist.Routes = append(rlist.Routes, r)
		}
	}
	return nil, scanner.Err()
}

func NewBirdSource(ctlPath string) *BirdSource {
	return &BirdSource{
		ctlPath: ctlPath,
	}
}
