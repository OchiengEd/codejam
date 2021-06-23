package checks

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type ImageClient struct {
	*client.Client
}

// DockerClient get returns client to interact with images
func DockerClient() *ImageClient {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	return &ImageClient{
		Client: cli,
	}
}

func (r *ImageClient) Operations() {
	ctx := context.Background()
	r.pullImage(ctx, "docker.io/busybox")
	r.listImages(ctx)
	r.saveImage(ctx, []string{"docker.io/busybox"})
}

func (r *ImageClient) pullImage(ctx context.Context, nameOrID string) {
	fmt.Printf("pull %s image...", nameOrID)
	reader, err := r.Client.ImagePull(ctx, nameOrID, types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, reader)
	defer reader.Close()
	fmt.Println(reader)
}

func (r *ImageClient) listImages(ctx context.Context) {
	fmt.Println("list images...")
	summary, err := r.Client.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		panic(err)
	}

	for _, img := range summary {
		fmt.Println(img.ID)
	}
}

func (r *ImageClient) saveImage(ctx context.Context, imgIDs []string) {
	fmt.Println("saving image locally")
	outfile, err := os.Create("image.tar.gz")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	reader, err := r.Client.ImageSave(ctx, imgIDs)
	if err != nil {
		panic(err)
	}
	defer reader.Close()

	outfile.ReadFrom(reader)
	fmt.Println("finished!")
}
