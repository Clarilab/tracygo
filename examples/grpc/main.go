package main

import (
	"context"
	"log"
	"net"
	"strings"
	"sync"

	"github.com/Clarilab/tracygo/v2"
	"github.com/Clarilab/tracygo/v2/examples/grpc/filetransfer"
	tracygrpc "github.com/Clarilab/tracygo/v2/middleware/grpc"
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

const (
	headerCorrelationID = "X-Correlation-ID"
	headerRequestID     = "X-Request-ID"

	correlationValue = "Zitronenbaum"
)

var wg *sync.WaitGroup //nolint:gochecknoglobals // intended use

func main() {
	wg = new(sync.WaitGroup)

	go runServer()

	// Set up a connection to the server
	grpcClient, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer grpcClient.Close()

	// Create a context with metadata
	md := metadata.Pairs(
		"authorization", "Bearer some-token",
		"custom-header", "custom-value",
		headerCorrelationID, correlationValue,
	)

	ctx := metadata.NewOutgoingContext(context.Background(), md)

	client := filetransfer.NewFileServiceClient(grpcClient)

	// Prepare to receive response metadata
	var trailer metadata.MD

	wg.Add(1)

	// Call the RPC method with the context
	res, err := client.UploadFile(ctx, &filetransfer.UploadFileRequest{
		FileName:    "example.txt",
		FileContent: []byte("Hello, gRPC!"),
	},
		grpc.Trailer(&trailer), // Capture the trailing metadata
	)
	if err != nil {
		log.Printf("could not upload file: %v", err)

		return
	}

	wg.Wait()

	var correlationID string
	var requestID string

	if values := trailer[strings.ToLower(headerCorrelationID)]; len(values) == 1 {
		correlationID = values[0]
	}

	if values := trailer[strings.ToLower(headerRequestID)]; len(values) == 1 {
		requestID = values[0]
	}

	if correlationID != correlationValue {
		panic("could not retrieve correlation id")
	}

	if requestID == "" {
		panic("could not retrieve request id")
	}

	log.Printf("UploadFile response: %v", res)
}

type server struct {
	filetransfer.UnimplementedFileServiceServer
}

func (s *server) UploadFile(ctx context.Context, req *filetransfer.UploadFileRequest) (*filetransfer.UploadFileResponse, error) {
	defer wg.Done()

	// make sure the context has the correlation and request id
	// so f.e. the logger can extract them later
	ensureInContext(ctx)

	// make sure the gRPC contexts have the correlation and request id
	ensureInOutgoingGRPCContext(ctx)

	// Process the file upload
	log.Printf("Received file: %s", req.GetFileName())

	return &filetransfer.UploadFileResponse{
		Success: true,
		Message: "File uploaded successfully",
	}, nil
}

func ensureInContext(ctx context.Context) {
	correlationID, ok := ctx.Value(headerCorrelationID).(string)
	if !ok || correlationID != correlationValue {
		panic("could not retrieve correlation id")
	}

	requestID, ok := ctx.Value(headerRequestID).(string)
	if !ok || requestID == "" {
		panic("could not retrieve request id")
	}
}

func ensureInOutgoingGRPCContext(ctx context.Context) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		panic("could not retrieve metadata")
	}

	var correlationID string
	var requestID string

	if values := md[strings.ToLower(headerCorrelationID)]; len(values) == 1 {
		correlationID = values[0]
	}

	if values := md[strings.ToLower(headerRequestID)]; len(values) == 1 {
		requestID = values[0]
	}

	if correlationID != correlationValue {
		panic("could not retrieve correlation id")
	}

	if requestID == "" {
		panic("could not retrieve request id")
	}
}

func runServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	tracer := tracygo.New()

	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpcmiddleware.ChainUnaryServer(tracygrpc.CheckTracingIDs(tracer)),
		),
	)

	filetransfer.RegisterFileServiceServer(s, new(server))

	log.Printf("Server is listening on port 50051...")

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
