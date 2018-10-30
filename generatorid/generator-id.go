package generatorid

import (
	"fmt"
	"strconv"

	"github.com/golang/glog"
	"github.com/sony/sonyflake"
)

func NewIDGenerator() *IDGenerator {
	settings := sonyflake.Settings{}

	return &IDGenerator{
		generator: sonyflake.NewSonyflake(settings),
	}
}

type IDGenerator struct {
	generator *sonyflake.Sonyflake
}

func (o *IDGenerator) Next() (string, error) {
	id, err := o.generator.NextID()
	if err != nil {
		msg := fmt.Sprintf("flake.NextID() error generate next random id. Error: %v", err)

		return "", fmt.Errorf(msg)
	}

	glog.Infof("IDGenerator.Next() Generate ID: %v", id)

	return strconv.FormatUint(id, 10), nil
}
