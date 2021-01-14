package dao

import (
	"bytes"
	"context"
	"encoding/gob"
	"flag"
	"testing"

	"github.com/go-kratos/kratos/pkg/cache/redis"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/log"
)

var d *dao
var ctx = context.Background()
var myredis *redis.Redis

// func TestMain(m *testing.M) {
// 	flag.Set("conf", "../../test")
// 	flag.Set("f", "../../test/docker-compose.yaml")
// 	flag.Parse()
// 	disableLich := os.Getenv("DISABLE_LICH") != ""
// 	if !disableLich {
// 		if err := lich.Setup(); err != nil {
// 			panic(err)
// 		}
// 	}
// 	var err error
// 	if err = paladin.Init(); err != nil {
// 		panic(err)
// 	}
// 	var cf func()
// 	if d, cf, err = newTestDao(); err != nil {
// 		panic(err)
// 	}
// 	ret := m.Run()
// 	cf()
// 	if !disableLich {
// 		_ = lich.Teardown()
// 	}
// 	os.Exit(ret)
// }

func TestRedis(t *testing.T) {
	flag.Set("conf", "../../configs")
	flag.Parse()
	paladin.Init()
	// _, closeFunc, err := di.InitApp()
	// if err != nil {
	// 	panic(err)
	// }
	r, cleanup, err := NewRedis()
	if err != nil {
		panic(err)
	}
	myredis = r
	defer cleanup()
	// 取集合
	// listGet(t)
	// 删集合
	// listRemove(t)
	// 删除指定元素
	listTargetRemove(t)
}

func TestGetRedis(t *testing.T) {
	flag.Set("conf", "../../configs")
	flag.Parse()
	paladin.Init()
	r, cleanup, err := NewRedis()
	if err != nil {
		panic(err)
	}
	myredis = r
	defer cleanup()
	// 存集合
	listGet(t)
}

func TestAddRedis(t *testing.T) {
	flag.Set("conf", "../../configs")
	flag.Parse()
	paladin.Init()
	r, cleanup, err := NewRedis()
	if err != nil {
		panic(err)
	}
	myredis = r
	defer cleanup()
	// 存集合
	listSet(t)
}

func listTargetRemove(t *testing.T) {
	key := "ItemList"
	conn := myredis.Conn(ctx)
	defer conn.Close()
	target := Item{
		Name: "marsonshine",
		Sex:  "man",
		Age:  28,
	}
	if d, err := myredis.Do(ctx, "LREM", key, 0, target); err != nil {
		t.Fatalf("获取列表失败 错误信息：%v", err)
	} else {
		log.Info("成功删除元素 %s 数量为 %v", key, d.(int64))
	}
}

func listRemove(t *testing.T) {
	key := "ItemList"
	conn := myredis.Conn(ctx)
	defer conn.Close()
	if d, err := myredis.Do(ctx, "DEL", key); err != nil {
		t.Fatalf("获取列表失败")
	} else {
		t.Logf("成功删除元素 %s 数量为 %v", key, d)
	}
}

func listGet(t *testing.T) {
	key := "ItemList"
	conn := myredis.Conn(ctx)
	defer conn.Close()
	// 取集合
	if buffer, err := redis.ByteSlices((myredis.Do(ctx, "LRANGE", key, -1, 0))); err != nil {
		t.Error("获取列表失败")
	} else {
		albums := new(struct {
			Name string
			Age  int
			Sex  string
		})
		buf := bytes.NewBuffer(buffer[0])
		dec := gob.NewDecoder(buf)
		dec.Decode(albums)

		// redis.ScanSlice(d, &albums)
		// (data, d, `redis:"Name"`, `redis:"Sex"`, `redis:"Age"`)
		t.Logf("成功取出集合 %v", d)
	}
}

func listSet(t *testing.T) {
	key := "ItemList"
	c := ctx
	conn := myredis.Conn(c)
	defer conn.Close()
	cacheItem := []interface{}{
		key,
		Item{
			Name: "marsonshine",
			Sex:  "man",
			Age:  28,
		},
		Item{
			Name: "marsonshine",
			Sex:  "man",
			Age:  28,
		},
	}
	if count, err := myredis.Do(c, "LPUSH", cacheItem...); err != nil {
		t.Errorf("conn.LPUSH(%s) error(%v)", key, err)
	} else {
		if count.(int64) > 0 {
			t.Log("插入list成功")
		} else {
			t.Errorf("插入失败，没有全部插入, 预期值 s=%d 实际值 s=%v", len(cacheItem), count)
		}
	}
}

type Item struct {
	Name string
	Sex  string
	Age  int
}
