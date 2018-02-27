package econerra

type orderMinHeap []*MarketOrder
type orderMaxHeap []*MarketOrder

func (h *orderMinHeap) Len() int { return len(h) }
func (h *orderMaxHeap) Len() int { return len(h) }

func (h *orderMinHeap) Less(i, j int) bool {
	return h[i].Price < h[j].Price
}
func (h *orderMaxHeap) Less(i, j int) bool {
	return h[i].Price > h[j].Price
}

func (h *orderMinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h *orderMaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *orderMinHeap) Push(x interface{}) {
	o := x.(*MarketOrder)
	h = append(h, o)
}
func (h *orderMaxHeap) Push(x interface{}) {
	o := x.(*MarketOrder)
	h = append(h, o)
}

func (h *orderMinHeap) Pop() {
	n := len(h)
	o := h[n-1]
	*h = (*h)[0 : n-1]
	return o
}
func (h *orderMaxHeap) Pop() {
	n := len(h)
	o := h[n-1]
	*h = (*h)[0 : n-1]
	return o
}

func (h *orderMinHeap) Peek() *MarketOrder { return h[0] }
func (h *orderMaxHeap) Peek() *MarketOrder { return h[0] }
