// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package southbound

import (
	"github.com/golang/mock/gomock"
	"github.com/onosproject/onos-topo/pkg/northbound/device"
	"github.com/onosproject/onos-ztp/pkg/southbound/mock_device"
	"gotest.tools/assert"
	"io"
	"testing"
	"time"
)

func Test_Basics(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_device.NewMockDeviceServiceClient(ctrl)
	stream := mock_device.NewMockDeviceService_ListClient(ctrl)

	m.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(stream, nil)
	stream.EXPECT().Recv().
		Return(&device.ListResponse{
			Type: device.ListResponse_ADDED,
			Device: &device.Device{
				ID: "foobar",
			},
		}, nil)
	stream.EXPECT().Recv().
		Return(nil, io.ErrClosedPipe)
	stream.EXPECT().Recv().
		Return(nil, io.EOF)

	monitor := DeviceMonitor{m}
	ch := make(chan *device.Device)
	err := monitor.Start(ch)
	assert.NilError(t, err, "unexpected error")

	dev := <-ch
	assert.Assert(t, dev.GetID() == "foobar", "incorrect device")

	time.Sleep(100 * time.Millisecond)
}

func Test_ListError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_device.NewMockDeviceServiceClient(ctrl)

	m.EXPECT().List(gomock.Any(), gomock.Any()).
		Return(nil, io.ErrUnexpectedEOF)

	monitor := DeviceMonitor{m}
	ch := make(chan *device.Device)
	err := monitor.Start(ch)
	assert.Error(t, err, "unexpected EOF", "wrong error")
	time.Sleep(100 * time.Millisecond)
}
