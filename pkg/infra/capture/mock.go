package capture

import "github.com/google/gopacket"

type Mock struct {
	Packets []gopacket.Packet
}

func (x *Mock) Read() chan gopacket.Packet {
	ch := make(chan gopacket.Packet)
	go func() {
		for i := range x.Packets {
			ch <- x.Packets[i]
		}
		close(ch)
	}()
	return ch
}
