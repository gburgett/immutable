package queue

import "testing"

func BenchmarkPush_EmptyQueue(b *testing.B) {
	q := NewQueue()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Push(i)
	}
}

func BenchmarkPush_SingleItem(b *testing.B) {
	q := NewQueue().Push(17)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Push(i)
	}
}

func BenchmarkPush_BalancedQueue(b *testing.B) {
	q := NewQueue().Push(17).Push(18)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = q.Push(i)
	}
}

func BenchmarkPop_EmptyQueue(b *testing.B) {
	q := NewQueue()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Pop()
	}
}

func BenchmarkPop_SingleItem(b *testing.B) {
	q := NewQueue().Push(21)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Pop()
	}
}

func benchmarkPop_MustRebuild(b *testing.B, items int) {
	q := NewQueue().Push(21)

	for i := 0; i < items; i++ {
		q = q.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.Pop()
	}
}

func BenchmarkPop_MustRebuild10Items(b *testing.B) {
	benchmarkPop_MustRebuild(b, 10)
}

func BenchmarkPop_MustRebuild1kItems(b *testing.B) {
	benchmarkPop_MustRebuild(b, 1000)
}

func BenchmarkPop_MustRebuild100kItems(b *testing.B) {
	benchmarkPop_MustRebuild(b, 100*1000)
}

func BenchmarkPeekNext_BestCase(b *testing.B) {
	q := NewQueueFrom([]interface{}{17, 18, 19})

	it := q.Peek()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.PeekNext(it)
	}
}

func benchmarkPeekNext_MustRebuild(b *testing.B, items int) {
	q := NewQueue().Push(17)

	for i := 0; i < items; i++ {
		q = q.Push(i)
	}

	it := q.Peek()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = q.PeekNext(it)
	}
}

func BenchmarkPeekNext_MustRebuild10Items(b *testing.B) {
	benchmarkPeekNext_MustRebuild(b, 10)
}

func BenchmarkPeekNext_MustRebuild1kItems(b *testing.B) {
	benchmarkPeekNext_MustRebuild(b, 1000)
}

func BenchmarkPeekNext_MustRebuild100kItems(b *testing.B) {
	benchmarkPeekNext_MustRebuild(b, 100*1000)
}

func benchmarkIterateWholeQueue(b *testing.B, items int) {
	queue := NewQueue()

	for i := 0; i < items; i++ {
		queue = queue.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for i, q := queue.Peek(), queue; i.HasValue; i, q = q.PeekNext(i) {
			//do nothing - time iteration loop only
		}
	}
}

func BenchmarkIterateWholeQueue_10Items(b *testing.B) {
	benchmarkIterateWholeQueue(b, 10)
}

func BenchmarkIterateWholeQueue_1kItems(b *testing.B) {
	benchmarkIterateWholeQueue(b, 1000)
}

func BenchmarkIterateWholeQueue_100kItems(b *testing.B) {
	benchmarkIterateWholeQueue(b, 100*1000)
}
