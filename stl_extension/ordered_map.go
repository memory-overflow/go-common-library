package stlextension

import (
	"fmt"
	"math"
	"reflect"
	"sync"

	jsoniter "github.com/json-iterator/go"
)

func max(a, b int) (m int) {
	if a > b {
		return a
	} else {
		return b
	}
}

// Key key 必须要实现 Less 比较方法
type Key interface {
	// Less 是否小于 v
	Less(v Key) (bool, error)
}

// IntKey ...
type IntKey int64

// Less ...
func (i IntKey) Less(v Key) (bool, error) {
	x, ok := v.(IntKey)
	if ok {
		return i < x, nil
	}
	y, ok := v.(*IntKey)
	if ok {
		return i < *y, nil
	}
	return false, fmt.Errorf("can not compare between IntKey and %s", reflect.ValueOf(v).Type().Name())
}

// IntKey ...
type FloatKey float64

func sgn(x float64) int {
	if math.Abs(x) <= 1e-15 {
		return 0
	}
	if x < 0 {
		return -1
	}
	return 1
}

// Less ...
func (f FloatKey) Less(v Key) (bool, error) {

	x, ok := v.(FloatKey)
	if ok {
		return sgn(float64(f)-float64(x)) < 0, nil
	}
	y, ok := v.(*FloatKey)
	if ok {
		return sgn(float64(f)-float64(*y)) < 0, nil
	}
	return false, fmt.Errorf("can not compare between FloatKey and %s", reflect.ValueOf(v).Type().Name())
}

// StringKey ...
type StringKey string

// Less ...
func (s StringKey) Less(v Key) (bool, error) {
	x, ok := v.(StringKey)
	if ok {
		return s < x, nil
	}
	y, ok := v.(*StringKey)
	if ok {
		return s < *y, nil
	}
	return false, fmt.Errorf("can not compare between StringKey and %s", reflect.ValueOf(v).Type().Name())

}

// OrderedMap 按照 key 排序的 map，可以实现顺序遍历。
// OrderedMap 使用 avl 树实现，是线程安全的。
// 和 c++ map 对比测试，1000 万的随机的增删查操作：
// OrderedMap: 21806 ms, c++ map: 11592ms，效率比 c++ map 慢两倍。
// next 遍历整个 map, OrderedMap: 676ms, c++ map: 171ms；
// prev 遍历整个 map, OrderedMap: 663ms, c++ map: 198ms；
// 遍历的效率大概比 c++ map 的慢三倍。
type OrderedMap struct {
	root  *node      // 根节点指针
	size  int        // 存储元素数量
	mutex sync.Mutex // 并发控制锁
}

type node struct {
	key   Key         // 节点中存储的元素 key
	value interface{} // 节点中存储的元素 value
	num   int         // 该元素数量
	depth int         // 该节点的深度
	left  *node       // 左节点指针
	right *node       // 右节点指针
}

func (n *node) getDepth() (depth int) {
	if n == nil {
		return 0
	}
	return n.depth
}

func newNode(key Key, v interface{}) (n *node) {
	return &node{
		key:   key,
		value: v,
		num:   1,
		depth: 1,
		left:  nil,
		right: nil,
	}
}

// leftRotate 左旋转
func (n *node) leftRotate() (m *node) {

	headNode := n.right
	n.right = headNode.left
	headNode.left = n
	// 更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

// rightRotate 右旋转
func (n *node) rightRotate() (m *node) {

	headNode := n.left
	n.left = headNode.right
	headNode.right = n
	// 更新结点高度
	n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
	headNode.depth = max(headNode.left.getDepth(), headNode.right.getDepth()) + 1
	return headNode
}

// rightLeftRotate 右旋转, 之后左旋转
func (n *node) rightLeftRotate() (m *node) {
	sonHeadNode := n.right.rightRotate()
	n.right = sonHeadNode
	return n.leftRotate()
}

// leftRightRotate 左旋转, 之后右旋转
func (n *node) leftRightRotate() (m *node) {
	sonHeadNode := n.left.leftRotate()
	n.left = sonHeadNode
	return n.rightRotate()
}

// getNextAndSwap 交换后继结点
func (n *node) getNextAndSwap(swap *node) {
	if n == nil {
		return
	}
	if n.left == nil {
		swap.key, swap.value, swap.num = n.key, n.value, n.num
		n.num = 1
	} else {
		n.left.getNextAndSwap(swap)
	}
}

// adjust 以 node 平衡二叉树节点做接收者，对 n 节点进行旋转以保持节点左右子树平衡。
func (n *node) adjust() (m *node) {
	if n.right.getDepth()-n.left.getDepth() >= 2 {
		// 右子树高于左子树且高度差超过 2,此时应当对 n 进行左旋
		if n.right.right.getDepth() > n.right.left.getDepth() {
			// 由于右右子树高度超过右左子树,故可以直接左旋
			n = n.leftRotate()
		} else {
			// 由于右右子树不高度超过右左子树
			// 所以应该先右旋右子树使得右子树高度不超过左子树
			// 随后 n 节点左旋
			n = n.rightLeftRotate()
		}
	} else if n.left.getDepth()-n.right.getDepth() >= 2 {
		// 左子树高于右子树且高度差超过 2,此时应当对 n 进行右旋
		if n.left.left.getDepth() > n.left.right.getDepth() {
			// 由于左左子树高度超过左右子树,故可以直接右旋
			n = n.rightRotate()
		} else {
			// 由于左左子树高度不超过左右子树
			// 所以应该先左旋左子树使得左子树高度不超过右子树
			// 随后n节点右旋
			n = n.leftRightRotate()
		}
	}
	return n
}

func (n *node) insert(key Key, value interface{}) (m *node, isSame bool, err error) {
	// 节点不存在,应该创建并插入二叉树中
	if n == nil {
		return newNode(key, value), false, nil
	}
	ok, err := key.Less(n.key)
	if err != nil {
		return n, false, err
	}
	if ok {
		// 从左子树继续插入
		n.left, isSame, err = n.left.insert(key, value)
		if err != nil {
			return n, false, err
		}
		if !isSame {
			// 插入成功,对该节点进行平衡
			n = n.adjust()
			n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		}
		return n, isSame, nil
	}
	ok, err = n.key.Less(key)
	if err != nil {
		return n, false, err
	}
	if ok {
		// 从右子树继续插入
		n.right, isSame, err = n.right.insert(key, value)
		if err != nil {
			return n, isSame, err
		}
		if !isSame {
			// 插入成功, 对该节点进行平衡
			n = n.adjust()
			n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		}
		return n, isSame, nil
	}
	// 不允许重复,对值进行覆盖
	n.value = value
	return n, true, nil
}

func (n *node) erase(key Key) (m *node, b bool, err error) {
	if n == nil {
		// 待删除值不存在,删除失败
		return n, false, nil
	}
	less, err := key.Less(n.key)
	if err != nil {
		return n, false, err
	}
	big, err := n.key.Less(key)
	if err != nil {
		return n, false, err
	}
	if less {
		// 从左子树继续删除
		n.left, b, err = n.left.erase(key)
	} else if big {
		//从右子树继续删除
		n.right, b, err = n.right.erase(key)
	} else {
		// 存在相同值,从该节点删除
		b = true
		if n.num > 1 {
			// 有重复值,节点无需删除,直接 -1 即可
			n.num--
		} else {
			// 该节点需要被删除
			if n.left != nil && n.right != nil {
				// 找到该节点后继节点进行交换
				n.right.getNextAndSwap(n)
				//从右节点继续删除,同时可以保证删除的节点必然无左节点
				n.right, b, err = n.right.erase(n.key)
			} else if n.left != nil {
				n = n.left
			} else {
				n = n.right
			}
		}
	}
	// 当 n 节点仍然存在时,对其进行调整
	if n != nil && err == nil {
		n.depth = max(n.left.getDepth(), n.right.getDepth()) + 1
		n = n.adjust()
	}
	return n, b, err
}

func (n *node) count(key Key) (num int, err error) {
	if n == nil {
		return 0, nil
	}
	less, err := key.Less(n.key)
	if err != nil {
		return 0, err
	}
	big, err := n.key.Less(key)
	if err != nil {
		return 0, err
	}
	if less {
		return n.left.count(key)
	}
	if big {
		return n.right.count(key)
	}
	return n.num, nil
}

func (n *node) find(key Key) (value interface{}, b bool, err error) {
	if n == nil {
		return nil, false, nil
	}
	less, err := key.Less(n.key)
	if err != nil {
		return nil, false, err
	}
	big, err := n.key.Less(key)
	if err != nil {
		return nil, false, err
	}
	if less {
		return n.left.find(key)
	}
	if big {
		return n.right.find(key)
	}
	return n.value, true, nil
}

func (n *node) prev(key Key) (preNode *node, err error) {
	if n == nil || key == nil {
		return nil, nil
	}
	big, err := n.key.Less(key)
	if err != nil {
		return nil, err
	}
	if !big {
		return n.left.prev(key)
	}
	preNode = n
	if n.right != nil {
		tmpNode, err := n.right.prev(key)
		if err != nil {
			return nil, err
		}
		if tmpNode != nil {
			preNode = tmpNode
		}
	}
	return preNode, err
}

func (n *node) next(key Key) (nextNode *node, err error) {
	if n == nil || key == nil {
		return nil, nil
	}
	less, err := key.Less(n.key)
	if err != nil {
		return nil, err
	}
	if !less {
		return n.right.next(key)
	}
	nextNode = n
	if n.left != nil {
		tmpNode, err := n.left.next(key)
		if err != nil {
			return nil, err
		}
		if tmpNode != nil {
			nextNode = tmpNode
		}
	}
	return nextNode, err
}

func (n *node) findMinNode() *node {
	if n == nil {
		return nil
	}
	if n.left != nil {
		return n.left.findMinNode()
	}
	return n
}

func (n *node) findMaxNode() *node {
	if n == nil {
		return nil
	}
	if n.right != nil {
		return n.right.findMaxNode()
	}
	return n
}

// Insert 插入元素
func (om *OrderedMap) Insert(key Key, value interface{}) (err error) {
	if om == nil {
		return fmt.Errorf("OrderedMap is nil")
	}
	if key == nil {
		return fmt.Errorf("key is nil")
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	if om.root == nil {
		// 二叉树为空, 用根节点承载元素 key
		om.root = newNode(key, value)
		om.size = 1
		return
	}
	// 从根节点进行插入, 并返回节点, 同时返回是否插入成功
	var isSame bool
	om.root, isSame, err = om.root.insert(key, value)
	if err != nil {
		// 失败
		return err
	}
	if !isSame {
		om.size++
	}
	return nil
}

// Erase 删除元素
func (om *OrderedMap) Erase(key Key) (err error) {
	if om.Empty() {
		return
	}
	less, err := key.Less(om.root.key)
	if err != nil {
		return err
	}
	big, err := om.root.key.Less(key)
	if err != nil {
		return err
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	if om.size == 1 && !less && !big {
		// 二叉树仅持有一个元素且根节点等价于待删除元素, 将二叉树根节点置为 nil
		om.root = nil
		om.size = 0
		return
	}
	var b bool
	om.root, b, err = om.root.erase(key)
	if err != nil {
		return err
	}
	if b {
		om.size--
	}
	return nil
}

// Count 查找 key 的个数
func (om *OrderedMap) Count(key Key) (num int, err error) {
	if om.Empty() {
		return 0, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	num, err = om.root.count(key)
	return num, nil
}

// Find 查询 key 的 value
func (om *OrderedMap) Find(key Key) (value interface{}, found bool, err error) {
	if om.Empty() {
		return nil, false, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	return om.root.find(key)
}

// Begin 获取最小的元素的 key 和 value，配合 Next 进行迭代
func (om *OrderedMap) Begin() (key Key, value interface{}) {
	if om.Empty() {
		return nil, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	n := om.root.findMinNode()
	if n == nil {
		return nil, nil
	}
	return n.key, n.value
}

// Next 寻找比 key 大的下一个节点
func (om *OrderedMap) Next(key Key) (nextkey Key, nextvalue interface{}, err error) {
	if om.Empty() {
		return nil, nil, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	n, err := om.root.next(key)
	if err != nil {
		return nil, nil, err
	}
	if n == nil {
		return nil, nil, nil
	}
	return n.key, n.value, nil
}

// RBegin 获取最大的元素的 key 和 value，配合 Prev 进行迭代
func (om *OrderedMap) RBegin() (key Key, value interface{}) {
	if om == nil {
		return nil, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	n := om.root.findMaxNode()
	if n == nil {
		return nil, nil
	}
	return n.key, n.value
}

// Prev 寻找比 key 小的下一个节点
func (om *OrderedMap) Prev(key Key) (prekey Key, prevalue interface{}, err error) {
	if om.Empty() {
		return nil, nil, nil
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	n, err := om.root.prev(key)
	if err != nil {
		return nil, nil, err
	}
	if n == nil {
		return nil, nil, nil
	}
	return n.key, n.value, nil
}

// Clear clear order map
func (om *OrderedMap) Clear() {
	if om == nil {
		return
	}
	om.mutex.Lock()
	defer om.mutex.Unlock()
	om.root = nil
	om.size = 0
}

// Empty 是否为空
func (om *OrderedMap) Empty() bool {
	if om == nil {
		return true
	}
	return om.size == 0 || om.root == nil
}

// Size 返回 map 大小
func (om *OrderedMap) Size() int {
	if om == nil {
		return 0
	}
	return om.size
}

// String 转换成字符串
func (om *OrderedMap) String() string {
	if om.Empty() {
		return "{}"
	}
	m := map[interface{}]interface{}{}
	for key, value := om.Begin(); key != nil; key, value, _ = om.Next(key) {
		m[key] = value
	}
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	bs, _ := json.MarshalToString(m)
	return bs
}
