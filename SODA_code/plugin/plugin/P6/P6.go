package main

import (
	// "../../pluginlog"
	"encoding/hex"
	"github.com/ethereum/collector"
	// "math/big"
	// "fmt"
	"github.com/json-iterator/go"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary
// var logger pluginlog.ErrTxLog

// var	log_map map[int]map[string]int   // store standard result
var event_flag int
var standard_func_flag int

type RegisterInfo struct {
	PluginName string   `json:"pluginname"`
	OpCode     map[string]string `json:"option"`
}

func Register() []byte {
	// fmt.Println("enter run")
	var data = RegisterInfo{
		PluginName: "P6",
		OpCode: map[string]string{"EXTERNALINFOSTART":"handle_EXTERNALINFOSTART", "EXTERNALINFOEND":"handle_EXTERNALINFOEND", "EVENT":"handle_EVENT"},
	}

	standard_func_flag = 0
	event_flag = 0

	retInfo, err := json.Marshal(&data)
	if err != nil {
		return nil
	}

	return retInfo
}

func handle_EXTERNALINFOSTART(m *collector.CollectorDataT) (byte ,string){
	standard_func_flag = 0
	event_flag = 0
	if m.TransInfo.CallType == "CALL"{   // external call, get contract name and input, check if the method is in the jumptable
		input := hex.EncodeToString(m.TransInfo.CallInfo.InputData)
		ll := len(input)
		if ll >= 8{
			methodid := strings.ToLower(input[0:8])
			if methodid == "a9059cbb" || methodid == "23b872dd"{
				standard_func_flag = 1
			}
		}
	}
	return 0x00,""
}

func handle_EVENT(m *collector.CollectorDataT) (byte ,string){
	if len(m.InsInfo.OpArgs) < 3{
		return 0x00,""
	}
	len_data := (4 - (len(m.InsInfo.OpArgs) - 2)) * 64
	data := hex.EncodeToString(m.InsInfo.RetArgs)
	true_len_data := len(data)
	if true_len_data == len_data{
		event := strings.ToLower(m.InsInfo.OpArgs[2])
		if event == "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"{
			event_flag = 1
		}
	}
	return 0x00,""
}

func handle_EXTERNALINFOEND(m *collector.CollectorDataT) (byte ,string){
	if m.TransInfo.IsSuccess{
		if standard_func_flag == 1 && event_flag == 0{
			return 0x01,""
		}
	}
	return 0x00,""
}



