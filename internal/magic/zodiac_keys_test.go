package magic

import "testing"

func TestZodiacKeysConfirmedTwelve(t *testing.T) {
	if len(ZodiacKeys) != 12 {
		t.Fatalf("len(ZodiacKeys) = %d, want 12 (one per zodiac sign)", len(ZodiacKeys))
	}
	seen := map[string]bool{}
	for _, k := range ZodiacKeys {
		if k.Sign == "" || k.Metal == "" {
			t.Errorf("ZodiacKey %+v missing expected fields", k)
		}
		seen[k.Sign] = true
	}
	if len(seen) != 12 {
		t.Errorf("ZodiacKeys has %d distinct Signs, want 12 (no duplicates)", len(seen))
	}
}
