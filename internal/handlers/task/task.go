package task

import (
	"github.com/Citadelas/api-gateway/internal/helpers/grpc"
	taskv1 "github.com/Citadelas/protos/golang/task"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
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
		const op = "handlers.sso.Login"
		var req Task
		c.ShouldBind(&req)
		log.Info("aa", req)
		grpcReq := taskv1.CreateTaskRequest{
			UserId:      req.UserId,
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
