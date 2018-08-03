package main;

import (
	"io"
	"google.golang.org/grpc"
	"log"
	"net"
	"grpc-practice/runningmaxstream/runningmaxstreampb"
)

type server struct{}

func (*server) CurrentMax(stream runningmaxstreampb.Runningmaxstream_CurrentMaxServer) error{
	max:=int64(0)
	for{
		request,err:=stream.Recv()
		if err==io.EOF{
			break
		}
		if err!=nil{
			log.Fatalf("stream processing failed %v",err)
		}
		num:=request.GetRequest()
		if max<num{
			stream.Send(
				&runningmaxstreampb.NumberResponse{
					Response:num,
				})
			max=num
		}
	}
	return nil
}

func main(){
	lis,err:=net.Listen("tcp","0.0.0.0:50052")
	if err!=nil{
		log.Fatalf("listener creation failed %v",err)
	}
	s:=grpc.NewServer()
	runningmaxstreampb.RegisterRunningmaxstreamServer(s,&server{})
	if err:=s.Serve(lis);err!=nil{
		log.Fatalf("Server creation failed %v",err)
	}
}
