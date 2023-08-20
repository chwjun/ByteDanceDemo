package util

import "testing"

func BenchmarkGetTotalLikeCounts(b *testing.B) {

	videoIDs := []uint{1, 2, 3}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := GlobalRedisClient.GetTotalLikeCounts(videoIDs)
		if err != nil {
			b.Fatal(err)
		}
	}
}
