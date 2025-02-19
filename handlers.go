package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func createTask(c *fiber.Ctx) error {
	// Выводим сырое тело запроса для отладки
	rawBody := c.Body()
	fmt.Println("Raw Body:", string(rawBody))

	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		// Выводим ошибку парсинга
		fmt.Println("BodyParser Error:", err)
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Выводим распарсенные данные
	fmt.Printf("Parsed Task: %+v\n", task)

	// Остальной код для создания задачи
	db, err := connectDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer db.Close()

	query := `INSERT INTO tasks (title, description, status) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`
	err = db.QueryRow(query, task.Title, task.Description, task.Status).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func getTasks(c *fiber.Ctx) error {
	db, err := connectDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer db.Close()

	rows, err := db.Query("SELECT id, title, description, status, created_at, updated_at FROM tasks")
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}
		tasks = append(tasks, task)
	}

	return c.JSON(tasks)
}

func updateTask(c *fiber.Ctx) error {
	db, err := connectDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer db.Close()

	id := c.Params("id")
	task := new(Task)
	if err := c.BodyParser(task); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	query := `UPDATE tasks SET title=$1, description=$2, status=$3, updated_at=now() WHERE id=$4 RETURNING updated_at`
	err = db.QueryRow(query, task.Title, task.Description, task.Status, id).Scan(&task.UpdatedAt)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(task)
}

func deleteTask(c *fiber.Ctx) error {
	db, err := connectDB()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	defer db.Close()

	id := c.Params("id")
	_, err = db.Exec("DELETE FROM tasks WHERE id=$1", id)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	c.SendStatus(204) // Отправляем статус 204 No Content
	return nil        // Возвращаем nil, так как ошибок нет
}
