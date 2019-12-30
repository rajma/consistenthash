package consistenthash

import (
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"
)

//ConsistentHash ...
type ConsistentHash struct {
	circle     map[*big.Int]string
	hasher     Hasher
	replicas   int32
	nodes      []string
	sortedKeys []*big.Int
	mux        sync.Mutex
}

//NewMD5Hash ...
func NewMD5Hash(replicas int32, nodes []string) *ConsistentHash {
	c := make(map[*big.Int]string)
	ch := ConsistentHash{
		circle:   c,
		hasher:   &MD5Hasher{},
		replicas: InitializeMin(replicas, 100),
		nodes:    nodes,
	}
	initializeCircle(ch.circle, ch.nodes, ch.replicas)
	ch.sortedKeys = make([]*big.Int, 0, len(ch.circle))
	for k := range ch.circle {
		ch.sortedKeys = append(ch.sortedKeys, k)
	}
	//Sort the keys in Ascending Order
	sort.Slice(ch.sortedKeys,
		func(i, j int) bool {
			return ch.sortedKeys[i].Cmp(ch.sortedKeys[j]) == -1
		})
	return &ch
}

//AddNode ...
func (ch *ConsistentHash) AddNode(node string) {
	if len(node) != 0 {
		//We need to do this with a mutex
		ch.mux.Lock()
		{
			addNode(node, ch.replicas, ch.circle)
			ch.sortedKeys = make([]*big.Int, 0, len(ch.circle))
			for k := range ch.circle {
				ch.sortedKeys = append(ch.sortedKeys, k)
			}
			//Sort the keys in Ascending Order
			sort.Slice(ch.sortedKeys, func(i, j int) bool {
				return ch.sortedKeys[i].Cmp(ch.sortedKeys[j]) == -1
			})
		}
		ch.mux.Unlock()
	}
}

//GetShard ...
func (ch *ConsistentHash) GetShard(key string) string {
	if len(key) != 0 && len(ch.circle) != 0 {
		keyHash := ch.hasher.GetHash(key)
		insertIdx := sort.Search(len(ch.sortedKeys), func(i int) bool {
			return ch.sortedKeys[i].Cmp(keyHash) > -1
		})
		if insertIdx == len(ch.sortedKeys) {
			insertIdx = 0
		}
		return ch.circle[ch.sortedKeys[insertIdx]]
	}
	return ""
}

func initializeCircle(circle map[*big.Int]string, nodes []string, replicas int32) {
	for _, val := range nodes {
		addNode(val, replicas, circle)
	}
}

func addNode(node string, replicas int32, circle map[*big.Int]string) error {
	var i int32
	hasher := &MD5Hasher{}
	for i = 0; i < replicas; i++ {
		shardPoint := fmt.Sprintf("%s%d", node, i)
		shardPointHash := hasher.GetHash(shardPoint)
		if _, exists := circle[shardPointHash]; !exists {
			circle[shardPointHash] = node
		} else {
			return &ApplicationError{
				When: time.Now(),
				What: fmt.Sprintf("Duplicate Key in circle:%s for %s", shardPointHash.String(), shardPoint),
			}
		}
	}
	return nil
}
