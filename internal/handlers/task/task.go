package task

import (
	"github.com/Citadelas/api-gateway/internal/helpers/grpc"
	"github.com/Citadelas/api-gateway/internal/lib/logger/sl"
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
	const op = "handlers.task.Create"
	log = log.With("op", op)
	return func(c *gin.Context) {
		var req Task
		uid, _ := c.Get("userID")
		c.ShouldBind(&req)
		grpcReq := taskv1.CreateTaskRequest{
			UserId:      uid.(uint64),
			Title:       req.Title,
			Description: req.Description,
			Priority:    taskv1.TaskPriority(taskv1.TaskPriority_value[req.Priority]),
			DueDate:     timestamppb.New(req.DueDate),
		}
		resp, err := client.CreateTask(c, &grpcReq)
		if err != nil {
			log.Error("Error making grpc create task request", sl.Err(err))
			grpc.HandleGRPCError(c, err)
			return
		}
		log.Info("Task created successfully")
		c.JSON(200, resp)
	}
}

func UpdateTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	const op = "handlers.task.Update"
	log = log.With("op", op)
	return func(c *gin.Context) {
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
			log.Error("Error making grpc update task request", sl.Err(err))
			grpc.HandleGRPCError(c, err)
			return
		}
		log.Info("Task updated successfully")
		c.JSON(200, resp)
	}
}

func GetTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	const op = "handlers.task.Get"
	log = log.With("op", op)
	return func(c *gin.Context) {
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
			log.Error("Error making grpc get task request", sl.Err(err))
			grpc.HandleGRPCError(c, err)
			return
		}
		log.Info("Task get successfully")
		c.JSON(200, resp)
	}
}

func DeleteTaskHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	const op = "handlers.task.Delete"
	log = log.With("op", op)
	return func(c *gin.Context) {
		uid, _ := c.Get("userID")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		grpcReq := taskv1.DeleteTaskRequest{
			Id:     id,
			UserId: uid.(uint64),
		}
		resp, err := client.DeleteTask(c, &grpcReq)
		if err != nil {
			log.Error("Error making grpc delete task request", sl.Err(err))
			grpc.HandleGRPCError(c, err)
			return
		}
		log.Info("Task deleted successfully")
		c.JSON(200, resp)
	}
}

type UpdateStatusReq struct {
	Status string `json:"status"`
}

func UpdateStatusHandler(log *slog.Logger, client taskv1.TaskServiceClient) gin.HandlerFunc {
	const op = "handlers.task.UpdateStatus"
	log = log.With("op", op)
	return func(c *gin.Context) {
		var req UpdateStatusReq
		if err := c.ShouldBind(&req); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		uid, _ := c.Get("userID")
		id, err := strconv.ParseUint(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(400, gin.H{"error": "invalid id"})
			return
		}
		grpcReq := taskv1.UpdateStatusRequest{
			Id:     id,
			UserId: uid.(uint64),
			Status: taskv1.TaskStatus(taskv1.TaskStatus_value[req.Status]),
		}
		resp, err := client.UpdateStatus(c, &grpcReq)
		if err != nil {
			log.Error("Error making grpc update task status request", sl.Err(err))
			grpc.HandleGRPCError(c, err)
			return
		}
		log.Info("Task status updated successfully")
		c.JSON(200, resp)
	}
}
