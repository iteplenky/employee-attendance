package attendance

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/iteplenky/employee-attendance/database"
	"github.com/iteplenky/employee-attendance/internal/bot"
	"log"
	"strconv"
	"time"
)

type Event struct {
	ID            int    `json:"id"`
	EmpID         string `json:"emp_id"`
	PunchTime     string `json:"punch_time"`
	TerminalAlias string `json:"terminal_alias"`
	Processed     bool   `json:"processed"`
}

func HandleAttendanceEvents(ctx context.Context, cache database.Cache, b *bot.Bot) {
	ch := cache.Subscribe(ctx, "attendance_events")
	for msg := range ch {
		log.Printf("received message: %v", msg)

		var event Event
		if err := json.Unmarshal([]byte(msg), &event); err != nil {
			log.Printf("failed to unmarshal message: %v", err)
			continue
		}

		members, err := cache.HGetAll(ctx, "subscribed_users")
		if err != nil {
			log.Printf("failed to get subscribed users: %v", err)
			continue
		}

		tgID, exists := members[event.EmpID]
		if !exists {
			continue
		}

		tg, err := strconv.ParseInt(tgID, 10, 64)
		if err != nil {
			log.Printf("failed to convert tgID to int: %v", err)
			continue
		}

		if _, err = b.Bot.SendMessage(tg, FormatAttendanceMessage(event), nil); err != nil {
			log.Printf("failed to send message: %v", err)
		}
	}
}

func FormatAttendanceMessage(event Event) string {
	return fmt.Sprintf(
		"✅ Вы отметились: \n⏰ Время: %s\n🏢 Терминал: %s",
		formatTime(event.PunchTime), event.TerminalAlias,
	)
}

func formatTime(punchTime string) string {
	t, err := time.Parse(time.RFC3339Nano, punchTime)
	if err != nil {
		log.Printf("failed to parse time: %v", err)
		return punchTime
	}
	return t.Format("15:04")
}
