package checks

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/containers/podman/v3/pkg/bindings"
	"github.com/containers/podman/v3/pkg/bindings/images"
)

func socket() string {
	xdgPath, ok := os.LookupEnv("XDG_RUNTIME_DIR")
	if !ok {
		xdgPath = "/run"
	}

	socketPath := []string{
		xdgPath,
		"podman",
		"podman.sock",
	}
	return fmt.Sprintf("unix://%s", filepath.Join(socketPath...))
}

func connection() context.Context {
	ctx, err := bindings.NewConnection(context.Background(), socket())
	if err != nil {
		fmt.Println("error:")
		panic(err)
	}
	return ctx
}

func getImage(ctx context.Context, image string) {
	_, err := images.Pull(ctx, image, &images.PullOptions{})
	if err != nil {
		panic(err)
	}
}

func listImages(ctx context.Context) {
	results, err := images.List(ctx, &images.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, img := range results {
		fmt.Println(img.Names)
	}
}

func inspectImage(ctx context.Context, nameOrID string) {
	report, err := images.GetImage(ctx, nameOrID, &images.GetOptions{})
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", report)
}

func saveImage(ctx context.Context, nameOrIDs []string) {
	outfile, err := os.Create("image.tar.gz")
	if err != nil {
		panic(err)
	}

	var compress bool = true
	err = images.Export(ctx, nameOrIDs, outfile, &images.ExportOptions{
		Compress: &compress,
	})
	if err != nil {
		panic(err)
	}
}

func PodmanSocket() {
	ctx := connection()
	getImage(ctx, "docker.io/busybox")
	listImages(ctx)
	inspectImage(ctx, "docker.io/busybox")
	saveImage(ctx, []string{"docker.io/busybox"})
}
