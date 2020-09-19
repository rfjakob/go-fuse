package fs

// returns the new value
func (b *rawBridge) nlookupAdd(ino uint64, delta int) int {
	b.nlookupMu.Lock()
	defer b.nlookupMu.Unlock()
	n := b.nlookupMap[ino]
	n += delta
	b.nlookupMap[ino] = n
	if n < 0 || n > 10000 {
		panic(n)
	}
	return n
}
