package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/hanwen/go-fuse/v2/fs"
	"github.com/hanwen/go-fuse/v2/fuse"
)

type Object struct {
	fs.Inode
	json map[string]interface{}
}

func (o *Object) OnAdd(ctx context.Context) {
	for key, val := range o.json {
		handle(&o.Inode, ctx, key, val)
	}
}

type Array struct {
	fs.Inode
	json []interface{}
}

func (a *Array) OnAdd(ctx context.Context) {
	for i, val := range a.json {
		handle(&a.Inode, ctx, fmt.Sprintf("%v", i), val)
	}
}

func handle(i *fs.Inode, ctx context.Context, name string, val interface{}) {
	switch v := val.(type) {
	case nil:
		newFile(i, ctx, name, "")
	case string:
		newFile(i, ctx, name, v)
	case bool:
		newFile(i, ctx, name, fmt.Sprintf("%v", v))
	case float64:
		newFile(i, ctx, name, fmt.Sprintf("%v", v))
	case []interface{}:
		i.AddChild(name, i.NewPersistentInode(
			ctx,
			&Array{json: v},
			fs.StableAttr{Mode: fuse.S_IFDIR},
		), false)
	case map[string]interface{}:
		i.AddChild(name, i.NewPersistentInode(
			ctx,
			&Object{json: v},
			fs.StableAttr{Mode: fuse.S_IFDIR},
		), false)
	}
}

func newFile(i *fs.Inode, ctx context.Context, name, data string) {
	i.AddChild(name, i.NewPersistentInode(ctx, &fs.MemRegularFile{
		Data: []byte(data),
		Attr: fuse.Attr{
			Mode: 0444,
		},
	}, fs.StableAttr{}), false)
}

func main() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s <json> <mountpoint>\n", os.Args[0])
		os.Exit(1)
	}

	start := time.Now()

	contents, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var server *fuse.Server
	object := Object{}
	array := Array{}

	if err1 := json.Unmarshal(contents, &object.json); err1 != nil {
		if err2 := json.Unmarshal(contents, &array.json); err2 != nil {
			fmt.Fprintln(os.Stderr, "Failed to load JSON")
			fmt.Fprintln(os.Stderr, "Load as object: ", err1)
			fmt.Fprintln(os.Stderr, "Load as array:  ", err2)
			os.Exit(1)
		} else {
			server, err = fs.Mount(os.Args[2], &array, &fs.Options{})
		}
	} else {
		server, err = fs.Mount(os.Args[2], &object, &fs.Options{})
	}

	if err != nil {
		panic(err)
	}

	stop := time.Now()
	diff := stop.Sub(start)
	fmt.Println("Startup took", diff)

	server.Wait()
}
