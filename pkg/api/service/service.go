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
			Interest:    uploadData.Intrest,
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
	videoID, err := c.Repo.CreateVideoid(request)
	if err != nil {
		// Handle the error
		return err
	}

	// Sending a response and closing the sending stream of bytes
	return stream.SendAndClose(&pb.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video successfully uploaded.",
		VideoId: videoID,
	})
}

// func (c *VideoServer) UploadVideo(stream pb.VideoService_UploadVideoServer) error {

// 	var uploadData *pb.UploadVideoRequest
// 	var buffer bytes.Buffer

// 	uploadData, err := stream.Recv()
// 	if err != nil {
// 		return err
// 	}

// 	for {
// 		chunk, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		_, err = buffer.Write(chunk.Data)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	fileUID := uuid.New()
// 	fileName := fileUID.String()
// 	s3Path := "reanweb/" + fileName + ".mp4"

// 	err = utils.UploadVideoToS3(buffer.Bytes(), s3Path)
// 	if err != nil {
// 		fmt.Println("Error uploading video to S3:", err)
// 		if awsErr, ok := err.(awserr.Error); ok {
// 			fmt.Println("AWS Error Code:", awsErr.Code())
// 			fmt.Println("AWS Error Message:", awsErr.Message())
// 		}
// 		return err
// 	}
// 	request := domain.ToSaveVideo{
// 		S3Path:      s3Path,
// 		UserName:    uploadData.UserName,
// 		AvatarId:    uploadData.AvatarId,
// 		Title:       uploadData.ThumbnailId,
// 		Discription: uploadData.Discription,
// 		Interest:    uploadData.Intrest,
// 		ThumbnailId: uploadData.ThumbnailId,
// 	}

// 	videoID, err := c.Repo.CreateVideoid(request)
// 	if err != nil {
// 		fmt.Println("Error saving video ID:", err)
// 		return err
// 	}

// 	fmt.Println("Video ID:", videoID)

// 	return stream.SendAndClose(&pb.UploadVideoResponse{
// 		Status:  http.StatusOK,
// 		Message: "Video successfully uploaded.",
// 		VideoId: videoID,
// 	})
// }

// func (c *VideoServer) UploadVideo(stream pb.VideoService_UploadVideoServer) error {
// 	var uploadData *pb.UploadVideoRequest
// 	fmt.Println("\nuploaded data struct empty\n\n", uploadData, "\n\n\n\n.")
// 	request := domain.ToSaveVideo{
// 		UserName:    uploadData.UserName,
// 		AvatarId:    uploadData.AvatarId,
// 		Title:       uploadData.ThumbnailId,
// 		Discription: uploadData.Discription,
// 		Interest:    uploadData.Intrest,
// 		ThumbnailId: uploadData.ThumbnailId,
// 	}
// 	var buffer bytes.Buffer
// 	for {
// 		chunk, err := stream.Recv()
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			return err
// 		}
// 		_, err = buffer.Write(chunk.Data)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	fileUID := uuid.New()
// 	fileName := fileUID.String()
// 	s3Path := "reanweb/" + fileName + ".mp4"

// 	err := utils.UploadVideoToS3(buffer.Bytes(), s3Path)
// 	if err != nil {
// 		fmt.Println("Error uploading video to S3:", err)
// 		if awsErr, ok := err.(awserr.Error); ok {
// 			fmt.Println("AWS Error Code:", awsErr.Code())
// 			fmt.Println("AWS Error Message:", awsErr.Message())
// 		}
// 		return err
// 	}
// 	request.S3Path = s3Path
// 	videoID, err := c.Repo.CreateVideoid(request)
// 	if err != nil {
// 		fmt.Println("Error saving video ID:", err)
// 		return err
// 	}

// 	fmt.Println("Video ID:", videoID)

// 	// Sending a response and closing the sending stream of bytes
// 	return stream.SendAndClose(&pb.UploadVideoResponse{
// 		Status:  http.StatusOK,
// 		Message: "Video successfully uploaded.",
// 		VideoId: videoID,
// 	})
// }
