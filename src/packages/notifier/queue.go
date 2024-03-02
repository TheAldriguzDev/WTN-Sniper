package notifier

import (
	"sync"
	"time"
)

/*
Logic FIFO (First in first out)
*/

type Queue struct {
	Webhooks   []Discord_Params
	MaxRetries int
	mu         sync.Mutex
}

func CreateQueue(max_retries int) *Queue {
	return &Queue{
		Webhooks:   make([]Discord_Params, 0),
		MaxRetries: max_retries,
	}
}

func (q *Queue) Enqueue(discord_data Discord_Params) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.Webhooks = append(q.Webhooks, discord_data)
}

/*
Queue is active only if the element in Queue (every session is indipendent) is more than 20 webhooks
*/
func (q *Queue) StartQueue(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		if len(q.Webhooks) < 20 {
			q.processQueue()
		} else {
			select {
			case <-ticker.C:
				q.processQueue()
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

// Process the Queue by analizying chunks of 5 elements
func (q *Queue) processQueue() {
	if len(q.Webhooks) != 0 {
		var slicer int
		if len(q.Webhooks) > 5 {
			slicer = 5
		} else {
			slicer = len(q.Webhooks)
		}
		chunk := q.Webhooks[0:slicer]

		for i := len(chunk) - 1; i >= 0; i-- {
			webhook := chunk[i]
			err := SendWebhook(webhook)

			if err != nil {
				l.Error(err.Error())
				q.Webhooks[i].Retries++

				if q.Webhooks[i].Retries > q.MaxRetries {
					q.Webhooks = append(q.Webhooks[:i], q.Webhooks[i+1:]...)
				}
			} else {
				q.Webhooks = append(q.Webhooks[:i], q.Webhooks[i+1:]...)
			}
		}
	}
}
