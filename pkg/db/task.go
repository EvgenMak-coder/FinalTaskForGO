package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Добавить задачу
func AddTask(task *Task) (int64, error) {
	q := `INSERT INTO scheduler(date, title, comment, repeat) VALUES(?,?,?,?)`
	res, err := DB.Exec(q, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// Получить список задач
func Tasks(limit int) ([]*Task, error) {
	rows, err := DB.Query(`SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT ?`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	tasks := []*Task{}
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}
	if tasks == nil {
		tasks = []*Task{}
	}
	return tasks, nil
}

// Получить одну задачу
func GetTask(id string) (*Task, error) {
	row := DB.QueryRow(`SELECT id, date, title, comment, repeat FROM scheduler WHERE id=?`, id)
	var t Task
	if err := row.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("task not found")
		}
		return nil, err
	}
	return &t, nil
}

// Обновить задачу
func UpdateTask(task *Task) error {
	res, err := DB.Exec(`UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`,
		task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("task not found")
	}
	return nil
}

// Удалить задачу
func DeleteTask(id string) error {
	res, err := DB.Exec(`DELETE FROM scheduler WHERE id=?`, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return errors.New("task not found")
	}
	return nil
}

// Обновить только дату для задачи
func UpdateDate(next string, id string) error {
	res, err := DB.Exec(`UPDATE scheduler SET date=? WHERE id=?`, next, id)
	if err != nil {
		return err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if n == 0 {
		return errors.New("task not found")
	}
	return nil
}
