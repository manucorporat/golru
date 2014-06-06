#GoLRU
## Memory optimized Least-Recently-Used cache

###Install it
```
$ go get github.com/manucorporat/golru
```


### Check out the API reference
[https://godoc.org/github.com/manucorporat/golru](https://godoc.org/github.com/manucorporat/golru)

### API example

```
    // Creates a new cache with capacity=100 and 5
    cache := golru.New(100, golru.DefaultLRUSamples)
    
    // Add a new entry
    cache.Set("foo", []byte("bar"))
    
    // Read "foo" --> it prints "bar"
    fmt.Println(string(cache.Get("foo")))
    
    // Check existence
    doc := cache.Get("something")
    if doc == nil {
        fmt.Println("something doesn't exist")
    } else {
        fmt.Println("something exist: "+string(doc))
    }
    
    // Delete an entry
    cache.Del("foo")
    
    // Delete all the entries (flush cache)
    cache.Flush()
    
    // Return current number of entries and capacity
    size := cache.Len()
    capacity := cache.Capacity()
    fmt.Printf("Size is %d and capacity is %s", size, capacity)
```
    
    
Benchmarks:

2,3 GHz Intel Core i7,
8 GB 1600 MHz DDR3

```
BenchmarkSetFullSize	 5000000	       734 ns/op
BenchmarkSetHalfSize	 2000000	       961 ns/op
BenchmarkSet10000	     2000000	       950 ns/op
BenchmarkGetFullSize	 5000000	       347 ns/op
```
- 2.900.000 read operations per second  
- 1.002.000 write operations per second