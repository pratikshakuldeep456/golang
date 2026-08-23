package ratelimiter

import "time"

func RLTest() {

	Zohoc1 := &ZohoProvider{}
	ProviderClient := &ProviderClient{Provider: Zohoc1, RL: &TokenBucket{
		Cap:            10,
		Token:          10,
		RefillRate:     2,
		LastRefillTime: time.Now(),
	}}

	ProviderClient.Call()

}
