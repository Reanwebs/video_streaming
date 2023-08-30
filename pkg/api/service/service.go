package service

import (
	"bytes"
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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
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

		// LogLevel: aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	if err != nil {
		return err
	}
	// Create an S3 uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	fileuid := uuid.New()
	fileName := fileuid.String()

	var buffer bytes.Buffer

	s3Path := "reanweb/" + fileName + ".mp4"

	// Upload the video data from the buffer to S3
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String("reanweb"),
		Key:    aws.String(s3Path),
		Body:   bytes.NewReader(buffer.Bytes()),
	})
	if err != nil {
		fmt.Println("error here")
		return err
	}

	// Saving the video id to a database
	video_id, err := c.Repo.CreateVideoid(fileName, s3Path)
	if err != nil {
		fmt.Println("\n\n..\t\t", err)
		return err
	}

	fmt.Println("err nil", video_id)
	// Sending a response and closing the sending stream of bytes
	return stream.SendAndClose(&pb.UploadVideoResponse{
		Status:  http.StatusOK,
		Message: "Video successfully uploaded.",
		VideoId: video_id,
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

func (c *VideoServer) FindAllVideo(ctx context.Context, req *pb.FindAllRequest) (*pb.FindAllResponse, error) {
	// res, err := c.Repo.FindAllVideo()
	// if err != nil {
	// 	return nil, err
	// }
	// return &pb.FindAllResponse{
	// 	Status: http.StatusOK,
	// 	Videos: res,
	// }, nil
	// Set your AWS credentials
	// \AWS_REGION=ap-south-1
	// AWS_ACCESS_KEY_ID=AKIARQS5LHBSUWBUPHPM
	// AWS_SECRET_ACCESS_KEY=tf3hsQxBLxpFLeS+fLycrjzwwo4+XIDKztcj2rok
	creds := credentials.NewStaticCredentials("AKIARQS5LHBSUWBUPHPM", "tf3hsQxBLxpFLeS+fLycrjzwwo4+XIDKztcj2rok", "")

	// Create a new session using your credentials
	sess := session.Must(session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("ap-south-1"), // Change to your desired region
	}))

	// Create an S3 client
	s3Client := s3.New(sess)

	// Specify the bucket and object key
	bucketName := "reanweb"
	objectKey := "reanweb/7cd03aef-21f3-41e5-b235-21ccdd5c9152.mp4"

	// Create a file to save the downloaded video
	file, err := os.Create("downloaded_video.mp4")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return nil, err
	}
	defer file.Close()

	// Get the video object from S3
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}

	// Download the object and write it to the file
	result, err := s3Client.GetObject(getObjectInput)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok && aerr.Code() == s3.ErrCodeNoSuchKey {
			fmt.Println("Object not found:", err)
		} else {
			fmt.Println("Error downloading object:", err)
		}
		return nil, err
	}
	defer result.Body.Close()

	_, err = file.ReadFrom(result.Body)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return nil, err
	}

	res, err := c.Repo.FindAllVideo()
	if err != nil {
		return nil, err
	}
	fmt.Println("Video downloaded successfully!")
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
