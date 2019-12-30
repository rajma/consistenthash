package consistenthash

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetShard(t *testing.T) {
	shards := []string{"redis1", "redis2", "redis3", "redis4"}
	ch := NewMD5Hash(100, shards)
	users := []string{"-84942321036308", "-76029520310209", "-68343931116147", "-54921760962352"}
	expected := []string{"redis4", "redis4", "redis2", "redis3"}
	actual := []string{}
	for _, user := range users {
		shard := ch.GetShard(user)
		actual = append(actual, shard)
	}
	fmt.Println("Actual:", actual)
	assert.Equal(t, actual, expected, "The two slices should be the same.")

}
