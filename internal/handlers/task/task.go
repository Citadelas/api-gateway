package task

import (
	"github.com/Citadelas/api-gateway/internal/helpers/grpc"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"strconv"
	"time"
)

type Task struct {
	Id          uint64    `json:"id,omitempty"`
	UserId      uint64    `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	DueDate     time.Time `json:"due_date"`
}

func CreateTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.Create"
		var req Task
		uid, _ := c.Get("userID")
		c.ShouldBind(&req)
		log.Info("aa", req)
		grpcReq := taskv1.CreateTaskRequest{
			UserId:      uid.(uint64),
			Title:       req.Title,
			Description: req.Description,
			Priority:    taskv1.TaskPriority(taskv1.TaskPriority_value[req.Priority]),
			DueDate:     timestamppb.New(req.DueDate),
		}
		resp, err := client.CreateTask(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}

func UpdateTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.Update"
		var req Task
		uid, _ := c.Get("userID")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		c.ShouldBind(&req)
		grpcReq := taskv1.UpdateTaskRequest{
			Id:          id,
			UserId:      uid.(uint64),
			Title:       req.Title,
			Description: req.Description,
			Priority:    taskv1.TaskPriority(taskv1.TaskPriority_value[req.Priority]),
			DueDate:     timestamppb.New(req.DueDate),
		}
		resp, err := client.UpdateTask(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}

func GetTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.task.Get"
		uid, _ := c.Get("userID")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		grpcReq := taskv1.GetTaskRequest{
			Id:     id,
			UserId: uid.(uint64),
		}
		resp, err := client.GetTask(c, &grpcReq)
		if err != nil {
			grpc.HandleGRPCError(c, err)
			return
		}
		c.JSON(200, resp)
	}
}
