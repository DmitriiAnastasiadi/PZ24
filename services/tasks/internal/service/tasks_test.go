package service

import (
	"testing"
)

func resetTasks() {
	mu.Lock()
	defer mu.Unlock()
	tasks = make(map[string]Task)
	counter = 0
}

func TestCreateAndGetTask(t *testing.T) {
	resetTasks()

	created := CreateTask("Test task", "Description", "2026-05-13")
	if created.ID == "" {
		t.Fatal("expected task ID to be set")
	}
	if created.Title != "Test task" {
		t.Fatalf("expected title %q, got %q", "Test task", created.Title)
	}

	got, ok := GetTask(created.ID)
	if !ok {
		t.Fatal("expected task to be found")
	}
	if got.Title != created.Title {
		t.Fatalf("expected fetched title %q, got %q", created.Title, got.Title)
	}
}

func TestUpdateAndDeleteTask(t *testing.T) {
	resetTasks()

	created := CreateTask("Initial", "Desc", "")

	updated, ok := UpdateTask(created.ID, map[string]interface{}{"title": "Updated", "done": true})
	if !ok {
		t.Fatal("expected update to succeed")
	}
	if updated.Title != "Updated" {
		t.Fatalf("expected updated title %q, got %q", "Updated", updated.Title)
	}
	if !updated.Done {
		t.Fatal("expected done flag to be true")
	}

	deleted := DeleteTask(created.ID)
	if !deleted {
		t.Fatal("expected task to be deleted")
	}
	_, ok = GetTask(created.ID)
	if ok {
		t.Fatal("expected deleted task to be absent")
	}
}

func TestGetTasksReturnsSummary(t *testing.T) {
	resetTasks()

	CreateTask("First", "", "")
	CreateTask("Second", "", "")

	list := GetTasks()
	if len(list) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(list))
	}
	if list[0].Description != "" || list[1].Description != "" {
		t.Fatal("expected summary list to omit full description")
	}
}
