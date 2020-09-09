package market

type orderMinHeap []*Order
type orderMaxHeap []*Order

func (h orderMinHeap) Len() int { return len(h) }
func (h orderMaxHeap) Len() int { return len(h) }

func (h orderMinHeap) Less(i, j int) bool {
	return h[i].Price < h[j].Price
}
func (h orderMaxHeap) Less(i, j int) bool {
	return h[i].Price > h[j].Price
}

func (h orderMinHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }
func (h orderMaxHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *orderMinHeap) Push(x interface{}) {
	o := x.(*Order)
	*h = append(*h, o)
}
func (h *orderMaxHeap) Push(x interface{}) {
	o := x.(*Order)
	*h = append(*h, o)
}

func (h *orderMinHeap) Pop() interface{} {
	arr := *h
	n := len(arr)
	o := arr[n-1]
	*h = arr[0 : n-1]
	return o
}
func (h *orderMaxHeap) Pop() interface{} {
	arr := *h
	n := len(arr)
	o := arr[n-1]
	*h = arr[0 : n-1]
	return o
}

func (h orderMinHeap) Peek() *Order { return h[0] }
func (h orderMaxHeap) Peek() *Order { return h[0] }
