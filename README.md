<div align="center">
<img alt="gophercises_jumping" src="https://user-images.githubusercontent.com/89729679/163836571-5b2fc3a6-2208-43e2-9478-56986ec73ab4.gif" />
  
# gobounce
Debounce &amp; Throttle inspired by Lodash's debounce and throttle
</div>

## Usage

### Debounce
```go
foo := func() {
  // do somthing anoying
}

var options *DebounceOptions // if nil, will have default options

debounced, cancel := gobounce.Debounce(foo, 2*time.Second, options)
defer cancel() // should be called to dispose timers and flush

debounced()
debounced()
debounced()
debounced()
debounced()
...
// after 2 seconds...
...
// foo called only once!
```

## Throttle

```go

foo := func() {
  // do somthing exausting
}

var options *ThrottleOptions // if nil, will have default options

throttled, cancel := gobounce.Throttle(foo, 2*time.Second, options)
defer cancel() // should be called to dispose timers and flush

// ------------- 00:00:00
throttled() // - 00:01:00
throttled() // - 00:02:00 - called!
throttled() // - 00:03:00
throttled() // - 00:04:00 - called!
throttled() // - 00:05:00
throttled() // - 00:06:00 - called!
// ------------- 00:07:00
```

## Author
Idan Izicovich
