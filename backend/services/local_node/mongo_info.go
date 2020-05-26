package local_node

import (
	"crawlab/model"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var locker atomic.Int32

type mongo struct {
	node *model.Node
	sync.RWMutex
}

func (n *mongo) load(retry bool) (err error) {
	n.Lock()
	defer n.Unlock()
	var node model.Node
	if retry {
		b := backoff.NewConstantBackOff(1 * time.Second)
		err = backoff.Retry(func() error {
			node, err = model.GetNodeByKey(GetLocalNode().Identify)
			if err != nil {
				log.WithError(err).Warnf("Get current node info from database failed.  Will after %f seconds, try again.", b.NextBackOff().Seconds())
			}
			return err
		}, b)
	} else {
		node, err = model.GetNodeByKey(localNode.Identify)
	}

	if err != nil {
		return
	}
	n.node = &node
	return nil
}
func (n *mongo) watch() {
	timer := time.NewTicker(time.Second * 5)
	for range timer.C {
		if locker.CAS(0, 1) {

			err := n.load(false)

			if err != nil {
				log.WithError(err).Errorf("load current node from database failed")
			}
			locker.Store(0)
		}
		continue
	}
}

func (n *mongo) Current() *model.Node {
	n.RLock()
	defer n.RUnlock()
	return n.node
}
