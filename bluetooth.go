package main

import (
	"context"
	"flag"
	"strings"
	"time"

	log "github.com/cihub/seelog"

	"github.com/currantlabs/ble"
	"github.com/currantlabs/ble/examples/lib/dev"
)

var (
	device = "default"
	addr   = ""
	sub    = 0
	sd     = 10 * time.Second
)

func startBluetooth(name string) (err error) {
	flag.Parse()

	d, err := dev.NewDevice(device)
	if err != nil {
		log.Errorf("can't new device : %s", err)
		return
	}
	ble.SetDefaultDevice(d)

	// Default to search device with name of Gopher (or specified by user).
	filter := func(a ble.Advertisement) bool {
		return strings.ToUpper(a.LocalName()) == strings.ToUpper(name)
	}

	// If addr is specified, search for addr instead.
	if len(addr) != 0 {
		filter = func(a ble.Advertisement) bool {
			return strings.ToUpper(a.Address().String()) == strings.ToUpper(addr)
		}
	}

	// Scan for specified durantion, or until interrupted by user.
	log.Debug("Scanning for %s...", sd)
	ctx := ble.WithSigHandler(context.WithTimeout(context.Background(), sd))
	cln, err := ble.Connect(ctx, filter)
	if err != nil {
		log.Errorf("can't connect : %s", err)
		return
	}

	// Make sure we had the chance to print out the message.
	done := make(chan struct{})
	// Normally, the connection is disconnected by us after our exploration.
	// However, it can be asynchronously disconnected by the remote peripheral.
	// So we wait(detect) the disconnection in the go routine.
	go func() {
		<-cln.Disconnected()
		log.Debugf("[ %s ] is disconnected ", cln.Address())
		close(done)
	}()

	log.Debugf("Discovering profile...")
	p, err := cln.DiscoverProfile(true)
	if err != nil {
		log.Errorf("can't discover profile: %s", err)
		return
	}

	// Start the exploration.
	explore(cln, p)

	// Disconnect the connection. (On OS X, this might take a while.)
	log.Debugf("Disconnecting [ %s ]... (this might take up to few seconds on OS X)", cln.Address())
	cln.CancelConnection()

	<-done
	return nil
}

func explore(cln ble.Client, p *ble.Profile) error {
	for _, s := range p.Services {
		log.Debugf("    Service: %s %s, Handle (0x%02X)", s.UUID, ble.Name(s.UUID), s.Handle)

		for _, c := range s.Characteristics {
			log.Debugf("      Characteristic: %s %s",
				c.UUID, ble.Name(c.UUID))
			if (c.Property & ble.CharRead) != 0 {
				b, err := cln.ReadCharacteristic(c)
				if err != nil {
					log.Debugf("Failed to read characteristic: %s", err)
					continue
				}
				log.Debugf("        Value         %x | %q", b, b)
			}

		}
	}
	return nil
}
