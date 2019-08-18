package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAvailable(t *testing.T) {
	assert.True(t, NVidiaSMIAvailable())
}

func TestTempOutput(t *testing.T) {
	input := []byte(`==============NVSMI LOG==============

Timestamp                           : Tue Jan  1 19:03:42 2019
Driver Version                      : 415.25
CUDA Version                        : 10.0

Attached GPUs                       : 1
GPU 00000000:01:00.0
    Temperature
        GPU Current Temp            : 31 C
        GPU Shutdown Temp           : 102 C
        GPU Slowdown Temp           : 99 C
        GPU Max Operating Temp      : N/A
        Memory Current Temp         : N/A
        Memory Max Operating Temp   : N/A
			`)
	assert.Equal(t, parseNvidiaTemperatureOutput(input), "31C")
}

func TestUtilizationOutput(t *testing.T) {
	input := []byte(`

==============NVSMI LOG==============

Timestamp                           : Tue Jan  1 19:13:44 2019
Driver Version                      : 415.25
CUDA Version                        : 10.0

Attached GPUs                       : 1
GPU 00000000:01:00.0
    Utilization
        Gpu                         : 6 %
        Memory                      : 6 %
        Encoder                     : 0 %
        Decoder                     : 0 %
    GPU Utilization Samples
        Duration                    : 18446744073709.21 sec
        Number of Samples           : 99
        Max                         : 30 %
        Min                         : 0 %
        Avg                         : 0 %
    Memory Utilization Samples
        Duration                    : 18446744073709.21 sec
        Number of Samples           : 99
        Max                         : 8 %
        Min                         : 5 %
        Avg                         : 0 %
    ENC Utilization Samples
        Duration                    : 18446744073709.21 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %
    DEC Utilization Samples
        Duration                    : 18446744073709.21 sec
        Number of Samples           : 99
        Max                         : 0 %
        Min                         : 0 %
        Avg                         : 0 %`)
	assert.Equal(t, parseUtilizationOutput(input), "6%")
}

func TestUtilizationOutputNA(t *testing.T) {
	input := []byte(`

==============NVSMI LOG==============

Timestamp                           : Fri May 10 23:54:29 2019
Driver Version                      : 418.56
CUDA Version                        : 10.1

Attached GPUs                       : 1
GPU 00000000:08:00.0
    Utilization
        Gpu                         : N/A
        Memory                      : N/A
        Encoder                     : N/A
        Decoder                     : N/A
    GPU Utilization Samples
        Duration                    : N/A
        Number of Samples           : N/A
        Max                         : N/A
        Min                         : N/A
        Avg                         : N/A
    Memory Utilization Samples
        Duration                    : N/A
        Number of Samples           : N/A
        Max                         : N/A
        Min                         : N/A
        Avg                         : N/A
    ENC Utilization Samples
        Duration                    : N/A
        Number of Samples           : N/A
        Max                         : N/A
        Min                         : N/A
        Avg                         : N/A
    DEC Utilization Samples
        Duration                    : N/A
        Number of Samples           : N/A
        Max                         : N/A
        Min                         : N/A
        Avg                         : N/A`)
	assert.Equal(t, parseUtilizationOutput(input), "0%")
}
