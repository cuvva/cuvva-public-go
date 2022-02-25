package ksuid

import (
	"bytes"
	crypto_rand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	math_rand "math/rand"
	"net"
	"os"
)

var random *math_rand.Rand

func init() {
	var b [8]byte

	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot seed random bytes")
	}

	random = math_rand.New(math_rand.NewSource(int64(binary.LittleEndian.Uint64(b[:]))))
}

type InstanceID struct {
	SchemeData byte
	BytesData  [8]byte
}

func (i InstanceID) Scheme() byte {
	return i.SchemeData
}

func (i InstanceID) Bytes() [8]byte {
	return i.BytesData
}

// NewHardwareID returns a HardwareID for the current node.
func NewHardwareID() (InstanceID, error) {
	hwAddr, err := getHardwareAddr()
	if err != nil {
		return InstanceID{}, err
	}

	var bd [8]byte
	copy(bd[:], hwAddr)
	binary.BigEndian.PutUint16(bd[6:], uint16(os.Getpid()))

	return InstanceID{
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

// NewDockerID returns a DockerID for the current Docker container.
func NewDockerID() (InstanceID, error) {
	cid, err := getDockerID()
	if err != nil {
		return InstanceID{}, err
	}

	var b [8]byte
	copy(b[:], cid)

	return InstanceID{
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

// NewRandomID returns a RandomID initialized by a PRNG.
func NewRandomID() InstanceID {
	tmp := make([]byte, 8)
	random.Read(tmp)

	var b [8]byte
	copy(b[:], tmp)

	return InstanceID{
		SchemeData: 'R',
		BytesData:  b,
	}
}
