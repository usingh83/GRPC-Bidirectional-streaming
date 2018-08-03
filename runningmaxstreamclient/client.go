package main;

import (
	"time"
	"fmt"
	"io"
	"context"
	"log"
	"google.golang.org/grpc"
	"grpc-practice/runningmaxstream/runningmaxstreampb"
)

func main(){
	cc,err:=grpc.Dial("localhost:50052",grpc.WithInsecure())
	if err!=nil{
		log.Fatalf("Listener creation failed %v",err)
	}
	defer cc.Close()
	c:=runningmaxstreampb.NewRunningmaxstreamClient(cc)
	stream,err:=c.CurrentMax(context.Background())
	if err!=nil{
		log.Fatalf("Stream creation failed %v",err)
	}
	x:=[]int64{1,5,3,6,2,20}
	waitc:=make(chan struct{})
	go func(){
		for _,n := range x{
			fmt.Println(n)
			stream.Send(&runningmaxstreampb.NumberRequest{
				Request:n,
			})
			time.Sleep(1000*time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func(){
		for{
			response,err:=stream.Recv()
			if err==io.EOF{
				break
			}
			if err!=nil{
				log.Fatalf("Stream response failure %v",err)
				break
			}
			fmt.Println("reponse is ",response.GetResponse())
		}
		close(waitc)
	}()
	<-waitc
}