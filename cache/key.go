package cache

import (
	"fmt"
	"strconv"
)

const (
	RankKey = "rank"
)

// ProductViewKey 将商品 id 映射成 redis 的相应键 key
func ProductViewKey(id uint) string {
	return fmt.Sprintf("view:product:%s", strconv.Itoa(int(id)))
}
