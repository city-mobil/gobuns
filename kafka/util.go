package kafka

import "github.com/segmentio/kafka-go"

func balancerToKafkaBalancer(typ string) kafka.Balancer {
	switch typ {
	case BalancerRoundRobin:
		return &kafka.RoundRobin{}
	case BalancerCRC32:
		return &kafka.CRC32Balancer{
			Consistent: true,
		}
	case BalancerMurMur2:
		return &kafka.Murmur2Balancer{
			Consistent: true,
		}
	}

	return &kafka.RoundRobin{}
}
