xlru
==========

This provides the `lru` package which implements a fixed-size
thread safe LRU cache with expire feature. It is based on [golang-lru](https://github.com/hashicorp/golang-lru).

# example
- More examples can be found in xxx_test.go files.
## LRU without expiration
```go
// LRU without expiration
l, _ := New[int, int](128)
for i := 0; i < 256; i++ {
	l.Add(i, i)
}
if l.Len() != 128 {
	panic(fmt.Sprintf("bad len: %v", l.Len()))
}
```

## LRU with expiration
```go
// LRU with expiration
l2, _ := NewWithExpire[int, int](2, 2*time.Second)
l2.Add(1, 1)
if !l2.Contains(1) {
	panic("1 should be contained")
}
time.Sleep(2 * time.Second)
if l2.Contains(1) {
	panic("1 should not be contained")
}
```

## ARC without expiration
```go
l3, _ := NewARC[int, int](128)
for i := 0; i < 256; i++ {
	l3.Add(i, i)
}
if l3.Len() != 128 {
	panic(fmt.Sprintf("bad len: %v", l3.Len()))
}
```

## ARC with expiration
```go
// ARC with expiration
l4, _ := NewARCWithExpire[int, int](128, 30*time.Second)
for i := 0; i < 256; i++ {
	l4.Add(i, i)
}
if l4.Len() != 128 {
	panic(fmt.Sprintf("bad len: %v", l4.Len()))
}
```

## 2q without expiration
```go
// 2q without expiration
l5, _ := New2Q[int, int](128)
for i := 0; i < 256; i++ {
	l5.Add(i, i)
}
if l5.Len() != 128 {
	panic(fmt.Sprintf("bad len: %v", l5.Len()))
}
```

## 2q with expiration
```go
// 2q with expiration
l6, _ := New2QWithExpire[int, int](128, 30*time.Second)
for i := 0; i < 256; i++ {
	l6.Add(i, i)
}
if l6.Len() != 128 {
	panic(fmt.Sprintf("bad len: %v", l6.Len()))
}
```
