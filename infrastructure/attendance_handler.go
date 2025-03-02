package infrastructure

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iteplenky/employee-attendance/domain"
	"log"
	"strconv"
	"time"

	"github.com/iteplenky/employee-attendance/internal/bot"
)

func HandleAttendanceEvents(ctx context.Context, cache *RedisCache, b *bot.Bot) {
	ch := cache.Subscribe(ctx, "attendance_events")
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				log.Printf("attendance_events channel closed")
				return
			}
			log.Println("Received event:", msg)

			var event domain.AttendanceEvent
			if err := json.Unmarshal([]byte(msg), &event); err != nil {
				log.Printf("failed to unmarshal event: %v\n", err)
				continue
			}

			members, err := cache.HGetAll(context.Background(), "subscribed_users")
			if err != nil {
				log.Printf("failed to fetch subscribed_users: %v\n", err)
				continue
			}

			user, exists := members[event.IIN]
			if !exists {
				continue
			}

			tgID, err := strconv.ParseInt(user, 10, 64)
			if err != nil {
				log.Printf("failed to convert tgID to int: %v", err)
				continue
			}

			if _, err = b.Bot.SendMessage(tgID, formatAttendanceMessage(event), nil); err != nil {
				log.Printf("failed to send message: %v\n", err)
			}
		case <-ctx.Done():
			log.Println("Stopping attendance event handler...")
			return
		}
	}
}

func formatAttendanceMessage(event domain.AttendanceEvent) string {
	return fmt.Sprintf(
		"✨ Вы отметились в %s у терминала %s",
		formatTime(event.PunchTime), event.TerminalAlias,
	)
}

func formatTime(punchTime string) string {
	t, err := time.Parse("2006-01-02T15:04:05.999999", punchTime)
	if err != nil {
		log.Printf("failed to parse time: %v", err)
		return punchTime
	}
	return t.Format("15:04")
}
