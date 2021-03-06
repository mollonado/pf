# PF (Packet Filter)

[![GoDoc](https://godoc.org/github.com/go-freebsd/pf?status.svg)](https://godoc.org/github.com/go-freebsd/pf)
[![Coverage 84.7%](https://img.shields.io/badge/coverage-84.7%25-green.svg)]()
[![FreeBSD 10.3](https://img.shields.io/badge/freebsd-10.3-green.svg)](https://www.freebsd.org/releases/10.3R/announce.html)
[![FreeBSD 11](https://img.shields.io/badge/freebsd-11-green.svg)](https://www.freebsd.org/releases/11.0R/announce.html)
[![FreeBSD HEAD](https://img.shields.io/badge/freebsd-HEAD-green.svg)](https://svnweb.freebsd.org/base/head/)

The FreeBSD operating system has multiple packet filter build-in. One of
the packet filters was ported from OpenBSD and is called pf (packetfilter).

Packet filtering restricts the types of packets that pass through network
interfaces entering or leaving the host based on filter rules as
described in. The packet filter can also replace addresses
and ports of packets. Replacing source addresses and ports of outgoing
packets is called NAT (Network Address Translation) and is used to
connect an internal network (usually reserved address space) to an
external one (the Internet) by making all connections to external hosts
appear to come from the gateway. Replacing destination addresses and
ports of incoming packets is used to redirect connections to different
hosts and/or ports. A combination of both translations, bidirectional
NAT, is also supported.

This go module enables easy access to the packet filter inside the
kernel. The FreeBSD kernel module responsible for implementing pf is
called pf.ko.

Since the kernel interface is different between the operating
systems this version currently only works with FreeBSD.

The packet filter creates the pseudo-device node /dev/pf,
it allows userland processes to control the behavior of the packet filter
through an ioctl(2) interface. There are commands to enable and disable
the filter, load rulesets, add and remove individual rules or state table
entries, and retrieve statistics. The most commonly used functions are
covered by this library.

Manipulations like loading a ruleset that involve more than a single
ioctl(2) call require a so-called ticket, which prevents the occurrence
of multiple concurrent manipulations. Tickets are modeled as transaction
objects inside the library.

Working with pf directly on a remote connection can cause you to loose
the connection in case of a programming error. Make sure you have a
second way to access the system e.g. a serial console.

# Testing

You need to be root to execute the tests.

	make test
