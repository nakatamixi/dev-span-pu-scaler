package scaler

import (
	"context"
	"fmt"

	"github.com/nktks/dev-span-pu-scaler/internal/spanner"
)

type Scaler struct {
	client *spanner.Client
}

func NewScaler(client *spanner.Client) *Scaler {
	return &Scaler{
		client: client,
	}
}

func (s *Scaler) Execute(ctx context.Context, buffer int) error {
	ins, err := s.client.Instance(ctx)
	if err != nil {
		return fmt.Errorf("could not get instance %s: %v", s.client.FQDN(), err)
	}
	currentPU := int(ins.ProcessingUnits)
	fmt.Printf("current instance pu count %s [%d]\n", s.client.FQDN(), currentPU)
	dbCount, err := s.client.DBCount(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("current db count %s [%d]\n", s.client.FQDN(), dbCount)
	calc := NewPUCalculator(dbCount, buffer)
	desiredPU := calc.DesiredPU()
	fmt.Printf("desired pu count %s [%d]\n", s.client.FQDN(), desiredPU)
	if currentPU == desiredPU {
		switch {
		case calc.IsUpperLimit(currentPU):
			fmt.Printf("%s filled with many dbs.  [%d]\n", s.client.FQDN(), dbCount)
		case dbCount == 0:
			fmt.Printf("%s has no db.\n", s.client.FQDN())
		}
		return nil
	}

	ins.ProcessingUnits = int32(desiredPU)

	switch {
	case currentPU > desiredPU:
		fmt.Printf("scale down %s [%d] -> [%d]\n", s.client.FQDN(), currentPU, desiredPU)
		if err := s.client.UpdatePU(ctx, ins); err != nil {
			return err
		}
	case currentPU < desiredPU:
		fmt.Printf("scale up %s [%d] -> [%d]\n", s.client.FQDN(), currentPU, desiredPU)
		if err := s.client.UpdatePU(ctx, ins); err != nil {
			return err
		}
	}

	return nil
}
