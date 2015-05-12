package queue

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestPush(t *testing.T) {
	Convey("Given an empty queue", t, func() {
		queue := NewQueue()

		Convey("When a single item is pushed", func() {
			q := queue.Push(7)

			Convey("Then the count of the queue should be 1", func() {
				So(q.Count(), ShouldEqual, 1)
			})

			Convey("And the item should be at the head of the queue", func() {
				p := q.Peek()
				So(p.HasValue, ShouldBeTrue)
				So(p.Value, ShouldEqual, 7)
			})

			Convey("And the original queue should remain unchanged", func() {
				So(queue.Count(), ShouldEqual, 0)
				So(q, ShouldNotEqual, queue)
			})
		})

		Convey("When multiple items are pushed", func() {
			q := queue.Push(7).Push(8).Push(9)

			Convey("Then the queues count should equal the number of items added", func() {
				So(q.Count(), ShouldEqual, 3)
			})

			Convey("And each item should pop in order", func() {
				v, q := q.Pop()
				So(v, ShouldEqual, 7)
				v, q = q.Pop()
				So(v, ShouldEqual, 8)
				v, q = q.Pop()
				So(v, ShouldEqual, 9)
			})

			Convey("And the original queue should remain unchanged", func() {
				So(queue.Count(), ShouldEqual, 0)
				So(q, ShouldNotEqual, queue)
			})
		})
	})
}

func TestPop(t *testing.T) {
	Convey("Given an empty queue", t, func() {
		queue := NewQueue()

		Convey("When an item is popped", func() {
			v, q := queue.Pop()

			Convey("Then an empty item should be returned", func() {
				So(v, ShouldBeNil)
			})

			Convey("And the queue should remain empty", func() {
				So(q.Count(), ShouldEqual, 0)
				So(q, ShouldEqual, queue)
			})
		})
	})

	Convey("Given a queue with one value", t, func() {
		queue := NewQueue().Push(10)

		Convey("When an item is popped", func() {
			v, q := queue.Pop()

			Convey("Then the single item should be returned", func() {
				So(v, ShouldEqual, 10)
			})

			Convey("And the new queue should be empty", func() {
				So(q.Count(), ShouldEqual, 0)
			})

			Convey("And the original queue should remain unchanged", func() {
				So(queue.Count(), ShouldEqual, 1)
				So(q, ShouldNotEqual, queue)
			})
		})
	})

	Convey("Given a queue with multiple values", t, func() {
		queue := NewQueue().Push(10).Push(11).Push(12)

		Convey("When an item is popped", func() {
			v, q := queue.Pop()

			Convey("Then the single item should be returned", func() {
				So(v, ShouldEqual, 10)
			})

			Convey("And the remaining items should pop in order", func() {
				v, q = q.Pop()
				So(v, ShouldEqual, 11)
				v, q = q.Pop()
				So(v, ShouldEqual, 12)
			})
		})
	})
}

func TestPeek(t *testing.T) {
	Convey("Given an empty queue", t, func() {
		queue := NewQueue()

		Convey("When the first item is peeked", func() {
			i := queue.Peek()

			Convey("Then the peek iterator should be empty", func() {
				So(i.HasValue, ShouldBeFalse)
				So(i.Value, ShouldBeNil)
			})
		})
	})

	Convey("Given a queue with one item", t, func() {
		queue := NewQueue().Push(25)

		Convey("When the first item is peeked", func() {
			i := queue.Peek()

			Convey("Then the peek iterator should contain the item", func() {
				So(i.HasValue, ShouldBeTrue)
				So(i.Value, ShouldEqual, 25)
			})
		})
	})
}

func TestPeekNext(t *testing.T) {
	Convey("Given an empty queue", t, func() {
		queue := NewQueue()

		Convey("And an empty peek iterator", func() {
			current := PeekIterator{}

			Convey("When the next item is peeked", func() {
				i, q := queue.PeekNext(current)

				Convey("Then an empty peek iterator is returned", func() {
					So(i.HasValue, ShouldBeFalse)
					So(i.Value, ShouldBeNil)
				})

				Convey("And the queue should not be modified", func() {
					So(q, ShouldEqual, queue)
				})
			})
		})
	})

	Convey("Given a queue with one item on it", t, func() {
		queue := NewQueue().Push(12)

		Convey("And a peek iterator for the first item", func() {
			current := queue.Peek()

			Convey("When the next item is peeked", func() {
				i, q := queue.PeekNext(current)

				Convey("Then an empty peek iterator is returned", func() {
					So(i.HasValue, ShouldBeFalse)
					So(i.Value, ShouldBeNil)
				})

				Convey("And the queue should not be modified", func() {
					So(q, ShouldEqual, queue)
				})

				Convey("And another peek returns an empty iterator", func() {
					i, q = q.PeekNext(i)
					So(i.HasValue, ShouldBeFalse)
				})
			})
		})
	})

	Convey("Given a queue with multiple items on it", t, func() {
		queue := NewQueue().Push(17).Push(18).Push(19)

		Convey("And a peek iterator for the first item", func() {
			current := queue.Peek()

			Convey("When the next item is peeked", func() {
				i, q := queue.PeekNext(current)

				Convey("Then the next item should be returned", func() {
					So(i.HasValue, ShouldBeTrue)
					So(i.Value, ShouldEqual, 18)
				})

				Convey("And the queue should have been rebuilt", func() {
					So(q, ShouldNotEqual, queue)
				})

				Convey("And peeking again with the original iterator should return the same result", func() {
					i, q2 := q.PeekNext(current)
					So(i.Value, ShouldEqual, 18)
					So(q2, ShouldEqual, q)
				})
			})

			Convey("When the first two items are popped", func() {
				_, q := queue.Pop()
				_, q = q.Pop()

				Convey("Then the next peeked item should be the head of the queue", func() {
					i, q2 := q.PeekNext(current)
					So(i.Value, ShouldEqual, 19)
					So(q2, ShouldEqual, q)
				})
			})
		})

		Convey("When a peek iterator crosses to the second generation of the queue", func() {
			current := queue.Peek()
			current, q2 := queue.PeekNext(current)
			So(q2, ShouldNotEqual, queue)

			Convey("Then peeking against the previous generation should panic", func() {
				So(func() { queue.PeekNext(current) }, ShouldPanicWith, "Unexpected queue generation - encountered 0 expected 1")
			})
		})
	})
}

func TestNewQueueFrom(t *testing.T) {
	Convey("Given a new queue from a slice", t, func() {
		items := []interface{}{1, 5, 2, 7, 4, 6}
		queue := NewQueueFrom(items...)

		Convey("When the items are popped", func() {
			popped := make([]int, 0, queue.Count())

			for v, q := queue.Pop(); v != nil; v, q = q.Pop() {
				popped = append(popped, v.(int))
			}

			Convey("Then the items should be in order", func() {
				So(popped, ShouldResemble, []int{1, 5, 2, 7, 4, 6})
			})
		})

		Convey("When the items are iterated", func() {
			iterated := make([]int, 0, queue.Count())

			for i, q := queue.Peek(), queue; i.HasValue; i, q = q.PeekNext(i) {
				iterated = append(iterated, i.Value.(int))
			}

			Convey("Then the items should be in order", func() {
				So(iterated, ShouldResemble, []int{1, 5, 2, 7, 4, 6})
			})

		})

	})
}
