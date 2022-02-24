package ksuid

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// InstanceID is an interface implemented to identify a instance
// for a node in a unique manor.
type InstanceID interface {
	// Scheme returns the single byte used to identify the InstanceID.
	Scheme() byte

	// Bytes returns the serialized form of the InstanceID.
	Bytes() [8]byte
}

type IID struct {
	SchemeData byte
	BytesData  [8]byte
}

func (i IID) Scheme() byte {
	return i.SchemeData
}

func (i IID) Bytes() [8]byte {
	return i.BytesData
}

// ParseInstanceID unmarshals a prefixed node ID into its dedicated type.
func ParseInstanceID(b []byte) (InstanceID, error) {
	if len(b) != 9 {
		return nil, fmt.Errorf("expected 9 bytes, got %d", len(b))
	}

	switch b[0] {
	case 'H':
		return ParseHardwareID(b[1:])

	case 'D':
		return ParseDockerID(b[1:])

	case 'R':
		return ParseRandomID(b[1:])

	default:
		return nil, fmt.Errorf("unknown node id '%c'", b[0])
	}
}

// NewHardwareID returns a HardwareID for the current node.
func NewHardwareID() (IID, error) {
	hwAddr, err := getHardwareAddr()
	if err != nil {
		return IID{}, err
	}

	var bd [8]byte
	copy(bd[:], hwAddr)
	binary.BigEndian.PutUint16(bd[6:], uint16(os.Getpid()))

	return IID{
		SchemeData: 'H',
		BytesData:  bd,
	}, nil
}

func getHardwareAddr() (net.HardwareAddr, error) {
	addrs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		// only return physical interfaces (i.e. not loopback)
		if len(addr.HardwareAddr) >= 6 {
			return addr.HardwareAddr, nil
		}
	}

	return nil, fmt.Errorf("no hardware addr available")
}

// ParseHardwareID unmarshals a HardwareID from a sequence of bytes.
func ParseHardwareID(b []byte) (IID, error) {
	if len(b) != 8 {
		return IID{}, fmt.Errorf("expected 8 bytes, got %d", len(b))
	}

	machineID := net.HardwareAddr(b[:6])
	processID := binary.BigEndian.Uint16(b[6:])

	var bd [8]byte
	copy(bd[:], machineID)
	binary.BigEndian.PutUint16(bd[6:], processID)

	return IID{
		SchemeData: 'H',
		BytesData:  bd,
	}, nil
}

// NewDockerID returns a DockerID for the current Docker container.
func NewDockerID() (IID, error) {
	cid, err := getDockerID()
	if err != nil {
		return IID{}, err
	}

	var b [8]byte
	copy(b[:], cid)

	return IID{
		SchemeData: 'D',
		BytesData:  b,
	}, nil
}

func getDockerID() ([]byte, error) {
	src, err := ioutil.ReadFile("/proc/1/cpuset")
	src = bytes.TrimSpace(src)
	if os.IsNotExist(err) || len(src) < 64 || !bytes.HasPrefix(src, []byte("/docker")) {
		return nil, fmt.Errorf("not a docker container")
	} else if err != nil {
		return nil, err
	}

	dst := make([]byte, 32)
	_, err = hex.Decode(dst, src[len(src)-64:])
	if err != nil {
		return nil, err
	}

	return dst, nil
}

// ParseDockerID unmarshals a DockerID from a sequence of bytes.
func ParseDockerID(b []byte) (IID, error) {
	if len(b) != 8 {
		return IID{}, fmt.Errorf("expected 8 bytes, got %d", len(b))
	}

	var bd [8]byte
	copy(bd[:], b)

	return IID{
		SchemeData: 'D',
		BytesData:  bd,
	}, nil
}

// NewRandomID returns a RandomID initialized by a PRNG.
func NewRandomID() (IID, error) {
	tmp := make([]byte, 8)
	rand.Read(tmp)

	var b [8]byte
	copy(b[:], tmp)

	return IID{
		SchemeData: 'R',
		BytesData:  b,
	}, nil
}

// ParseRandomID unmarshals a RandomID from a sequence of bytes.
func ParseRandomID(b []byte) (IID, error) {
	if len(b) != 8 {
		return IID{}, fmt.Errorf("expected 8 bytes, got %d", len(b))
	}

	var x [8]byte
	copy(x[:], b)

	return IID{
		SchemeData: 'R',
		BytesData:  x,
	}, nil
}
