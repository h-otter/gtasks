package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	lerrors "github.com/syndtr/goleveldb/leveldb/errors"
)

type Task struct {
	ID          int               `json:"id"`
	State       string            `json:"state"`
	Labels      map[string]string `json:"labels"`
	Description string            `json:"description"`
	DueDate     string            `json:"due_date"`
	DependsOn   []int             `json:"depends_on"`
}

type TaskHandler struct {
	db *leveldb.DB
}

func CreateTaskHandler(directory string) (*TaskHandler, error) {
	db, err := leveldb.OpenFile(directory, nil)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect database")
	}

	return &TaskHandler{
		db: db,
	}, nil
}

func (h *TaskHandler) ListTasks() echo.HandlerFunc {
	return func(c echo.Context) error {
		res := make([]*Task, 0)

		iter := h.db.NewIterator(nil, nil)
		for iter.Next() {
			t := &Task{}
			if err := json.Unmarshal(iter.Value(), t); err != nil {
				return c.JSON(http.StatusInternalServerError, "failed to unmarshal")
			}

			res = append(res, t)
		}
		iter.Release()

		if err := iter.Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		if len(res) == 0 {
			return c.NoContent(http.StatusNotFound)
		}

		return c.JSON(http.StatusOK, res)
	}
}
func (h *TaskHandler) GetTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		b, err := h.db.Get([]byte(id), nil)
		if err != nil {
			if err == lerrors.ErrNotFound {
				return c.NoContent(http.StatusNotFound)
			}

			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		t := &Task{}
		if err := json.Unmarshal(b, t); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to unmarshal")
		}

		return c.JSON(http.StatusOK, t)
	}
}
func (h *TaskHandler) CreateTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := 0

		// keyはソートされているため、最後を取ればいいはず
		// TODO: lock
		iter := h.db.NewIterator(nil, nil)
		for iter.Next() {
			if iter.Last() {
				sid := string(iter.Key())
				id, _ = strconv.Atoi(sid)
				id++
			}
		}
		iter.Release()

		if err := iter.Error(); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		t := &Task{}
		if err := c.Bind(t); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}
		t.ID = id

		b, err := json.Marshal(t)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to marshal")
		}

		if err := h.db.Put([]byte(strconv.Itoa(id)), b, nil); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		return c.JSON(http.StatusOK, t)
	}
}
func (h *TaskHandler) UpdateTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		i, err := strconv.Atoi(id)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "id is uint")
		}

		t := &Task{}
		if err := c.Bind(t); err != nil {
			return c.JSON(http.StatusBadRequest, "")
		}
		t.ID = i

		b, err := json.Marshal(t)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to marshal")
		}

		if err := h.db.Put([]byte(id), b, nil); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		return c.JSON(http.StatusOK, t)
	}
}
func (h *TaskHandler) DeleteTask() echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")
		if err := h.db.Delete([]byte(id), nil); err != nil {
			return c.JSON(http.StatusInternalServerError, "failed to operate on db")
		}

		return c.NoContent(http.StatusOK)
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	h, err := CreateTaskHandler("db")
	if err != nil {
		panic(err)
	}
	e.GET("/api/task/", h.ListTasks())
	e.GET("/api/task/:id", h.GetTask())
	e.POST("/api/task/", h.CreateTask())
	e.PUT("/api/task/:id", h.UpdateTask())
	e.DELETE("/api/task/:id", h.DeleteTask())

	e.Start(":2345")
}
