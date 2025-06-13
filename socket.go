package ethernet

import (
	"fmt"
	"net"

	"github.com/IvMaslov/netutils"
	"github.com/IvMaslov/socket"
)

type EtherSocket struct {
	// raw socket
	sock *socket.Interface
	// information about default gateway of socket
	defaultInfo *netutils.InterfaceInfo
	mtu         int
}

// NewEtherSocket create ether socket over device
func NewEtherSocket(sock *socket.Interface) (*EtherSocket, error) {
	es := &EtherSocket{
		sock: sock,
	}

	// get info about defualt gateway of this interface
	info, err := netutils.GetDefaultGatewayInfo(sock.Name())
	if err != nil {
		return nil, err
	}

	es.defaultInfo = &info

	mtu, err := netutils.GetInterfaceMTU(sock.Name())
	if err != nil {
		return nil, err
	}

	es.mtu = mtu

	return es, nil
}

// Name returns device name
func (es *EtherSocket) Name() string {
	return es.sock.Name()
}

// MTU returns maximum transmission unit of device
func (es *EtherSocket) MTU() int {
	return es.mtu
}

// GetHWAddr returns mac address of device
func (es *EtherSocket) GetHWAddr() net.HardwareAddr {
	return es.sock.GetHardwareAddr()
}

// GetGatewayHWAddr returns mac address of device's default gateway
func (es *EtherSocket) GetGatewayHWAddr() net.HardwareAddr {
	return es.defaultInfo.HardAddr
}

// Read returns payload of ethernet frame
func (es *EtherSocket) Read() ([]byte, error) {
	f, err := es.ReadFrame()
	if err != nil {
		return nil, err
	}

	return f.Payload, nil
}

// ReadFrame returns full ethernet frame with payload
func (es *EtherSocket) ReadFrame() (*Frame, error) {
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

// WriteTo writes data to socket, with certain destination mac addr. Use Ether Type IPv4
func (es *EtherSocket) WriteTo(to net.HardwareAddr, data []byte) error {
	f := Frame{
		DestHarwAddr: to,
		SrcHarwAddr:  es.sock.GetHardwareAddr(),
		EtherType:    EtherTypeIPv4,
		Payload:      data,
	}

	return es.WriteFrame(&f)
}

// WriteTo writes data to socket, destination is default gateway. Use Ether Type IPv4
func (es *EtherSocket) Write(data []byte) error {
	f := Frame{
		DestHarwAddr: es.defaultInfo.HardAddr,
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
