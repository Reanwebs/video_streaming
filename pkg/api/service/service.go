package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
	clientinterfaces "videoStreaming/pkg/client/clientInterfaces"
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb"
	"videoStreaming/pkg/respository/interfaces"
	"videoStreaming/pkg/utils"

	"github.com/google/uuid"
)

type VideoServer struct {
	Repo        interfaces.VideoRepo
	MonitClient clientinterfaces.MonitClient
	pb.VideoServiceServer
}

func NewVideoServer(repo interfaces.VideoRepo, monitClient clientinterfaces.MonitClient) pb.VideoServiceServer {
	return &VideoServer{
		Repo:        repo,
		MonitClient: monitClient,
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

		videoId := utils.GenerateUniqueString()

		request = domain.ToSaveVideo{
			UserName:    uploadData.UserName,
			AvatarId:    uploadData.AvatarId,
			Title:       uploadData.Title,
			Discription: uploadData.Discription,
			Intrest:     uploadData.Intrest,
			ThumbnailId: uploadData.ThumbnailId,
			Video_id:    videoId,
			UserId:      uploadData.UserId,
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
			VideoId:     v.Video_id,
			AvatarId:    v.Avatar_id,
			S3Path:      v.S3_path,
			UserName:    v.User_name,
			ThumbnailId: v.Thumbnail_id,
			Title:       v.Title,
			Intrest:     v.Interest,
			Archived:    v.Archived,
			Views:       uint32(v.Views),
			Starred:     uint32(v.Starred),
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
			VideoId:     v.Video_id,
			AvatarId:    v.Avatar_id,
			S3Path:      v.S3_path,
			UserName:    v.User_name,
			ThumbnailId: v.Thumbnail_id,
			Title:       v.Title,
			Intrest:     v.Interest,
			Archived:    v.Archived,
			Views:       uint32(v.Views),
			Starred:     uint32(v.Starred),
		}
	}

	response := &pb.FindArchivedVideoByUserIdResponse{
		Videos: data,
	}
	return response, err
}

func (c *VideoServer) FetchAllVideo(ctx context.Context, input *pb.FetchAllVideoRequest) (*pb.FetchAllVideoResponse, error) {

	res, err := c.Repo.FetchAllVideos()
	if err != nil {
		return nil, err
	}
	data := make([]*pb.FetchVideo, len(res))
	for i, v := range res {
		data[i] = &pb.FetchVideo{
			VideoId:     v.Video_id,
			AvatarId:    v.Avatar_id,
			S3Path:      v.S3_path,
			UserName:    v.User_name,
			ThumbnailId: v.Thumbnail_id,
			Title:       v.Title,
			Intrest:     v.Interest,
			Archived:    v.Archived,
			Views:       uint32(v.Views),
			Starred:     uint32(v.Starred),
		}
	}

	response := &pb.FetchAllVideoResponse{
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

func (c *VideoServer) GetVideoById(ctx context.Context, input *pb.GetVideoByIdRequest) (*pb.GetVideoByIdResponse, error) {

	res, isStarred, err := c.Repo.GetVideoById(input.VideoId, input.UserName)
	if err != nil {
		return nil, err
	}
	response := &pb.GetVideoByIdResponse{
		VideoId:     res.Video_id,
		UserName:    res.User_name,
		AvatarId:    res.Avatar_id,
		Archived:    res.Archived,
		Intrest:     res.Interest,
		ThumbnailId: res.Thumbnail_id,
		Title:       res.Title,
		S3Path:      res.S3_path,
		IsStarred:   isStarred,
		Views:       uint32(res.Views),
		Starred:     uint32(res.Starred),
	}
	if res.Views >= 100 {
		reward := res.Views % 100
		if reward == 0 {
			err := c.MonitClient.VideoReward(ctx, domain.VideoRewardRequest{
				UserID:    res.UserId,
				VideoID:   response.VideoId,
				Reason:    "views",
				Views:     response.Views,
				PaidCoins: 0,
			})
			if err != nil {
				fmt.Println(" failed to add reward")
			}
		}
	}
	return response, err

}

func (c *VideoServer) ToggleStar(ctx context.Context, input *pb.ToggleStarRequest) (*pb.ToggleStarResponse, error) {

	res, err := c.Repo.ToggleStar(input.VideoId, input.UserNAme, input.Starred)
	if err != nil {
		return nil, err
	}

	response := &pb.ToggleStarResponse{
		Status: res,
	}

	return response, nil
}

func (c *VideoServer) BlockVideo(ctx context.Context, input *pb.BlockVideoRequest) (*pb.BlockVideoResponse, error) {
	res, err := c.Repo.BlockVideo(domain.BlockedVideo{
		VideoID:   input.VideoId,
		Reason:    input.Reason,
		Timestamp: time.Now(),
	})
	if err != nil {
		return nil, err
	}

	response := &pb.BlockVideoResponse{
		Blocked: res,
	}

	return response, nil
}
