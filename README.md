# 🦜 Gocache

![MIT](https://img.shields.io/badge/license-MIT-blue.svg)

gocache is a distributed caching library, adapted from groupcache, intended as a replacement for memcached in many cases.

## Installation

Install and update this package with `go get -u github.com/sinomoe/gocache`.

## Usage

1. define getter function(used when cache missed) implement Getter interface

    ```golang
    getFromDB := gocache.GetterFunc(func(key string) ([]byte, error) {
        data, err :=  db.Get(key) // get data from database by the key
        if err == nil {
            return data, nil
        }
		return nil, fmt.Errorf("%s not exist", key)
	})
    ```

2. initialize a group named as "scores" with max 2048 cached items

    ```golang
    goc := gocache.NewGroup("scores", 2<<10, getFromDB)
    ```

3. create httpPool instance using local addr

    ```golang
    peers := gocache.NewHTTPPool("http://localhost:8001")
    ```

4. set peers on httpPool instance, there are three peers, every addr is meant to a gocache peer

    ```golang
    peers.SetPeers("http://localhost:8002", "http://localhost:8002", "http://localhost:8003")
    ```

5. register httpPool(peers) instance on group instance

    ```golang
    goc.RegisterPeers(peers)
    ```

## License

MIT © sino
