package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNvmeAvailable(t *testing.T) {
	assert.True(t, NvmeAvailable())
}

func TestNvmeTempOutput(t *testing.T) {
	input := []byte(`Smart Log for NVME device:nvme0 namespace-id:ffffffff
critical_warning                    : 0
temperature                         : 37 C
available_spare                     : 100%
available_spare_threshold           : 10%
percentage_used                     : 0%
data_units_read                     : 553,959
data_units_written                  : 1,175,224
host_read_commands                  : 21,683,776
host_write_commands                 : 21,524,281
controller_busy_time                : 77
power_cycles                        : 20
power_on_hours                      : 81
unsafe_shutdowns                    : 5
media_errors                        : 0
num_err_log_entries                 : 0
Warning Temperature Time            : 0
Critical Composite Temperature Time : 0
Temperature Sensor 1                : 36 C
Temperature Sensor 2                : 38 C
Temperature Sensor 5                : 45 C
Thermal Management T1 Trans Count   : 7
Thermal Management T2 Trans Count   : 0
Thermal Management T1 Total Time    : 5309
Thermal Management T2 Total Time    : 0`)

	assert.Equal(t, parseNvmeTemperatureOutput(input), "37C")
}
