package pf

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"text/scanner"
)

// #include <net/if.h>
// #include <net/pfvar.h>
import "C"

// ParseSource sets the source ip (inet and inet6) based on the
// passed strings, if parsing failes err is returned
func (r *Rule) ParseSource(src, port string, neg bool) error {
	err := parsePort(&r.wrap.rule.src, port)
	if err != nil {
		return err
	}

	err = parseAddress(&r.wrap.rule.src, src, neg)
	if err != nil {
		return err
	}

	// determine if it is IPv6 or IPv4
	if strings.ContainsRune(src, ':') {
		r.SetAddressFamily(AddressFamilyInet6)
	} else {
		r.SetAddressFamily(AddressFamilyInet)
	}

	return nil
}

// ParseDestination sets the destination (inet and inet6) based on
// the passed strings, if parsing failes err returned
func (r *Rule) ParseDestination(dst, port string, neg bool) error {
	err := parsePort(&r.wrap.rule.dst, port)
	if err != nil {
		return err
	}

	return parseAddress(&r.wrap.rule.dst, dst, neg)
}

// parseAddress parses the passed string into the addr structure
func parseAddress(addr *C.struct_pf_rule_addr, address string, negative bool) error {
	if negative {
		addr.neg = 1
	}

	if strings.ContainsRune(address, '/') {
		_, n, err := net.ParseCIDR(address)
		if err != nil {
			return err
		}
		if strings.ContainsRune(address, ':') {
			copy(addr.addr.v[0:16], n.IP)
			copy(addr.addr.v[16:32], n.Mask)
		} else {
			copy(addr.addr.v[0:4], n.IP.To4())
			copy(addr.addr.v[16:20], n.Mask)
		}
	} else {
		n := net.ParseIP(address)
		if strings.ContainsRune(address, ':') {
			copy(addr.addr.v[0:16], n)
			copy(addr.addr.v[16:32], net.CIDRMask(128, 128))
		} else {
			copy(addr.addr.v[0:4], n.To4())
			copy(addr.addr.v[16:20], net.CIDRMask(32, 32))
		}
	}

	return nil
}

// parsePort parses the passed port into the address structure port section
func parsePort(addr *C.struct_pf_rule_addr, port string) error {
	s := scanner.Scanner{}
	s.Init(strings.NewReader(port))
	addr.port_op = C.PF_OP_NONE

	var tok rune
	curPort := 0
	for tok != scanner.EOF {
		tok = s.Scan()
		switch tok {
		case -3:
			if curPort >= 2 {
				return fmt.Errorf("Unexpected 3rd number in port range: %s",
					s.TokenText())
			}
			val, err := strconv.ParseUint(s.TokenText(), 10, 16)
			if err != nil {
				return fmt.Errorf("Number not allowed in port range: %s",
					s.TokenText())
			}
			if val < 0 {
				return fmt.Errorf("Port number can't be negative: %d", val)
			}

			addr.port[curPort] = C.u_int16_t(C.htons(C.uint16_t(val)))
			curPort++

			// if it is the first number and after there is nothing, set none
		case ':':
			addr.port_op = C.PF_OP_RRG
		case '!':
			if curPort != 0 {
				return fmt.Errorf("Unexpected number before '!'")
			}
			if s.Peek() == '=' {
				s.Next() // consume
				addr.port_op = C.PF_OP_NE
			} else {
				return fmt.Errorf("Expected '=' after '!'")
			}
		case '<':
			if s.Peek() == '>' {
				s.Next() // consume
				addr.port_op = C.PF_OP_XRG
			} else if s.Peek() == '=' {
				s.Next() // consume
				addr.port_op = C.PF_OP_LE
			} else if s.Peek() >= '0' && s.Peek() <= '9' { // int
				// next is port number continue
				addr.port_op = C.PF_OP_LT
			} else {
				return fmt.Errorf("Expected port number not '%c'", s.Peek())
			}
		case '>':
			if s.Peek() == '<' {
				s.Next() // consume
				addr.port_op = C.PF_OP_IRG
			} else if s.Peek() == '=' {
				s.Next() // consume
				addr.port_op = C.PF_OP_GE
			} else if s.Peek() >= '0' && s.Peek() <= '9' { // int
				// next is port number continue
				addr.port_op = C.PF_OP_GT
			} else {
				return fmt.Errorf("Expected port number not '%c'", s.Peek())
			}
		case -1:
			// if no operation was set
			if curPort == 1 && addr.port_op == C.PF_OP_NONE { // one port
				addr.port_op = C.PF_OP_EQ
			}
		default:
			return fmt.Errorf("Unexpected char '%c'", s.Peek())
		}
	}
	return nil
}