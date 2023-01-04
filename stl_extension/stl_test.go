package stlextension_test

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"testing"

	stlextension "github.com/memory-overflow/go-common-library/stl_extension"
)

func TestOrderMap(t *testing.T) {
	intmp := stlextension.OrderedMap{}
	floatmp := stlextension.OrderedMap{}
	stringmp := stlextension.OrderedMap{}
	// test empty
	if intmp.Empty() != true {
		t.Fatalf("test empty true failed")
	}
	// test int key insert int map
	if err := intmp.Insert(stlextension.IntKey(5), "BBB"); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	if err := intmp.Insert(stlextension.IntKey(10), "aaaa"); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	t.Logf("int map: %s", intmp.String())
	// test float key insert int map
	if err := intmp.Insert(stlextension.FloatKey(10), "bbb"); err == nil {
		t.Fatalf("test insert float key to int map error.")
	}
	// test string key insert int map
	if err := intmp.Insert(stlextension.StringKey("10"), "bbb"); err == nil {
		t.Fatalf("test insert string key to int map error.")
	}

	// test float key insert float map
	if err := floatmp.Insert(stlextension.FloatKey(3.14), "bbb"); err != nil {
		t.Fatalf("test insert float key to float map error: %v", err)
	}
	if err := floatmp.Insert(stlextension.FloatKey(3.15), "aaa"); err != nil {
		t.Fatalf("test insert float key to float map error: %v", err)
	}
	t.Logf("float map: %s", floatmp.String())
	// test int key insert float map
	if err := floatmp.Insert(stlextension.IntKey(10), "bbb"); err == nil {
		t.Fatalf("test insert float key to float map error.")
	}
	// test string key insert float map
	if err := floatmp.Insert(stlextension.StringKey("10"), "bbb"); err == nil {
		t.Fatalf("test insert string key to float map error.")
	}

	// test string key insert string map
	if err := stringmp.Insert(stlextension.StringKey("china"), "beijing"); err != nil {
		t.Fatalf("test insert string key to string map error: %v", err)
	}
	if err := stringmp.Insert(stlextension.StringKey("U.S."), "Washington"); err != nil {
		t.Fatalf("test insert string key to string map error: %v", err)
	}
	t.Logf("string map: %s", stringmp.String())
	// test int key insert string map
	if err := stringmp.Insert(stlextension.IntKey(10), "bbb"); err == nil {
		t.Fatalf("test insert int key to string map error.")
	}
	// test float key insert string map
	if err := stringmp.Insert(stlextension.FloatKey(10.0), "bbb"); err == nil {
		t.Fatalf("test insert float key to string map error.")
	}

	if err := intmp.Insert(stlextension.IntKey(4), 5); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	if err := intmp.Insert(stlextension.IntKey(7), 355); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	if err := intmp.Insert(stlextension.IntKey(-34), 34); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	if err := intmp.Insert(stlextension.IntKey(0), 0); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	t.Logf("int map: %s", intmp.String())
	// test count true
	num, err := intmp.Count(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test count error: %v", err)
	}
	if num == 0 {
		t.Fatalf("test count true failed")
	}
	// test count false
	num, err = intmp.Count(stlextension.IntKey(10000))
	if err != nil {
		t.Fatalf("test count error: %v", err)
	}
	if num != 0 {
		t.Fatalf("test count false failed")
	}

	// test find true
	v, b, err := intmp.Find(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test find error: %v", err)
	}
	if !b || v.(int) != 5 {
		t.Fatalf("test find true failed")
	}
	// test count false
	v, b, err = intmp.Find(stlextension.IntKey(10000))
	if err != nil {
		t.Fatalf("test find error: %v", err)
	}
	if b {
		t.Fatalf("test find false failed")
	}

	// test insert same key
	if err := intmp.Insert(stlextension.IntKey(4), "bbbb"); err != nil {
		t.Fatalf("test insert int key to int map error: %v", err)
	}
	t.Logf("int map: %s", intmp.String())
	v, b, err = intmp.Find(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test find error: %v", err)
	}
	if !b || v.(string) != "bbbb" {
		t.Fatalf("test insert same key failed")
	}

	// test erase
	if err := intmp.Erase(stlextension.IntKey(4)); err != nil {
		t.Fatalf("test erase error: %v", err)
	}
	t.Logf("int map: %s", intmp.String())
	v, b, err = intmp.Find(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test find error: %v", err)
	}
	if b {
		t.Fatalf("test erase")
	}

	// test empty false
	if intmp.Empty() {
		t.Fatalf("test empty false failed")
	}

	// test begin and next
	for key, value := intmp.Begin(); key != nil; key, value, err = intmp.Next(key) {
		if err != nil {
			t.Fatalf("test next error: %v", err)
		}
		t.Logf("key: %v, value: %v", key, value)
	}

	// test rbegin and prev
	for key, value := intmp.RBegin(); key != nil; key, value, err = intmp.Prev(key) {
		if err != nil {
			t.Fatalf("test prev error: %v", err)
		}
		t.Logf("key: %v, value: %v", key, value)
	}

	// test next
	key, value, err := intmp.Next(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test next error: %v", err)
	}
	if key != stlextension.IntKey(5) && value.(string) != "BBB" {
		t.Fatalf("test next failed")
	}
	key, value, err = intmp.Next(stlextension.IntKey(5))
	if err != nil {
		t.Fatalf("test next error: %v", err)
	}
	if key != stlextension.IntKey(7) && value.(int) != 355 {
		t.Fatalf("test next failed")
	}

	// test prev
	key, value, err = intmp.Prev(stlextension.IntKey(4))
	if err != nil {
		t.Fatalf("test prev error: %v", err)
	}
	if key != stlextension.IntKey(0) && value.(int) != 0 {
		t.Fatalf("test prev failed")
	}
	key, value, err = intmp.Prev(stlextension.IntKey(5))
	if err != nil {
		t.Fatalf("test prev error: %v", err)
	}
	if key != stlextension.IntKey(0) && value.(int) != 0 {
		t.Fatalf("test prev failed")
	}

	intmp.Erase(stlextension.IntKey(0))
	// test size
	if intmp.Size() != 4 {
		t.Fatalf("test size failed")
	}
	// test size
	if floatmp.Size() != 2 {
		t.Fatalf("test size failed")
	}
}

func TestPoj2503(t *testing.T) {
	r, _ := os.Open("/data/workspace/go-library/in")
	reader := bufio.NewReader(r)
	dit := stlextension.OrderedMap{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		info := strings.Split(strings.Split(line, "\n")[0], " ")
		if info[0] == "" {
			continue
		}
		if len(info) == 2 {
			key, value := info[1], info[0]
			dit.Insert(stlextension.StringKey(key), value)
		}
		if len(info) == 1 {
			key := info[0]
			v, b, _ := dit.Find(stlextension.StringKey(key))
			if b {
				fmt.Println(v.(string))
			} else {
				fmt.Println("eh")
			}
		}
	}
}

func TestHackerRankcppmaps(t *testing.T) {
	r, _ := os.Open("/data/workspace/go-library/in")
	dit := stlextension.OrderedMap{}
	q := 0
	_, err := fmt.Fscanf(r, "%d", &q)
	if err != nil {
		fmt.Println(err)
	}
	for ; q > 0; q-- {
		k, key, val := 0, "", 0
		fmt.Fscanf(r, "%d %s", &k, &key)
		if k == 1 {
			fmt.Fscanf(r, "%d", &val)
			v, b, _ := dit.Find(stlextension.StringKey(key))
			if !b {
				dit.Insert(stlextension.StringKey(key), val)
			} else {
				dit.Insert(stlextension.StringKey(key), v.(int)+val)
			}
		} else if k == 2 {
			dit.Erase(stlextension.StringKey(key))
		} else if k == 3 {
			v, b, _ := dit.Find(stlextension.StringKey(key))
			if !b {
				fmt.Println(0)
			} else {
				fmt.Println(v.(int))
			}
		}
	}
}
