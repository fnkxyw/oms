package main

import (
	"context"
	"log"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "gitlab.ozon.dev/akugnerevich/homework.git/pkg/PuP-service/v1"
)

const (
	numClients     = 1
	numRequests    = 100000
	requestTimeout = 5 * time.Second
)

func stressTestAcceptOrder(address string, wg *sync.WaitGroup, goroutineID int) {
	defer wg.Done()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPupServiceClient(conn)

	for i := 1; i <= numRequests; i++ {
		uniqueID := (goroutineID * numRequests) + i

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		req := &pb.AcceptOrderRequest{
			OrderId:       uint32(uniqueID),
			UserId:        uint32(uniqueID),
			KeepUntilDate: timestamppb.New(time.Now().Add(24 * time.Hour)),
			Weight:        1,
			Price:         1,
			PackageType:   "box",
			NeedWrapping:  false,
		}

		_, err := client.AcceptOrder(ctx, req)
		if err != nil {
			log.Printf("AcceptOrder request failed: %v", err)
		}
	}
}

func stressTestListOrders(address string, wg *sync.WaitGroup, goroutineID int) {
	defer wg.Done()

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewPupServiceClient(conn)

	for i := 1; i <= numRequests; i++ {
		uniqueID := (goroutineID * numRequests) + i

		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		defer cancel()

		req := &pb.ListOrdersRequest{
			UserId: uint32(uniqueID),
			InPup:  true,
			Count:  int32(0),
		}

		_, err := client.ListOrders(ctx, req)
		if err != nil {
			log.Printf("ListOrders request failed: %v", err)
		}
	}
}

func main() {
	cpuProfile, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatalf("could not create CPU profile: %v", err)
	}
	defer cpuProfile.Close()

	if err := pprof.StartCPUProfile(cpuProfile); err != nil {
		log.Fatalf("could not start CPU profile: %v", err)
	}
	defer pprof.StopCPUProfile()

	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go stressTestAcceptOrder("localhost:7002", &wg, i+1)
	}
	wg.Wait()

	durationAcceptOrder := time.Since(start)
	rpsAcceptOrder := float64(numClients*numRequests) / durationAcceptOrder.Seconds()
	log.Printf("RPS AcceptOrder: %.2f", rpsAcceptOrder)
	log.Printf("Общее время выполнения на запись: %v", durationAcceptOrder)

	start = time.Now()
	for i := 0; i < numClients; i++ {
		wg.Add(1)
		go stressTestListOrders("localhost:7002", &wg, i+1)
	}
	wg.Wait()

	durationListOrders := time.Since(start)
	rpsListOrders := float64(numClients*numRequests) / durationListOrders.Seconds()
	log.Printf("RPS ListOrders: %.2f", rpsListOrders)
	log.Printf("Общее время выполнения на чтение: %v", durationListOrders)

	memProfile, err := os.Create("mem.prof")
	if err != nil {
		log.Fatalf("could not create memory profile: %v", err)
	}
	defer memProfile.Close()

	if err := pprof.WriteHeapProfile(memProfile); err != nil {
		log.Fatalf("could not write memory profile: %v", err)
	}

	log.Println("Профилирование завершено. Профили сохранены в файлы cpu.prof и mem.prof")
}
