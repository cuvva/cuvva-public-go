package ksuid

import (
	"context"
	"sync"
	"time"

	"github.com/cuvva/cuvva-public-go/lib/servicecontext"
)

var exportedNode *Node

func init() {
	var iid InstanceID
	var err error

	iid, err = NewDockerID()
	if err != nil {
		iid, err = NewHardwareID()
		if err != nil {
			iid = NewRandomID()
		}
	}

	exportedNode = NewNode(Production, iid)
}

// Production is the internal name for production ksuid, but is omitted
// during marshaling.
const Production = "prod"

// Node contains metadata used for ksuid generation for a specific machine.
type Node struct {
	InstanceID InstanceID

	ts  uint64
	seq uint32
	mu  sync.Mutex
}

// NewNode returns a ID generator for the current machine.
func NewNode(environment string, instanceID InstanceID) *Node {
	return &Node{
		InstanceID: instanceID,
	}
}

// Generate returns a new ID for the machine and resource configured.
func (n *Node) Generate(ctx context.Context, resource string) (id ID) {
	if info := servicecontext.GetContext(ctx); info != nil {
		id.Environment = info.Environment
	} else {
		id.Environment = Production
	}

	id.Resource = resource
	id.InstanceID = n.InstanceID

	n.mu.Lock()

	ts := uint64(time.Now().UTC().Unix())
	if (ts - n.ts) >= 1 {
		n.ts = ts
		n.seq = 0
	} else {
		n.seq++
	}

	id.Timestamp = ts
	id.SequenceID = n.seq

	n.mu.Unlock()

	return
}

// SetInstanceID overrides the default instance id in the exported node.
// This will effect all invocations of the Generate function.
func SetInstanceID(instanceID InstanceID) {
	exportedNode.InstanceID = instanceID
}

// Generate returns a new ID for the current machine and resource configured.
func Generate(ctx context.Context, resource string) ID {
	return exportedNode.Generate(ctx, resource)
}
