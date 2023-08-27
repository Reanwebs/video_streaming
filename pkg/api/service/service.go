package service

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"videoStreaming/pkg/pb"
	"videoStreaming/pkg/respository/interfaces"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

const storageLocation = "storage"

type VideoServer struct {
	Repo interfaces.VideoRepo
	pb.UnimplementedVideoServiceServer
}

func NewVideoServer(repo interfaces.VideoRepo) pb.VideoServiceServer {
	return &VideoServer{
		Repo: repo,
	}
}

func (c *VideoServer) UploadVideo(stream pb.VideoService_UploadVideoServer) error {
	fmt.Println("service called")
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
	// Create a new AWS session with the loaded access keys
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-south-1"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
		LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return err
	}
	fmt.Println("\n\n\n", sess, "\n\n\n.")
	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	fileuid := uuid.New()
	fileName := fileuid.String()
	folderpath := storageLocation + "/" + fileName
	filepath := folderpath + "/" + fileName + ".mp4"

	if err := os.MkdirAll(folderpath, 0755); err != nil {
		return err
	}

	newfile, err1 := os.Create(filepath)
	if err1 != nil {
		return err1
	}
	defer newfile.Close()

	//receiving from the streamed bytes from the api_gateway
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if _, err := newfile.Write(chunk.Data); err != nil {
			return err
		}
	}

	s3Path := "reanweb/" + fileName + ".mp4"
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("reanweb"),
		Key:    aws.String(s3Path),
		Body:   newfile,
	})
	if err != nil {
		fmt.Println(
			"config.AWS_ACCESS_KEY_ID\t", os.Getenv("AWS_ACCESS_KEY_ID"),
			"\nconfig.AWS_SECRET_ACCESS_KEY\t", os.Getenv("AWS_SECRET_ACCESS_KEY"),
		)
		fmt.Println("error here")
		return err
	}

	chanerr := make(chan error, 2)
	go func() {
		//call to segment the file to hls format using ffmpeg
		err := CreatePlaylistAndSegments(filepath, folderpath)
		chanerr <- err
	}()
	go func() {
		//saving the video id to a database
		err := c.Repo.CreateVideoid(fileName)
		chanerr <- err
	}()
	for i := 1; i <= 2; i++ {
		err := <-chanerr
		if err != nil {
			return err
		}
	}

	// sending a response and closing the sending stream of bytes
	return stream.SendAndClose(&pb.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video successfully uploaded.",
		VideoId: fileName,
	})
}

func (c *VideoServer) StreamVideo(req *pb.StreamVideoRequest, stream pb.VideoService_StreamVideoServer) error {
	chunkSize := 4096 // Set your desired chunk size
	buffer := make([]byte, chunkSize)
	playlistPath := fmt.Sprintf("storage/%s/%s", req.Videoid, req.Playlist)
	plalistfile, _ := os.Open(playlistPath)
	defer plalistfile.Close()
	for {
		n, err := plalistfile.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		// Send the video chunk as a response to the client
		if err := stream.Send(&pb.StreamVideoResponse{
			VideoChunk: buffer[:n],
		}); err != nil {
			return err
		}
	}
	return nil
}

// to find all the video id
func (c *VideoServer) FindAllVideo(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	res, err := c.Repo.FindAllVideo()
	if err != nil {
		return nil, err
	}
	return &pb.FindAllResponse{
		Status: http.StatusOK,
		Videos: res,
	}, nil
}

// function to segment the video using ffmpeg and storing it as playlist
func CreatePlaylistAndSegments(filePath string, folderPath string) error {
	segmentDuration := 3
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", filePath,
		"-profile:v", "baseline", // baseline profile is compatible with most devices
		"-level", "3.0",
		"-start_number", "0", // start number segments from 0
		"-hls_time", strconv.Itoa(segmentDuration), //duration of each segment in second
		"-hls_list_size", "0", // keep all segments in the playlist
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", folderPath),
	)
	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v \nOutput: %s ", err, string(output))
	}
	return nil
}
