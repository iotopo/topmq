package utils

import "testing"

func makeTestDebugBytesData(size int) (data []byte) {
	data = make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i % 256)
	}
	return
}

func testDebugBytesEqual(t *testing.T, size int, expected string) {
	built := DebugBytes(makeTestDebugBytesData(size))
	t.Logf("DebugBytes(%d):\n%s\n", size, built)
	if built != expected {
		t.Errorf("failed test DebugBytes for size %d", size)
	}
}

func TestDebugBytes0(t *testing.T) {
	expected := `
(0 bytes)`
	testDebugBytesEqual(t, 0, expected[1:])
}

func TestDebugBytes1(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00          |             |             |             |
   |-------------|-------------|-------------|-------------|
(1 bytes)`
	testDebugBytesEqual(t, 1, expected[1:])
}

func TestDebugBytes4(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00 01 02 03 |             |             |             |
   |-------------|-------------|-------------|-------------|
(4 bytes)`
	testDebugBytesEqual(t, 4, expected[1:])
}

func TestDebugBytes7(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00 01 02 03 | 04 05 06    |             |             |
   |-------------|-------------|-------------|-------------|
(7 bytes)`
	testDebugBytesEqual(t, 7, expected[1:])
}

func TestDebugBytes8(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00 01 02 03 | 04 05 06 07 |             |             |
   |-------------|-------------|-------------|-------------|
(8 bytes)`
	testDebugBytesEqual(t, 8, expected[1:])
}

func TestDebugBytes12(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00 01 02 03 | 04 05 06 07 | 08 09 0a 0b |             |
   |-------------|-------------|-------------|-------------|
(12 bytes)`
	testDebugBytesEqual(t, 12, expected[1:])
}

func TestDebugBytes16(t *testing.T) {
	expected := `
   |-------------|-------------|-------------|-------------|
0: | 00 01 02 03 | 04 05 06 07 | 08 09 0a 0b | 0c 0d 0e 0f |
   |-------------|-------------|-------------|-------------|
(16 bytes)`
	testDebugBytesEqual(t, 16, expected[1:])
}

func TestDebugBytes17(t *testing.T) {
	expected := `
    |-------------|-------------|-------------|-------------|
00: | 00 01 02 03 | 04 05 06 07 | 08 09 0a 0b | 0c 0d 0e 0f |
10: | 10          |             |             |             |
    |-------------|-------------|-------------|-------------|
(17 bytes)`
	testDebugBytesEqual(t, 17, expected[1:])
}
