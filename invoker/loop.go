package invoker

import (
	"context"
	"task-scheduler/data"
	"task-scheduler/database"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jasonlvhit/gocron"
)

var queries *database.Queries
var queue *data.ScheduleQueue

func New(conn *pgx.Conn) {
	queries = database.New(conn)

	gocron.Every(10).Minute().Do(SecondaryLoop)

	gocron.Every(1).Minute().Do(PrimaryLoop)
}

func PrimaryLoop() {
	toCheck, err := queue.Peek()

	if err != nil {
		return
	}

	isAfter := toCheck.InvocationTimestamp.Time.After(time.Now())

	if !isAfter {
		return
	}

	schedule, _ := queue.Dequeue()

	go invoke(schedule)
}

func SecondaryLoop() {
	ctx := context.Background()

	// fetch most recent schedule
	schedules, err := queries.ListSchedule(ctx, 10)

	if err != nil {
		return
	}

	queue := data.GetQueueInstance()
	queue.SetQueue(schedules)
}
