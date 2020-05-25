package local_node

import (
	"crawlab/model"
	"github.com/apex/log"
	"github.com/cenkalti/backoff/v4"
	"go.uber.org/atomic"
	"sync"
	"time"
)

var localNode *LocalNode
var locker atomic.Int32
var once sync.Once

type LocalNode struct {
	node *model.Node
	sync.RWMutex
}

func (n *LocalNode) load(retry bool) (err error) {
	n.Lock()
	defer n.Unlock()
	var node model.Node
	if retry {
		b := backoff.NewConstantBackOff(1 * time.Second)
		err = backoff.Retry(func() error {
			node, err = model.GetCurrentNode()
			if err != nil {
				log.WithError(err).Warnf("Get current node info from database failed.  Will after %f seconds, try again.", b.NextBackOff().Seconds())
			}
			return err
		}, b)
	} else {
		node, err = model.GetCurrentNode()
	}

	if err != nil {
		return
	}
	n.node = &node
	return nil
}
func (n *LocalNode) watch() {
	timer := time.NewTicker(time.Second * 5)
	for range timer.C {
		if locker.CAS(0, 1) {

			err := n.load(false)

			if err != nil {
				log.WithError(err).Errorf("load current node from database failed,")
			}
			locker.Store(0)
		}
		continue
	}
}

func (n *LocalNode) Current() *model.Node {
	n.RLock()
	defer n.RUnlock()
	return n.node
}
func CurrentNode() *model.Node {
	once.Do(func() {
		localNode = &LocalNode{}
		_ = localNode.load(true)
		go localNode.watch()
	})
	return localNode.Current()
}
