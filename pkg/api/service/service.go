package service

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb"
	"videoStreaming/pkg/respository/interfaces"
	"videoStreaming/pkg/utils"

	"github.com/google/uuid"
)

const storageLocation = "storage"

type VideoServer struct {
	Repo interfaces.VideoRepo
	pb.VideoServiceServer
}

func NewVideoServer(repo interfaces.VideoRepo) pb.VideoServiceServer {
	return &VideoServer{
		Repo: repo,
	}
}

func (c *VideoServer) HealthCheck(ctx context.Context, input *pb.Request) (*pb.Response, error) {

	if input == nil {
		return nil, errors.New("empty input")
	}
	a := &pb.Response{
		Result: "success",
	}
	return a, nil
}

func (c *VideoServer) UploadVideo(stream pb.VideoService_UploadVideoServer) error {
	var request domain.ToSaveVideo
	var buffer bytes.Buffer
	fileUID := uuid.New()
	fileName := fileUID.String()
	s3Path := "reanweb/" + fileName + ".mp4"

	for {
		// Receive the next message from the client
		uploadData, err := stream.Recv()
		if err == io.EOF {
			break // End of stream
		}
		if err != nil {
			return err
		}

		// Process the received message and populate the request struct
		request = domain.ToSaveVideo{
			UserName:    uploadData.UserName,
			AvatarId:    uploadData.AvatarId,
			Title:       uploadData.Title,
			Discription: uploadData.Discription,
			Intrest:     uploadData.Intrest,
			ThumbnailId: uploadData.ThumbnailId,
		}

		// Append the received data to the buffer
		_, err = buffer.Write(uploadData.Data)
		if err != nil {
			return err
		}
	}

	// Upload the video data to S3
	err := utils.UploadVideoToS3(buffer.Bytes(), s3Path)
	if err != nil {
		// Handle the error
		return err
	}

	// Set the S3Path in the request
	request.S3Path = s3Path

	// Create the video record and get the video ID
	_, err = c.Repo.CreateVideoid(request)
	if err != nil {
		// Handle the error
		return err
	}

	// Sending a response and closing the sending stream of bytes
	return stream.SendAndClose(&pb.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video successfully uploaded.",
		VideoId: "",
	})
}

func (c *VideoServer) FindUserVideo(ctx context.Context, input *pb.FindUserVideoRequest) (*pb.FindUserVideoResponse, error) {
	res, err := c.Repo.FetchUserVideos(input.UserName)
	if err != nil {
		return nil, err
	}

	data := make([]*pb.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &pb.FetchVideo{
			AvatarId:    v.Avatar_id,
			S3Path:      v.S3_path,
			UserName:    v.User_name,
			ThumbnailId: v.Thumbnail_id,
			Title:       v.Title,
			Intrest:     v.Interest,
		}
	}

	response := &pb.FindUserVideoResponse{
		Videos: data,
	}
	return response, err
}
