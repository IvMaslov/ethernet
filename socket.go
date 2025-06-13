package ethernet

import (
	"fmt"
	"net"

	"github.com/IvMaslov/netutils"
	"github.com/IvMaslov/socket"
)

const defaultMTU = 1500

type EtherSocket struct {
	sock   *socket.Interface
	inInfo *netutils.InterfaceInfo
	mtu    int
}

// NewEtherSocket create ether socket over interface with mtu.
// if mtu == 0, default value is 1500
func NewEtherSocket(sock *socket.Interface, mtu int) (*EtherSocket, error) {
	es := &EtherSocket{
		sock: sock,
		mtu:  mtu,
	}

	if mtu == 0 {
		es.mtu = defaultMTU
	}

	// get info about defualt gateway of this interface
	info, err := netutils.GetDefaultGatewayInfo(sock.Name())
	if err != nil {
		return nil, err
	}

	es.inInfo = &info

	return es, nil
}

// Read returns payload of ethernet frame
func (es *EtherSocket) Read() ([]byte, error) {
	f, err := es.readFrame()
	if err != nil {
		return nil, err
	}

	return f.Payload, nil
}

// ReadFrame returns full thernet frame with payload
func (es *EtherSocket) ReadFrame() (*Frame, error) {
	return es.readFrame()
}

func (es *EtherSocket) readFrame() (*Frame, error) {
	buf := make([]byte, es.mtu)

	n, err := es.sock.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read from socket: %w", err)
	}

	f := &Frame{}

	err = f.Unmarshal(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal ethernet packet: %w", err)
	}

	return f, nil
}

func (es *EtherSocket) WriteTo(to net.HardwareAddr, data []byte) error {
	f := Frame{
		DestHarwAddr: to,
		SrcHarwAddr:  es.sock.GetHardwareAddr(),
		EtherType:    EtherTypeIPv4,
		Payload:      data,
	}

	return es.WriteFrame(&f)
}

func (es *EtherSocket) Write(data []byte) error {
	f := Frame{
		DestHarwAddr: es.inInfo.HardAddr,
		SrcHarwAddr:  es.sock.GetHardwareAddr(),
		EtherType:    EtherTypeIPv4,
		Payload:      data,
	}

	return es.WriteFrame(&f)
}

func (es *EtherSocket) WriteFrame(f *Frame) error {
	bts, err := f.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal ethernet frame: %w", err)
	}

	_, err = es.sock.Write(bts)

	return err
}
