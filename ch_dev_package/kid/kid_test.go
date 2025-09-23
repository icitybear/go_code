package kid

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestGetMachineID(t *testing.T) {
	// 固定种子，用于单元测试
	rand.Seed(1)
	testCases := []struct {
		exp      int64
		hostname string
	}{
		// bcdfghjklmnpqrstvwxz2456789
		{4300268, "chatroom-c89744558-ldq6h"}, // 8*27^4 + 2*27^3 + 12*27^2 + 23*27 + 5 = 4300268
		{3101377, "live-api-549ffcf84-h5tk5"}, // 5*27^4 + 22*27^3 + 15*27^2 + 7*27 + 22 = 3101377
		{10436627, "g1-5966bfcdff-zwjm2"},     // 19*27^4 + 17*27^3 + 6*27^2 + 9*27 + 20 = 10436627
		{8350655, "local_hostname"},           // 种子为1时，rand.Int63n(0xFFFFFF) 固定返回 8350655
	}

	for i, tc := range testCases {
		id := getMachineID(tc.hostname)
		if id != tc.exp {
			t.Errorf("[#%d]获取机器ID失败. hostname:%s, 期望:%d, 实际：%d", i, tc.hostname, tc.exp, id)
		}
	}
}

func TestNew(t *testing.T) {
	machineID = getMachineID("chatroom-c89744558-ldq6h")
	ts := time.Now().UnixNano() / nanoms
	ids := make([]ID, 0, 10)
	for i := 0; i < 10; i++ {
		ids = append(ids, New())
	}
	for i, id := range ids {
		if (id.Time().UnixNano() / nanoms) != ts {
			t.Errorf("[#%d]获取 kid 时间失败，期望:%d, 实际:%d", i, ts, id.Time().UnixNano()/nanoms)
		}
		if id.SequenceID() != int64(i+1) {
			t.Errorf("[#%d]获取 kid 自增ID失败, 期望:%d, 实际:%d", i, i, id.SequenceID())
		}
		v, err := json.Marshal(id)
		if err != nil {
			t.Errorf("[#%d]kid json序列化失败, err:%v", i, err)
		}
		if string(v) != fmt.Sprintf("%q", id) {
			t.Errorf("[#%d]kid json序列化非预期，期望:%q, 实际:%q", i, id, v)
		}
		x, err := xml.Marshal(id)
		if err != nil {
			t.Errorf("[#%d]kid text序列化失败, err：%v", i, err)
		}
		if string(x) != fmt.Sprintf("<ID>%s</ID>", id) {
			t.Errorf("[#%d]kid text序列化非预期，期望:<ID>%s</ID>, 实际:%s", i, id, x)
		}
		if id.MachineID() != machineID {
			t.Errorf("[#%d]获取kid 机器ID失败，期望：%d, 实际:%d", i, machineID, id.MachineID())
		}
	}
}

func TestEopch(t *testing.T) {
	eopch := Eopch
	var t42bit int64 = 2 << 40
	now := time.Now().UnixNano() / nanoms
	if now-eopch < t42bit {
		Eopch = -(t42bit - now)
	}

	id := New()
	if (id.Time().UnixNano() / nanoms) != now {
		t.Fatalf("设置纪元失败，期望:%d, 实际:%d", now, id.Time().UnixNano()/nanoms)
	}
	if len(id.String()) != 16 {
		t.Errorf("测试kid字符最大长度失败，期望:16，实际:%d", len(id.String()))
	}

	Eopch = eopch
}

// 命令行相当 /usr/local/go/bin/go test -timeout 30s -run ^TestSetMachineID$ kid
func TestSetMachineID(t *testing.T) {
	var mid uint64 = 123
	SetMachineID(mid)
	id := New()

	if id.MachineID() != int64(mid) {
		t.Errorf("设置机器ID失败, 期望:%d; 实际:%d", mid, id.MachineID())
	}
}

func TestSetEncodingDict(t *testing.T) {
	dict := encoding

	err := SetEncodingDict("abc")
	if !errors.Is(err, ErrInvalidEncodingLength) {
		t.Errorf("设置编码字典失败，期望:%v; 实际:%v", ErrInvalidEncodingLength, err)
	}
	err = SetEncodingDict("0123456789abcdefghijklmnopqrstuu")
	if !errors.Is(err, ErrRepeatedCharacters) {
		t.Errorf("设置编码字典失败，期望:%v; 实际:%v", ErrRepeatedCharacters, err)
	}
	err = SetEncodingDict("0123456789abcdefghijklmnopqrstu\t")
	if !errors.Is(err, ErrOnlyPrintableCharacters) {
		t.Errorf("设置编码字典失败，期望:%v; 实际:%v", ErrOnlyPrintableCharacters, err)
	}

	now := time.Now().UnixNano() / nanoms
	id1 := New()
	sid1 := id1.String()

	ndict := "abcdefghijklmnopqrstuv0123456789"
	err = SetEncodingDict(ndict)
	if err != nil {
		t.Fatalf("设置编码字典失败，err:%v", err)
	}
	id2 := New()
	sid2 := id2.String()
	if (id1.Time().UnixNano()/nanoms) != now || (id2.Time().UnixNano()/nanoms) != now {
		t.Errorf("设置编码字典失败，期望时间一致，实际不一致. id1:%d; id2:%d", id1.Time().UnixNano()/nanoms, id2.Time().UnixNano()/nanoms)
	}

	for i, c1 := range sid1 {
		c2 := sid2[i]
		i1 := strings.IndexRune(dict, c1)
		i2 := strings.IndexByte(ndict, c2)
		if i1 != i2 {
			// 最后一位 c2 会 比 c1 大 1
			if i == len(sid2)-1 && i1+1 == i2 {
				continue
			}
			t.Errorf("设置编码字典失败，编码不正确. 期望:%d; 实际：%d; id1:%s; id2:%s", i1, i2, sid1, sid2)
		}
	}

	encoding = dict
}
