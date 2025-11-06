package main

import (
	"context"
	"net/http"
	"os"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/gin-gonic/gin"
)

type StudentCreate struct {
	FirstName string    `json:"first_name" binding:"required"`
	LastName  string    `json:"last_name" binding:"required"`
	Phone     string    `json:"phone" binding:"required"`
	Email     string    `json:"email" binding:"required"`
	BirthDate time.Time `json:"birth_date" binding:"required"`
}

type Student struct {
	StudentID string    `json:"student_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Phone     string    `json:"phone"`
	Email     string    `json:"email"`
	BirthDate time.Time `json:"birth_date"`
	CreatedAt time.Time `json:"created_at"`
}

type EnrollmentCreate struct {
	StudentID string `json:"student_id" binding:"required"`
	CourseID  string `json:"course_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

type Course struct {
	CourseID     string    `json:"course_id"`
	Title        string    `json:"title"`
	Level        string    `json:"level"`
	InstrumentID string    `json:"instrument_id"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`
}

func getCH() (driver.Conn, error) {
	host := getenv("CLICKHOUSE_HOST", "localhost")
	port := getenv("CLICKHOUSE_PORT", "9000")
	user := getenv("CLICKHOUSE_USER", "default")
	pass := getenv("CLICKHOUSE_PASSWORD", "")
	db := getenv("CLICKHOUSE_DATABASE", "music_school")

	opts := &ch.Options{
		Addr: []string{host + ":" + port},
		Auth: ch.Auth{
			Database: db,
			Username: user,
			Password: pass,
		},
		DialTimeout: 5 * time.Second,
		Settings: ch.Settings{
			"max_execution_time": 60,
		},
	}
	return ch.Open(opts)
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		conn, err := getCH()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": err.Error()})
			return
		}
		defer conn.Close()
		ctx := context.Background()
		var version string
		if err := conn.QueryRow(ctx, "SELECT version()").Scan(&version); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "error", "error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ok", "clickhouse_version": version})
	})

	r.GET("/students", func(c *gin.Context) {
		conn, err := getCH()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()
		ctx := context.Background()
		rows, err := conn.Query(ctx, `
            SELECT student_id, first_name, last_name, phone, email, birth_date, created_at
            FROM students
            ORDER BY created_at DESC
            LIMIT 100`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var list []Student
		for rows.Next() {
			var s Student
			if err := rows.Scan(&s.StudentID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &s.BirthDate, &s.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			list = append(list, s)
		}
		c.JSON(http.StatusOK, list)
	})

	r.POST("/students", func(c *gin.Context) {
		var in StudentCreate
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		conn, err := getCH()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()
		ctx := context.Background()
		// Insert
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO students (first_name, last_name, phone, email, birth_date)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := batch.Append(in.FirstName, in.LastName, in.Phone, in.Email, in.BirthDate); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := batch.Send(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		// Read back
		var s Student
		if err := conn.QueryRow(ctx, `
            SELECT student_id, first_name, last_name, phone, email, birth_date, created_at
            FROM students
            ORDER BY created_at DESC
            LIMIT 1`).Scan(&s.StudentID, &s.FirstName, &s.LastName, &s.Phone, &s.Email, &s.BirthDate, &s.CreatedAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, s)
	})

	r.GET("/courses", func(c *gin.Context) {
		conn, err := getCH()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()
		ctx := context.Background()
		rows, err := conn.Query(ctx, `
            SELECT course_id, title, level, instrument_id, description, created_at
            FROM courses
            ORDER BY created_at DESC
            LIMIT 100`)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer rows.Close()
		var list []Course
		for rows.Next() {
			var x Course
			if err := rows.Scan(&x.CourseID, &x.Title, &x.Level, &x.InstrumentID, &x.Description, &x.CreatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			list = append(list, x)
		}
		c.JSON(http.StatusOK, list)
	})

	r.POST("/enrollments", func(c *gin.Context) {
		var in EnrollmentCreate
		if err := c.ShouldBindJSON(&in); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		conn, err := getCH()
		if err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}
		defer conn.Close()
		ctx := context.Background()
		var sc, cc uint64
		if err := conn.QueryRow(ctx, "SELECT count() FROM students WHERE student_id = ?", in.StudentID).Scan(&sc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := conn.QueryRow(ctx, "SELECT count() FROM courses WHERE course_id = ?", in.CourseID).Scan(&cc); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if sc == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "student not found"})
			return
		}
		if cc == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "course not found"})
			return
		}
		batch, err := conn.PrepareBatch(ctx, "INSERT INTO enrollments (student_id, course_id, status)")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := batch.Append(in.StudentID, in.CourseID, in.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if err := batch.Send(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"ok": true})
	})

	if err := r.Run(":8000"); err != nil {
		panic(err)
	}
}
