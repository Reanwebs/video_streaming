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
		uploadData, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		request = domain.ToSaveVideo{
			UserName:    uploadData.UserName,
			AvatarId:    uploadData.AvatarId,
			Title:       uploadData.Title,
			Discription: uploadData.Discription,
			Intrest:     uploadData.Intrest,
			ThumbnailId: uploadData.ThumbnailId,
		}

		_, err = buffer.Write(uploadData.Data)
		if err != nil {
			return err
		}
	}

	err := utils.UploadVideoToS3(buffer.Bytes(), s3Path)
	if err != nil {
		return err
	}

	request.S3Path = s3Path

	_, err = c.Repo.CreateVideoid(request)
	if err != nil {
		return err
	}

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
			VideoId:     uint32(v.ID),
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

func (c *VideoServer) FindArchivedVideoByUserId(ctx context.Context, input *pb.FindArchivedVideoByUserIdRequest) (*pb.FindArchivedVideoByUserIdResponse, error) {
	res, err := c.Repo.FindArchivedVideos(input.UserName)
	if err != nil {
		return nil, err
	}

	data := make([]*pb.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &pb.FetchVideo{
			VideoId:     uint32(v.ID),
			AvatarId:    v.Avatar_id,
			S3Path:      v.S3_path,
			UserName:    v.User_name,
			ThumbnailId: v.Thumbnail_id,
			Title:       v.Title,
			Intrest:     v.Interest,
		}
	}

	response := &pb.FindArchivedVideoByUserIdResponse{
		Videos: data,
	}
	return response, err
}

func (c *VideoServer) ArchiveVideo(ctx context.Context, input *pb.ArchiveVideoRequest) (*pb.ArchiveVideoResponse, error) {

	res, err := c.Repo.ArchivedVideos(uint(input.VideoId))
	if err != nil {
		return nil, err
	}
	response := &pb.ArchiveVideoResponse{
		Status: res,
	}
	return response, err

}
