package timecraft

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	v1 "github.com/stealthrocket/timecraft/gen/proto/go/timecraft/server/v1"
	"github.com/stealthrocket/timecraft/gen/proto/go/timecraft/server/v1/serverv1connect"
)

// NewClient creates a timecraft client.
func NewClient() (*Client, error) {
	grpcClient := serverv1connect.NewTimecraftServiceClient(
		httpClient,
		"http://timecraft/",
		connect.WithAcceptCompression("gzip", nil, nil), // disable gzip for now
		connect.WithProtoJSON(),                         // use JSON for now
		// connect.WithCodec(grpc.Codec{}),
	)

	return &Client{
		grpcClient: grpcClient,
	}, nil
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		DialContext: dialContext,
		// TODO: timeouts/limits
	},
}

// Client is a timecraft client.
type Client struct {
	grpcClient serverv1connect.TimecraftServiceClient
}

// SubmitTasks submits tasks to the timecraft runtime.
//
// The tasks are executed asynchronously. The method returns a set of TaskID
// that can be used to query task status and fetch task output when ready.
func (c *Client) SubmitTasks(ctx context.Context, requests []TaskRequest) ([]TaskID, error) {
	r := connect.NewRequest(&v1.SubmitTasksRequest{
		Requests: make([]*v1.TaskRequest, len(requests)),
	})
	for i := range requests {
		var err error
		r.Msg.Requests[i], err = c.makeTaskRequest(&requests[i])
		if err != nil {
			return nil, err
		}
	}
	res, err := c.grpcClient.SubmitTasks(ctx, r)
	if err != nil {
		return nil, err
	}
	taskIDs := make([]TaskID, len(res.Msg.TaskId))
	for i, taskID := range res.Msg.TaskId {
		taskIDs[i] = TaskID(taskID)
	}
	return taskIDs, nil
}

// LookupTasks retrieves task responses by ID.
func (c *Client) LookupTasks(ctx context.Context, taskIDs []TaskID) ([]TaskResponse, error) {
	req := connect.NewRequest(&v1.LookupTasksRequest{
		TaskId: make([]string, len(taskIDs)),
	})
	for i, taskID := range taskIDs {
		req.Msg.TaskId[i] = string(taskID)
	}
	res, err := c.grpcClient.LookupTasks(ctx, req)
	if err != nil {
		return nil, err
	}
	responses := make([]TaskResponse, len(res.Msg.Responses))
	for i, taskResponse := range res.Msg.Responses {
		responses[i], err = c.makeTaskResponse(taskResponse)
		if err != nil {
			return nil, err
		}
	}
	return responses, nil
}

// PollTasks retrieves tasks that are complete.
//
// Tasks are complete when task execution either succeeds or fails permanently.
//
// PollTasks will block the goroutine until either batchSize tasks are
// complete, or the timeout is reached, whichever comes first. If timeout
// is zero, the timecraft runtime won't block waiting for complete tasks.
// If the timeout is less than zero, the method will block until batchSize
// tasks are complete (and thus may block indefinitely).
func (c *Client) PollTasks(ctx context.Context, batchSize int, timeout time.Duration) ([]TaskResponse, error) {
	req := connect.NewRequest(&v1.PollTasksRequest{
		BatchSize: int32(batchSize),
		TimeoutNs: int64(timeout),
	})
	res, err := c.grpcClient.PollTasks(ctx, req)
	if err != nil {
		return nil, err
	}
	responses := make([]TaskResponse, len(res.Msg.Responses))
	for i, taskResponse := range res.Msg.Responses {
		responses[i], err = c.makeTaskResponse(taskResponse)
		if err != nil {
			return nil, err
		}
	}
	return responses, nil
}

// DiscardTasks discards a batch of tasks by ID.
func (c *Client) DiscardTasks(ctx context.Context, taskIDs []TaskID) error {
	req := connect.NewRequest(&v1.DiscardTasksRequest{
		TaskId: make([]string, len(taskIDs)),
	})
	for i, taskID := range taskIDs {
		req.Msg.TaskId[i] = string(taskID)
	}
	_, err := c.grpcClient.DiscardTasks(ctx, req)
	return err
}

func (c *Client) makeTaskRequest(req *TaskRequest) (*v1.TaskRequest, error) {
	r := &v1.TaskRequest{
		Module: &v1.ModuleSpec{Path: req.Module.Path, Args: req.Module.Args},
	}
	switch in := req.Input.(type) {
	case *HTTPRequest:
		headers := make([]*v1.Header, 0, len(in.Headers))
		for name, values := range in.Headers {
			for _, value := range values {
				headers = append(headers, &v1.Header{Name: name, Value: value})
			}
		}
		r.Input = &v1.TaskRequest_HttpRequest{HttpRequest: &v1.HTTPRequest{
			Method:  in.Method,
			Path:    in.Path,
			Body:    in.Body,
			Headers: headers,
		}}
	default:
		return nil, fmt.Errorf("invalid task input: %v", req.Input)
	}
	return r, nil
}

func (c *Client) makeTaskResponse(res *v1.TaskResponse) (TaskResponse, error) {
	taskResponse := TaskResponse{
		State:     TaskState(res.State),
		ProcessID: ProcessID(res.ProcessId),
	}
	if taskResponse.State == Error {
		taskResponse.Error = errors.New(res.ErrorMessage)
	}
	switch out := res.Output.(type) {
	case *v1.TaskResponse_HttpResponse:
		httpResponse := &HTTPResponse{
			StatusCode: int(out.HttpResponse.StatusCode),
			Body:       out.HttpResponse.Body,
			Headers:    make(http.Header, len(out.HttpResponse.Headers)),
		}
		for _, h := range out.HttpResponse.Headers {
			httpResponse.Headers[h.Name] = append(httpResponse.Headers[h.Name], h.Value)
		}
		taskResponse.Output = httpResponse
	}
	return taskResponse, nil
}

// Version fetches the timecraft version.
func (c *Client) Version(ctx context.Context) (string, error) {
	req := connect.NewRequest(&v1.VersionRequest{})
	res, err := c.grpcClient.Version(ctx, req)
	if err != nil {
		return "", err
	}
	return res.Msg.Version, nil
}