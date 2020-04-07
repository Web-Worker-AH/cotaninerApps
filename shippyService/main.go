package main 

import (

	"context"
	"log"
	"sync"
	
	micro "github.com/micro/go-micro/v2"
	pb "shippyService/proto"
)



type repository interface {

	Create(*pb.Consignment) (*pb.Consignment,error)
	GetAll() []*pb.Consignment
}

type Repository struct {

	mu	sync.Mutex
	consignments []*pb.Consignment

}


func (repo *Repository) Create(consignment *pb.Consignment) (*pb.Consignment,error) {
	repo.mu.Lock()
	updated := append(repo.consignments , consignment)
	repo.consignments = updated 
	repo.mu.Unlock()

	return consignment,nil

}
func (repo *Repository) GetAll() []*pb.Consignment {
	return repo.consignments
}

type service struct {
	repo repository
}

func (s *service) CreateConsignment (ctx context.Context , req *pb.Consignment,res *pb.Response) error {
	consignment,err := s.repo.Create(req)
	if err!=nil{
		return err
	}
	res.Created = true 
	res.Consignment = consignment
	return nil 
}

func (s *service) GetConsignment(ctx context.Context,req *pb.GetRequest,res *pb.Response) error {
	consignments := s.repo.GetAll()
	res.Consignments = consignments
	return nil 
}

func main(){
	repo := &Repository{}
	s := micro.NewService(micro.Name("shippyService"))
	pb.RegisterShippingServiceHandler(s.Server(), &service{repo})
	//registering our service implementation with the GRPC server



	if err := s.Run();err!=nil{
		log.Fatalf("failed to serve %v",err)
	}
}